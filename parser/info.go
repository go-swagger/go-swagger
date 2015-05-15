package parser

import (
	"go/ast"
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

type infoSectionTagger struct {
	*sectionTagger
	set infoSetter
}

func newInfoSection(name string, multiLine bool, setter infoSetter) (t *infoSectionTagger) {
	t = new(infoSectionTagger)
	t.sectionTagger = newSectionTagger(name, multiLine)
	t.set = setter
	return
}

func newAPIInfoParser() *apiInfoParser {
	taggers := []*infoSectionTagger{
		newInfoSection("Version", false, setInfoVersion),
		newInfoSection("Description", true, setInfoDescription),
		newInfoSection("Title", true, setInfoTitle),
		newInfoSection("TOS", true, setInfoTOS),
		newInfoSection("License", false, setInfoLicense),
		newInfoSection("Contact", false, setInfoContact),
	}
	return &apiInfoParser{taggers}
}

type apiInfoParser struct {
	taggers []*infoSectionTagger
}

func (ai *apiInfoParser) Parse(gofile *ast.File) (*spec.Info, error) {
	info := new(spec.Info)

	// var currentLines []string
	var selectedTagger *infoSectionTagger
	var otherTags []string
	taggers := ai.taggers
	for _, c := range gofile.Doc.List {
		text := c.Text
		lines := strings.Split(text, "\n")
	LINES:
		for _, line := range lines {
			// this is an aggregating tagger
			if selectedTagger != nil {
				switch res := selectedTagger.Tag(line, otherTags).(type) {
				case multiLineSectionPart:
					continue LINES
				case multiLineSectionTerminator:
					if err := selectedTagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
					selectedTagger = nil
					continue LINES
				case newTagSectionTerminator:
					if err := selectedTagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
				}
			}
			selectedTagger = nil
			for i, tagger := range taggers {
				switch res := tagger.Tag(line, nil).(type) {
				case singleLineSection:
					if err := tagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case multiLineSectionPart:
					selectedTagger = tagger
					otherTags = nil
					for _, t := range ai.taggers {
						if t.Name != tagger.Name {
							otherTags = append(otherTags, t.Name)
						}
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES
				case unmatchedSection:
					// TODO: something slightly smarter than nothing
				}
			}
		}
	}

	if selectedTagger != nil {
		if err := selectedTagger.set(info, selectedTagger.Lines); err != nil {
			return nil, err
		}
	}

	return info, nil
}

func setInfoVersion(info *spec.Info, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	info.Version = lines[0]
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
	info.Title = strings.Join(lines, "\n")
	return nil
}

func setInfoDescription(info *spec.Info, lines []string) error {
	info.Description = strings.Join(lines, "\n")
	return nil
}

func setInfoTOS(info *spec.Info, lines []string) error {
	info.TermsOfService = strings.Join(lines, "\n")
	return nil
}
