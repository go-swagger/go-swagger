---
title: strfmt
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:strfmt

A **swagger:strfmt** annotation names a type as a string formatter. The name is mandatory and that is
what will be used as format name for this particular string format.
String formats should only be used for **very** well known formats.

<!-- more -->

String formats are well-known items. These imply a common well-documented set of formats that can be validated. The toolkit allows for creating your own string formats too.

To create a custom string format you need to create a type that implements the (Unm/M)arshalText interfaces and the sql Scan and sql Value interfaces.  The SQL interfaces are not strictly necessary but allow other people to use the string format in structs that are used with databases

The default string formats for this toolkit are:

* uuid, uuid3, uuid4, uuid5
* email
* uri (absolute)
* hostname
* ipv4
* ipv6
* credit card
* isbn, isbn10, isbn13
* social security number
* hexcolor
* rgbcolor
* date
* date-time
* duration
* password
* custom string formats

##### Syntax

```go
swagger:strfmt [name]
```

##### Example

```go
func init() {
  eml := Email("")
  Default.Add("email", &eml, govalidator.IsEmail)
}

// Email represents the email string format as specified by the json schema spec
//
// swagger:strfmt email
type Email string

// MarshalText turns this instance into text
func (e Email) MarshalText() ([]byte, error) {
	return []byte(string(e)), nil
}

// UnmarshalText hydrates this instance from text
func (e *Email) UnmarshalText(data []byte) error { // validation is performed later on
	*e = Email(string(data))
	return nil
}

func (b *Email) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case []byte:
		*b = Email(string(v))
	case string:
		*b = Email(v)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Email from: %#v", v)
	}

	return nil
}

func (b Email) Value() (driver.Value, error) {
	return driver.Value(string(b)), nil
}
```

##### Result

```yaml
```
