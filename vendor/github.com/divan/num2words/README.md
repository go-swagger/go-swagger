num2words
=========

[![Build Status](https://drone.io/github.com/divan/num2words/status.png)](https://drone.io/github.com/divan/num2words/latest)
[![GoDoc](https://godoc.org/github.com/divan/num2words?status.svg)](https://godoc.org/github.com/divan/num2words)

num2words - Numbers to words converter in Go (Golang)

## Usage

First, import package num2words

```import github.com/divan/num2words```

Convert number
```go
  str := num2words.Convert(17) // outputs "seventeen"
  ...
  str := num2words.Convert(1024) // outputs "one thousand twenty four"
  ...
  str := num2words.Convert(-123) // outputs "minus one hundred twenty three"
```
