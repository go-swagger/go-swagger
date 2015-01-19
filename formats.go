package swagger

// URI represents the uri string format as specified by the json schema spec
type URI string

func (u URI) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *URI) UnmarshalText(data []byte) error { // validation is performed later on
	*u = URI(string(data))
	return nil
}

// Email represents the email string format as specified by the json schema spec
type Email string

func (e Email) MarshalText() ([]byte, error) {
	return []byte(string(e)), nil
}

func (e *Email) UnmarshalText(data []byte) error { // validation is performed later on
	*e = Email(string(data))
	return nil
}

// Hostname represents the hostname string format as specified by the json schema spec
type Hostname string

func (h Hostname) MarshalText() ([]byte, error) {
	return []byte(string(h)), nil
}

func (h *Hostname) UnmarshalText(data []byte) error { // validation is performed later on
	*h = Hostname(string(data))
	return nil
}

// IPv4 represents an IP v4 address
type IPv4 string

func (u IPv4) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *IPv4) UnmarshalText(data []byte) error { // validation is performed later on
	*u = IPv4(string(data))
	return nil
}

// IPv6 represents an IP v6 address
type IPv6 string

func (u IPv6) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *IPv6) UnmarshalText(data []byte) error { // validation is performed later on
	*u = IPv6(string(data))
	return nil
}

// UUID represents a uuid string format
type UUID string

func (u UUID) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *UUID) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID(string(data))
	return nil
}

// UUID3 represents a uuid3 string format
type UUID3 string

func (u UUID3) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *UUID3) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID3(string(data))
	return nil
}

// UUID4 represents a uuid4 string format
type UUID4 string

func (u UUID4) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *UUID4) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID4(string(data))
	return nil
}

// UUID5 represents a uuid5 string format
type UUID5 string

func (u UUID5) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *UUID5) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID5(string(data))
	return nil
}

// ISBN represents an isbn string format
type ISBN string

func (u ISBN) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *ISBN) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN(string(data))
	return nil
}

// ISBN10 represents an isbn 10 string format
type ISBN10 string

func (u ISBN10) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *ISBN10) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN10(string(data))
	return nil
}

// ISBN13 represents an isbn 13 string format
type ISBN13 string

func (u ISBN13) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *ISBN13) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN13(string(data))
	return nil
}

// CreditCard represents a credit card string format
type CreditCard string

func (u CreditCard) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *CreditCard) UnmarshalText(data []byte) error { // validation is performed later on
	*u = CreditCard(string(data))
	return nil
}

// SSN represents a social security string format
type SSN string

func (u SSN) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

func (u *SSN) UnmarshalText(data []byte) error { // validation is performed later on
	*u = SSN(string(data))
	return nil
}

// HexColor represents a hex color string format
type HexColor string

func (h HexColor) MarshalText() ([]byte, error) {
	return []byte(string(h)), nil
}

func (h *HexColor) UnmarshalText(data []byte) error { // validation is performed later on
	*h = HexColor(string(data))
	return nil
}

// RGBColor represents a RGB color string format
type RGBColor string

func (r RGBColor) MarshalText() ([]byte, error) {
	return []byte(string(r)), nil
}

func (r *RGBColor) UnmarshalText(data []byte) error { // validation is performed later on
	*r = RGBColor(string(data))
	return nil
}
