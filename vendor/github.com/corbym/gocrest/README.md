# gocrest

A hamcrest-like assertion library for Go. GoCrest matchers are composable, self-describing and
can be strung together in a more readable form to create flexible assertions. 

Inspired by [Hamcrest](https://github.com/hamcrest). 

[![Build status](https://travis-ci.org/corbym/gocrest.svg?branch=master)](https://github.com/corbym/gocrest)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbym/gocrest)](https://goreportcard.com/report/github.com/corbym/gocrest)
[![GoDoc](https://godoc.org/github.com/corbym/gocrest?status.svg)](http://godoc.org/github.com/corbym/gocrest)
[![BCH compliance](https://bettercodehub.com/edge/badge/corbym/gocrest?branch=master)](https://bettercodehub.com/)
[![Coverage Status](https://coveralls.io/repos/github/corbym/gocrest/badge.svg?branch=master)](https://coveralls.io/github/corbym/gocrest?branch=master)
## Package import

```
import (
  "github.com/corbym/gocrest/then"
  "github.com/corbym/gocrest/is"
  "github.com/corbym/gocrest/has"
)
```

## Example:
```go
then.AssertThat(testing, "hi", is.EqualTo("bye").Reason("we are going"))
```

output:

```
we are going
Expected: value equal to <bye>
     but: <hi>
```

Composed with AllOf:

```go
then.AssertThat(t, "abcdef", is.AllOf(is.ValueContaining("abc"), is.LessThan("ghi")))
```

# Matchers so far..

- is.EqualTo(x)
- is.EqualToIgnoringWhitespace(string) - compares two strings without comparing their whitespace characters.
- is.Nil() - value must be nil
- is.ValueContaining(expected) -- acts like containsAll
- is.Not(m *Matcher) -- logical not of matcher's result
- is.MatchForPattern(regex string) -- a string regex expression
- has.FunctionNamed(string x) - checks if an interface has a function (method)
- has.FieldNamed(string x) - checks if a struct has a field named x
- is.AllOf(... *Matcher) - returns true if all matchers match
- is.AnyOf(... *Matcher) - return true if any matcher matches
- is.GreaterThan(expected) - checks if actual > expected
- is.LessThan(expected)
- is.Empty() - matches if the actual is "", nil or len(actual)==0
- is.LessThan(x)
- is.LessThanOrEqualTo(x)
- is.GreaterThan(x)
- is.GreaterThanOrEqualTo(x)
- has.Length(x) - matcher if given value (int or matcher) matches the len of the given
- has.Prefix(x) - string starts with x
- has.Suffix(x) - string ends with x
- has.Key(x) - map has key x
- has.AllKeys(T x, T y) (or has.AllKeys([]T{x,y})) - finds key of type T in map

For more comprehensive documentation see [godoc](http://godoc.org/github.com/corbym/gocrest).
