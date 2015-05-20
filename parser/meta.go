package parser

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger/spec"
)

var allSwaggerTags = []string{
	"Consumes",
	"Produces",
	"Schemes",
	"Host",
	"BasePath",
	"Tags",
	"TOS",
	"Version",
	"License",
	"Contact",
}

type metaSetter func(*spec.Swagger, []string) error

func newMetaTitle(setter metaSetter) (t *sectionTagger) {
	t = newTitleTagger()
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Swagger), lines) }
	return
}
func newMetaDescription(setter metaSetter) (t *sectionTagger) {
	t = newDescriptionTagger()
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Swagger), lines) }
	return
}
func newMetaSection(name string, multiLine bool, setter metaSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Swagger), lines) }
	return
}

func newMetaParser() *docCommentParser {
	return newDocCommentParser(
		nil,
		newMetaTitle(setInfoTitle),
		newMetaDescription(setInfoDescription),
		newMetaSection("Version", false, setInfoVersion),
		newMetaSection("TOS", true, setInfoTOS),
		newMetaSection("License", false, setInfoLicense),
		newMetaSection("Contact", false, setInfoContact),
		newMetaSection("Consumes", true, setSwaggerConsumes),
		newMetaSection("Produces", true, setSwaggerProduces),
		newMetaSection("Schemes", false, setSwaggerSchemes),
		newMetaSection("Host", false, setSwaggerHost),
		newMetaSection("BasePath", false, setSwaggerBasePath),
		//newMetaSection("Tags", false, setSwaggerBasePath),
	)
}

func setSwaggerConsumes(swspec *spec.Swagger, lines []string) error {
	swspec.Consumes = lines
	return nil
}

func setSwaggerProduces(swspec *spec.Swagger, lines []string) error {
	swspec.Produces = lines
	return nil
}

func setSwaggerSchemes(swspec *spec.Swagger, lines []string) error {
	lns := lines
	if len(lns) == 0 || lns[0] == "" {
		lns = []string{"http"}
	}
	swspec.Schemes = strings.Split(lns[0], ",")
	return nil
}

func setSwaggerHost(swspec *spec.Swagger, lines []string) error {
	lns := lines
	if len(lns) == 0 {
		lns = []string{"localhost"}
	}
	swspec.Host = lns[0]
	return nil
}

func setSwaggerBasePath(swspec *spec.Swagger, lines []string) error {
	lns := lines
	if len(lns) == 0 {
		lns = []string{"/"}
	}
	swspec.BasePath = lns[0]
	return nil
}

func setInfoVersion(swspec *spec.Swagger, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	info := safeInfo(swspec)
	info.Version = strings.TrimSpace(lines[0])
	return nil
}

func setInfoContact(swspec *spec.Swagger, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	contact, err := parseContactInfo(lines[0])
	if err != nil {
		return err
	}
	info := safeInfo(swspec)
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

func setInfoLicense(swspec *spec.Swagger, lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	info := safeInfo(swspec)
	line := lines[0]
	name, url := splitURL(line)
	info.License = &spec.License{
		Name: name,
		URL:  url,
	}
	return nil
}

func safeInfo(swspec *spec.Swagger) *spec.Info {
	if swspec.Info == nil {
		swspec.Info = new(spec.Info)
	}
	return swspec.Info
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

func setInfoTitle(swspec *spec.Swagger, lines []string) error {
	info := safeInfo(swspec)
	info.Title = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}

func setInfoDescription(swspec *spec.Swagger, lines []string) error {
	info := safeInfo(swspec)
	info.Description = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}

func setInfoTOS(swspec *spec.Swagger, lines []string) error {
	info := safeInfo(swspec)
	info.TermsOfService = strings.TrimSpace(strings.Join(lines, "\n"))
	return nil
}
