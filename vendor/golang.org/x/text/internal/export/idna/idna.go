// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_trieval.go gen_common.go

// http://www.unicode.org/reports/tr46

// Package idna implements IDNA2008 using the compatibility processing
// defined by UTS (Unicode Technical Standard) #46, which defines a standard to
// deal with the transition from IDNA2003.
//
// IDNA2008 (Internationalized Domain Names for Applications), is defined in RFC
// 5890, RFC 5891, RFC 5892, RFC 5893 and RFC 5894.
// UTS #46 is defined in http://www.unicode.org/reports/tr46.
// See http://unicode.org/cldr/utility/idna.jsp for a visualization of the
// differences between these two standards.
package idna // import "golang.org/x/text/internal/export/idna"

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/secure/bidirule"
	"golang.org/x/text/unicode/norm"
)

// A Profile defines the configuration of a IDNA mapper.
type Profile struct {
	Transitional    bool
	IgnoreSTD3Rules bool
	IgnoreDNSLength bool
	// ErrHandler      func(error)
}

// String reports a string with a description of the profile for debugging
// purposes. The string format may change with different versions.
func (p *Profile) String() string {
	s := ""
	if p.Transitional {
		s = "Transitional"
	} else {
		s = "NonTraditional"
	}
	if p.IgnoreSTD3Rules {
		s += ":NoSTD3Rules"
	}
	return s
}

var (
	// Resolve is the recommended profile for resolving domain names.
	// The configuration of this profile may change over time.
	Resolve = resolve

	// Transitional defines a profile that implements the Transitional mapping
	// as defined in UTS #46 with no additional constraints.
	Transitional = transitional

	// NonTransitional defines a profile that implements the Transitional
	// mapping as defined in UTS #46 with no additional constraints.
	NonTransitional = nonTransitional

	resolve         = &Profile{Transitional: true}
	transitional    = &Profile{Transitional: true}
	nonTransitional = &Profile{}

	// TODO: profiles
	// V2008: strict IDNA2008
	// Registrar: recommended for approving domain names.
)

// TODO: rethink error strategy

var (
	// errDisallowed indicates a domain name contains a disallowed rune.
	errDisallowed = errors.New("idna: disallowed rune")

	// errEmptyLabel indicates a label was empty.
	errEmptyLabel = errors.New("idna: empty label")
)

// process implements the algorithm described in section 4 of UTS #46,
// see http://www.unicode.org/reports/tr46.
func (p *Profile) process(s string, toASCII bool) (string, error) {
	var (
		b    []byte
		err  error
		k, i int
	)
	for i < len(s) {
		v, sz := trie.lookupString(s[i:])
		start := i
		i += sz
		// Copy bytes not copied so far.
		switch p.simplify(info(v).category()) {
		case valid:
			continue
		case disallowed:
			if err == nil {
				err = errDisallowed
			}
			continue
		case mapped, deviation:
			b = append(b, s[k:start]...)
			b = info(v).appendMapping(b, s[start:i])
		case ignored:
			b = append(b, s[k:start]...)
			// drop the rune
		case unknown:
			b = append(b, s[k:start]...)
			b = append(b, "\ufffd"...)
		}
		k = i
	}
	if k == 0 {
		// No changes so far.
		s = norm.NFC.String(s)
	} else {
		b = append(b, s[k:]...)
		if norm.NFC.QuickSpan(b) != len(b) {
			b = norm.NFC.Bytes(b)
		}
		// TODO: the punycode converters requires strings as input.
		s = string(b)
	}
	// TODO(perf): don't split.
	labels := strings.Split(s, ".")
	// Remove leading empty labels
	for len(labels) > 0 && labels[0] == "" {
		labels = labels[1:]
	}
	if len(labels) == 0 {
		return "", errors.New("idna: there are no labels")
	}
	// Find the position of the root label.
	root := len(labels) - 1
	if labels[root] == "" {
		root--
	}
	for i, label := range labels {
		// Empty labels are not okay, unless it is the last.
		if label == "" {
			if i <= root && err == nil {
				err = errEmptyLabel
			}
			continue
		}
		if strings.HasPrefix(label, acePrefix) {
			u, err2 := decode(label[len(acePrefix):])
			if err2 != nil {
				if err == nil {
					err = err2
				}
				// Spec says keep the old label.
				continue
			}
			labels[i] = u
			if err == nil {
				err = p.validateFromPunycode(u)
			}
			if err == nil {
				err = NonTransitional.validate(u)
			}
		} else if err == nil {
			err = p.validate(labels[i])
		}
	}
	if toASCII {
		for i, label := range labels {
			if !ascii(label) {
				a, err2 := encode(acePrefix, label)
				if err == nil {
					err = err2
				}
				labels[i] = a
			}
			n := len(labels[i])
			if !p.IgnoreDNSLength && err == nil && (n == 0 || n > 63) {
				if n != 0 || i != len(labels)-1 {
					err = fmt.Errorf("idna: label with invalid length %d", n)
				}
			}
		}
	}
	s = strings.Join(labels, ".")
	if toASCII && !p.IgnoreDNSLength && err == nil {
		// Compute the length of the domain name minus the root label and its dot.
		n := len(s) - 1 - len(labels[len(labels)-1])
		if len(s) < 1 || n > 253 {
			err = fmt.Errorf("idna: doman name with invalid length %d", n)
		}
	}
	return s, err
}

// acePrefix is the ASCII Compatible Encoding prefix.
const acePrefix = "xn--"

func (p *Profile) simplify(cat category) category {
	switch cat {
	case disallowedSTD3Mapped:
		if !p.IgnoreSTD3Rules {
			cat = disallowed
		} else {
			cat = mapped
		}
	case disallowedSTD3Valid:
		if !p.IgnoreSTD3Rules {
			cat = disallowed
		} else {
			cat = valid
		}
	case deviation:
		if !p.Transitional {
			cat = valid
		}
	case validNV8, validXV8:
		// TODO: handle V2008
		cat = valid
	}
	return cat
}

func (p *Profile) validateFromPunycode(s string) error {
	if !norm.NFC.IsNormalString(s) {
		return errors.New("idna: punycode is not normalized")
	}
	for i := 0; i < len(s); {
		v, sz := trie.lookupString(s[i:])
		if c := p.simplify(info(v).category()); c != valid && c != deviation {
			return fmt.Errorf("idna: invalid character %+q in expanded punycode", s[i:i+sz])
		}
		i += sz
	}
	return nil
}

// validate validates the criteria from Section 4.1. Item 1, 4, and 6 are
// already implicitly satisfied by the overall implementation.
func (p *Profile) validate(s string) error {
	if len(s) > 4 && s[2] == '-' && s[3] == '-' {
		return errors.New("idna: label starts with ??--")
	}
	if s[0] == '-' || s[len(s)-1] == '-' {
		return errors.New("idna: label may not start or end with '-'")
	}
	// TODO: merge the use of this in the trie.
	r, _ := utf8.DecodeRuneInString(s)
	if unicode.Is(unicode.M, r) {
		return fmt.Errorf("idna: label starts with modifier %U", r)
	}
	if !bidirule.ValidString(s) {
		return fmt.Errorf("idna: label violates Bidi Rule", r)
	}
	return nil
}

func (p *Profile) ToASCII(s string) (string, error) {
	return p.process(s, true)
}

func (p *Profile) ToUnicode(s string) (string, error) {
	return NonTransitional.process(s, false)
}

func ascii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
