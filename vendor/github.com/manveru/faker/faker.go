// faker generates fake data in various languages.
package faker

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	keyPattern       = regexp.MustCompile("#{[^}]+}")
	numerifyPattern  = regexp.MustCompile("#+")
	removeFancyChars = regexp.MustCompile("\\W")
)

type Faker struct {
	Language string
	Dict     map[string][]string
	Rand     Rand
}

type Rand interface {
	Float32() float32
	Float64() float64
	Int63n(int64) int64
	Intn(int) int
}

func New(lang string) (*Faker, error) {
	subDict, ok := Dict[lang]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No such language: %q", lang))
	}

	source := rand.NewSource(time.Now().UnixNano())
	return &Faker{
		Language: lang,
		Dict:     subDict,
		Rand:     rand.New(source),
	}, nil
}

func (f *Faker) Words(count int, supplemental bool) []string {
	out := make([]string, 0, count)
	for n := 0; n < count; n++ {
		if supplemental && f.Rand.Float32() > 0.5 {
			out = append(out, f.parse("lorem.supplemental"))
		} else {
			out = append(out, f.parse("lorem.words"))
		}
	}
	return out
}

func (f *Faker) Characters(count int) string {
	if count < 1 {
		return ""
	}
	out := make([]rune, 0, count)
	for n := 0; n < count; n++ {
		r := []rune(strconv.FormatInt(f.Rand.Int63n(36), 36))[0]
		out = append(out, r)
	}
	return string(out)
}

func (f *Faker) Sentence(words int, supplemental bool) string {
	if f.Language == "zh-CN" {
		return strings.Join(f.Words(words+f.Rand.Intn(6), supplemental), "") + "。"
	}
	return capitalize(strings.Join(f.Words(words+f.Rand.Intn(6), supplemental), " ")) + "."
}

func (f *Faker) Sentences(count int, supplemental bool) []string {
	out := make([]string, 0, count)
	for n := 0; n < count; n++ {
		out = append(out, f.Sentence(3, supplemental))
	}
	return out
}

func (f *Faker) Paragraph(sentences int, supplemental bool) string {
	return strings.Join(f.Sentences(sentences, supplemental), " ")
}

func (f *Faker) Paragraphs(count int, supplemental bool) []string {
	out := make([]string, 0, count)
	for n := 0; n < count; n++ {
		out = append(out, f.Paragraph(3, supplemental))
	}
	return out
}

func (f *Faker) City() string             { return f.parse("address.city") }
func (f *Faker) StreetName() string       { return f.parse("address.street_name") }
func (f *Faker) StreetAddress() string    { return f.numerify(f.parse("address.street_address")) }
func (f *Faker) SecondaryAddress() string { return f.numerify(f.parse("address.secondary_address")) }
func (f *Faker) PostCode() string         { return f.bothify(f.parse("address.postcode")) }
func (f *Faker) StreetSuffix() string     { return f.parse("address.street_suffix") }
func (f *Faker) CitySuffix() string       { return f.parse("address.city_suffix") }
func (f *Faker) CityPrefix() string       { return f.parse("address.city_prefix") }
func (f *Faker) StateAbbr() string        { return f.parse("address.state_abbr") }
func (f *Faker) State() string            { return f.parse("address.state") }
func (f *Faker) Country() string          { return f.parse("address.country") }
func (f *Faker) Latitude() float64        { return (f.Rand.Float64() * 180) - 90 }
func (f *Faker) Longitude() float64       { return (f.Rand.Float64() * 360) - 180 }

func (f *Faker) CompanyName() string   { return f.parse("company.name") }
func (f *Faker) CompanySuffix() string { return f.parse("company.suffix") }
func (f *Faker) CompanyCatchPhrase() string {
	return f.combine("company.buzzwords.0", "company.buzzwords.1", "company.buzzwords.2")
}
func (f *Faker) CompanyBs() string { return f.combine("company.bs.0", "company.bs.1", "company.bs.2") }

func (f *Faker) PhoneNumber() string {
	return f.numerify(f.combine("phone_number.formats"))
}

func (f *Faker) CellPhoneNumber() string {
	_, got := f.Dict["phone_number.cell_phone"]
	if got {
		return f.numerify(f.combine("phone_number.cell_phone"))
	}
	return f.numerify(f.combine("phone_number.formats"))
}

func (f *Faker) Email() string     { return f.UserName() + "@" + f.DomainName() }
func (f *Faker) FreeEmail() string { return f.UserName() + "@" + f.parse("internet.free_email") }
func (f *Faker) SafeEmail() string { return f.UserName() + "@example." + f.sample("org", "com", "net") }
func (f *Faker) UserName() string {
	return fixUmlauts(strings.ToLower(f.sample(f.FirstName(), f.FirstName()+f.sample(".", "_")+f.LastName())))
}
func (f *Faker) DomainName() string {
	return f.DomainWord() + "." + f.DomainSuffix()
}

func (f *Faker) DomainWord() string {
	return removeFancyChars.ReplaceAllString(
		strings.ToLower(
			strings.SplitN(f.CompanyName(), " ", 2)[0]), "")
}
func (f *Faker) DomainSuffix() string { return f.parse("internet.domain_suffix") }
func (f *Faker) IPv4Address() net.IP {
	oct := func() int { return 2 + f.Rand.Intn(254) }
	ip := fmt.Sprintf("%d.%d.%d.%d", oct(), oct(), oct(), oct())
	return net.ParseIP(ip)
}
func (f *Faker) IPv6Address() net.IP {
	m := 65536
	ip := fmt.Sprintf("2001:cafe:%x:%x:%x:%x:%x:%x",
		f.Rand.Intn(m), f.Rand.Intn(m), f.Rand.Intn(m),
		f.Rand.Intn(m), f.Rand.Intn(m), f.Rand.Intn(m))
	return net.ParseIP(ip)
}
func (f *Faker) URL() string {
	return "http://" + f.DomainName() + "/" + f.UserName()
}

// Name returns a random personal name in various formats.
func (f *Faker) Name() string       { return f.parse("name.name") }
func (f *Faker) FirstName() string  { return f.parse("name.first_name") }
func (f *Faker) LastName() string   { return f.parse("name.last_name") }
func (f *Faker) NamePrefix() string { return f.parse("name.prefix") }
func (f *Faker) NameSuffix() string { return f.parse("name.suffix") }
func (f *Faker) JobTitle() string {
	return f.parse("name.title.descriptor") + " " +
		f.parse("name.title.level") + " " +
		f.parse("name.title.job")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func fixUmlauts(str string) string {
	out := make([]rune, 0, len(str))
	for _, r := range str {
		switch r {
		case 'ä':
			out = append(out, 'a', 'e')
		case 'ö':
			out = append(out, 'o', 'e')
		case 'ü':
			out = append(out, 'u', 'e')
		case 'ß':
			out = append(out, 's', 's')
		default:
			out = append(out, r)
		}
	}

	return string(out)
}

func (f *Faker) combine(keys ...string) string {
	tmp := make([]string, 0, len(keys))
	for _, key := range keys {
		tmp = append(tmp, f.parse(key))
	}
	return strings.Join(tmp, " ")
}

func (f *Faker) parse(key string) string {
	baseKeyIndex := strings.Index(key, ".")
	baseKey := key[0:baseKeyIndex]

	formats, found := f.Dict[key]
	if !found {
		panic("couldn't find key: " + key)
	}

	format := f.sample(formats...)

	return recGsub(keyPattern, format, func(s string) string {
		entryKey := strings.ToLower(s[2 : len(s)-1])

		if strings.Index(entryKey, ".") == -1 {
			entryKey = baseKey + "." + entryKey
		}

		entry, found := f.Dict[entryKey]
		if !found {
			panic("couldn't find entry key: " + entryKey)
		}
		return f.sample(entry...)
	})
}

func recGsub(r *regexp.Regexp, in string, f func(string) string) string {
	for keepRunning := true; keepRunning; {
		keepRunning = false
		in = r.ReplaceAllStringFunc(in, func(s string) string {
			keepRunning = true
			return f(s)
		})
	}

	return in
}

func (f *Faker) sample(set ...string) string {
	idx := f.Rand.Intn(len(set))
	return set[idx]
}

func (f *Faker) bothify(in string) string {
	return f.letterify(f.numerify(in))
}

func (f *Faker) letterify(in string) string {
	return in
}

func (f *Faker) numerify(in string) string {
	return recGsub(numerifyPattern, in, func(s string) string {
		return strings.Map(func(r rune) rune {
			return rune(48 + rand.Intn(9))
		}, s)
	})
}
