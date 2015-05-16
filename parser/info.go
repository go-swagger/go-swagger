package parser

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger/spec"
)

var infoTags = []string{
	"Description",
	"Title",
	"TOS",
	"Version",
	"License",
	"Contact",
}

type infoSetter func(*spec.Info, []string) error

func newInfoSection(name string, multiLine bool, setter infoSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Info), lines) }
	return
}

func newAPIInfoParser(otherTags []string) *docCommentParser {
	return newDocCommentParser(
		otherTags,
		newInfoSection("Version", false, setInfoVersion),
		newInfoSection("Description", true, setInfoDescription),
		newInfoSection("Title", true, setInfoTitle),
		newInfoSection("TOS", true, setInfoTOS),
		newInfoSection("License", false, setInfoLicense),
		newInfoSection("Contact", false, setInfoContact),
	)
}

func setInfoVersion(info *spec.Info, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	info.Version = strings.TrimSpace(lines[0])
	return nil
}

func setInfoContact(info *spec.Info, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	contact, err := parseContactInfo(lines[0])
	if err != nil {
		return err
	}
	info.Contact = contact
	return nil
}

func parseContactInfo(line string) (*spec.ContactInfo, error) {
	nameEmail, url := splitURL(line)
	var name, email string
	if len(nameEmail) > 0 {
		addr, err := mail.ParseAddress(nameEmail)
		if err != nil {
			return nil, err
		}
		name, email = addr.Name, addr.Address
	}
	return &spec.ContactInfo{
		URL:   url,
		Name:  name,
		Email: email,
	}, nil
}

func setInfoLicense(info *spec.Info, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	line := lines[0]
	name, url := splitURL(line)
	info.License = &spec.License{
		Name: name,
		URL:  url,
	}
	return nil
}

// httpFTPScheme matches http://, https://, ftp:// or ftps://
var httpFTPScheme = regexp.MustCompile("(?:ht|f)tp(?:s)?://")

func splitURL(line string) (notURL, url string) {
	str := strings.TrimSpace(line)
	parts := httpFTPScheme.FindStringIndex(str)
	if len(parts) == 0 {
		if len(str) > 0 {
			notURL = str
		}
		return
	}
	if len(parts) > 0 {
		notURL = strings.TrimSpace(str[:parts[0]])
		url = strings.TrimSpace(str[parts[0]:])
	}
	return
}

func setInfoTitle(info *spec.Info, lines []string) error {
	info.Title = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}

func setInfoDescription(info *spec.Info, lines []string) error {
	info.Description = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}

func setInfoTOS(info *spec.Info, lines []string) error {
	info.TermsOfService = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}
