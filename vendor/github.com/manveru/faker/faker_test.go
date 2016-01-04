package faker

import (
	"fmt"
	. "github.com/manveru/gobdd"
	"math/rand"
	"testing"
)

func ExampleFaker_Name() {
	fake, _ := New("en")
	fake.Rand = rand.New(rand.NewSource(42))
	fmt.Println(fake.Name())
	// Output: Adriana Crona
}

func TestEverything(*testing.T) {}

func init() {
	defer PrintSpecReport()

	Describe("english faker", func() {
		fake, _ := New("en")
		fake.Rand = rand.New(rand.NewSource(42))

		BeforeEach(func() {
			fake.Rand = rand.New(rand.NewSource(42))
		})

		It("makes fake city names", func() {
			Expect(func() string { return fake.City() },
				Returns, "West Kim", "Jackyland", "New Nathen")
		})
		It("makes fake street names", func() {
			Expect(func() string { return fake.StreetName() },
				Returns, "Willms Grove", "Jacky Harbor", "Lesch Parkway")
		})
		It("makes fake street addresses", func() {
			Expect(func() string { return fake.StreetAddress() }, Returns,
				"562 Dustin Prairie", "24678 Nathen Extension", "465 Altenwerth Hollow")
		})
		It("makes fake post codes", func() {
			Expect(func() string { return fake.PostCode() },
				Returns, "73248-0225", "20510-7716", "34655")
		})
		It("makes fake street suffixes", func() {
			Expect(func() string { return fake.StreetSuffix() },
				Returns, "Island", "Gardens", "Grove")
		})
		It("makes fake city suffixes", func() {
			Expect(func() string { return fake.CitySuffix() },
				Returns, "burgh", "ville", "mouth")
		})
		It("makes fake city prefixes", func() {
			Expect(func() string { return fake.CityPrefix() },
				Returns, "Lake", "West", "New")
		})
		It("makes random US state abbreviations", func() {
			Expect(func() string { return fake.StateAbbr() },
				Returns, "AZ", "FL", "VA")
		})
		It("makes random US state names", func() {
			Expect(func() string { return fake.State() },
				Returns, "Colorado", "Pennsylvania", "Maine")
		})
		It("makes random country names", func() {
			Expect(func() string { return fake.Country() },
				Returns, "Tajikistan", "Cameroon", "Cote d'Ivoire")
		})
		It("makes random latitude", func() {
			Expect(fake.Latitude(), ToEqual, -22.854895011606132)
			Expect(fake.Latitude(), ToEqual, -78.11991057716678)
			Expect(fake.Latitude(), ToEqual, 18.736893280555563)
		})
		It("makes random longitude", func() {
			Expect(fake.Longitude(), ToEqual, -45.709790023212264)
			Expect(fake.Longitude(), ToEqual, -156.23982115433355)
			Expect(fake.Longitude(), ToEqual, 37.473786561111126)
		})

		It("makes fake company names", func() {
			Expect(func() string { return fake.CompanyName() }, Returns,
				"Willms, Crona and Swift",
				"Schaefer-Maggio",
				"Hartmann, Huels and Wolff")
		})
		It("makes fake company suffixes", func() {
			Expect(func() string { return fake.CompanySuffix() },
				Returns, "and Sons", "Group", "Inc")
		})
		It("makes company catch phrases", func() {
			Expect(func() string { return fake.CompanyCatchPhrase() }, Returns,
				"Balanced next generation circuit",
				"Object-based high-level task-force",
				"Organized actuating intranet")
		})
		It("makes company bs", func() {
			Expect(func() string { return fake.CompanyBs() }, Returns,
				"deliver dynamic e-markets",
				"iterate impactful e-services",
				"utilize robust eyeballs")
		})

		It("makes names", func() {
			Expect(func() string { return fake.Name() }, Returns,
				"Adriana Crona", "Miss Stuart Maggio", "Nathen Huels")
		})
		It("makes first names", func() {
			Expect(func() string { return fake.FirstName() },
				Returns, "Jeromy", "Adriana", "Kim")
		})
		It("makes last names", func() {
			Expect(func() string { return fake.LastName() }, Returns,
				"Willms", "Willms", "Crona")
		})
		It("makes name prefixes", func() {
			Expect(func() string { return fake.NamePrefix() }, Returns,
				"Mr.", "Ms.", "Miss")
		})
		It("makes name suffixes", func() {
			Expect(func() string { return fake.NameSuffix() }, Returns,
				"Jr.", "PhD", "Sr.")
		})
		It("makes job titles", func() {
			Expect(func() string { return fake.JobTitle() }, Returns,
				"Human Factors Planner",
				"Legacy Mobility Representative",
				"Senior Accounts Architect")
		})

		It("makes ", func() {
			Expect(func() string { return fake.Email() }, Returns,
				"adriana.swift@maggiolesch.info",
				"bernadette@rau.org",
				"charity.brown@fritschbotsford.biz")
		})
		It("makes ", func() {
			Expect(func() string { return fake.FreeEmail() }, Returns,
				"adriana.swift@yahoo.com",
				"hertha.huels@hotmail.com",
				"susana@hotmail.com")
		})
		It("makes ", func() {
			Expect(func() string { return fake.SafeEmail() }, Returns,
				"adriana.swift@example.com",
				"hertha.huels@example.net",
				"susana@example.net")
		})
		It("makes ", func() {
			Expect(func() string { return fake.IPv4Address().String() },
				Returns, "249.195.38.38", "45.21.213.246", "82.77.201.235")
		})
		It("makes ", func() {
			Expect(func() string { return fake.IPv6Address().String() }, Returns,
				"2001:cafe:64b1:b44b:f284:923e:d7df:7a61",
				"2001:cafe:5ba5:8588:9b70:75d3:19f9:826f",
				"2001:cafe:9ee7:badc:4f8c:8f6d:ce04:cf79")
		})
		It("makes ", func() {
			Expect(func() string { return fake.URL() }, Returns,
				"http://willms.biz/keely.hartmann",
				"http://lueilwitzboehm.net/marilou.rodriguez",
				"http://jerde.name/abdul.jaskolski")
		})

		It("makes random words", func() {
			Expect(fake.Words(0, false), ToDeepEqual, []string{})
			Expect(fake.Words(1, false), ToDeepEqual, []string{"at"})
			Expect(fake.Words(2, false), ToDeepEqual, []string{"illum", "ut"})
			Expect(fake.Words(3, false), ToDeepEqual, []string{"est", "sit", "soluta"})
			Expect(fake.Words(0, true), ToDeepEqual, []string{})
			Expect(fake.Words(1, true), ToDeepEqual, []string{"stabilis"})
			Expect(fake.Words(2, true), ToDeepEqual, []string{"sunt", "theologus"})
			Expect(fake.Words(3, true), ToDeepEqual, []string{"facere", "super", "adipiscor"})
		})

		It("makes random characters", func() {
			Expect(fake.Characters(0), ToEqual, "")
			Expect(fake.Characters(50), ToEqual, "v7ottxf08ku5u1e8barnpqxfqhposxhm1ur5wc7jqy6ccjw8bx")
			for n := 0; n < 100; n++ {
				Expect(len(fake.Characters(n)), ToEqual, n)
			}
		})

		It("makes random sentences", func() {
			Expect(fake.Sentence(0, false), ToEqual, "Illum ut est sit soluta.")
			Expect(fake.Sentence(15, false), ToEqual, "Numquam nobis sunt quaerat ea dolores facere deleniti culpa numquam ut distinctio maxime consequatur est qui corporis sunt.")
			Expect(fake.Sentence(0, true), ToEqual, "Agnosco odit voluptas sumo ipsa.")
			Expect(fake.Sentence(15, true), ToEqual, "Tempus solutio umbra hic vulnero baiulus colo blanditiis circumvenio nostrum eius fugit cogo centum fuga.")
		})

		It("makes random paragraphs", func() {
			Expect(fake.Paragraph(0, false), ToEqual, "")
			Expect(fake.Paragraph(2, false), ToEqual, "Illum ut est sit soluta nulla numquam nobis. Quaerat ea dolores facere.")
			Expect(fake.Paragraph(4, false), ToEqual, "Culpa numquam ut distinctio maxime. Est qui corporis sunt officia odit et. Odit molestias voluptas porro. Magnam ipsa corporis.")
			Expect(fake.Paragraph(0, true), ToEqual, "")
			Expect(fake.Paragraph(2, true), ToEqual, "Non averto quisquam corpus. Baiulus colo blanditiis.")
			Expect(fake.Paragraph(4, true), ToEqual, "Tersus qui suscipit tenus et quod. Comprehendo coepi terminatio claudeo suscipio. Voluptas bis voluptatibus voluptatibus sol. Terebro arto autem canonicus stabilis defungo adnuo at.")
		})

		It("makes random US phone numbers", func() {
			Expect(func() string { return fake.PhoneNumber() }, Returns,
				"(386)730-7410", "107.113.0706", "478-375-8633 x188")
		})
	})

	Describe("german faker", func() {
		fake, _ := New("de")
		fake.Rand = rand.New(rand.NewSource(42))

		BeforeEach(func() {
			fake.Rand = rand.New(rand.NewSource(42))
		})

		It("makes fake city names", func() {
			Expect(func() string { return fake.City() },
				Returns, "West Collien", "Neostadt", "Neu Markus")
		})
		It("makes fake street names", func() {
			Expect(func() string { return fake.StreetName() },
				Returns, "Fichtestr.", "Eichenkamp", "Heinrich-Böll-Str.")
		})
		It("makes fake street addresses", func() {
			Expect(func() string { return fake.StreetAddress() }, Returns,
				"Eichenkamp 85a", "Große Kirchstr. 21c", "Händelstr. 0856")
		})
		It("makes fake post codes", func() {
			Expect(func() string { return fake.PostCode() },
				Returns, "01478", "60712", "72581")
		})
		It("makes fake street suffixes", func() {
			Expect(func() string { return fake.StreetSuffix() },
				Returns, "", "", "")
		})
		It("makes fake city suffixes", func() {
			Expect(func() string { return fake.CitySuffix() },
				Returns, "stadt", "land", "scheid")
		})
		It("makes fake city prefixes", func() {
			Expect(func() string { return fake.CityPrefix() },
				Returns, "Alt", "West", "Neu")
		})
		It("makes random German state abbreviations", func() {
			Expect(func() string { return fake.StateAbbr() },
				Returns, "BE", "SH", "HB")
		})
		It("makes random US state names", func() {
			Expect(func() string { return fake.State() },
				Returns, "", "", "")
		})
		It("makes random country names", func() {
			Expect(func() string { return fake.Country() },
				Returns, "Turks- und Caicosinseln", "Turks- und Caicosinseln", "Guinea")
		})
		It("makes random latitude", func() {
			Expect(fake.Latitude(), ToEqual, -22.854895011606132)
			Expect(fake.Latitude(), ToEqual, -78.11991057716678)
			Expect(fake.Latitude(), ToEqual, 18.736893280555563)
		})
		It("makes random longitude", func() {
			Expect(fake.Longitude(), ToEqual, -45.709790023212264)
			Expect(fake.Longitude(), ToEqual, -156.23982115433355)
			Expect(fake.Longitude(), ToEqual, 37.473786561111126)
		})

		It("makes fake company names", func() {
			Expect(func() string { return fake.CompanyName() }, Returns,
				"Lang, Schindzielorz und Bühler",
				"Döring-Spank",
				"Eberhard, Wimmer und Thust")
		})
		It("makes fake company suffixes", func() {
			Expect(func() string { return fake.CompanySuffix() },
				Returns, "Gruppe", "Gruppe", "Gruppe")
		})
		It("makes company catch phrases", func() {
			Expect(func() string { return fake.CompanyCatchPhrase() }, Returns,
				"Balanced next generation circuit",
				"Object-based high-level task-force",
				"Organized actuating intranet")
		})
		It("makes company bs", func() {
			Expect(func() string { return fake.CompanyBs() }, Returns,
				"deliver dynamic e-markets",
				"iterate impactful e-services",
				"utilize robust eyeballs")
		})

		It("makes names", func() {
			Expect(func() string { return fake.Name() }, Returns,
				"Adriana Crona", "Miss Stuart Maggio", "Nathen Huels")
		})
		It("makes first names", func() {
			Expect(func() string { return fake.FirstName() },
				Returns, "Dorian", "Leroy", "Collien")
		})
		It("makes last names", func() {
			Expect(func() string { return fake.LastName() }, Returns,
				"Urbansky", "Lang", "Schindzielorz")
		})
		It("makes name prefixes", func() {
			Expect(func() string { return fake.NamePrefix() }, Returns,
				"Fr.", "Prof. Dr.", "Hr.")
		})
		It("makes name suffixes", func() {
			Expect(func() string { return fake.NameSuffix() }, Returns,
				"von der", "von der", "von der")
		})
		It("makes job titles", func() {
			Expect(func() string { return fake.JobTitle() }, Returns,
				"Human Factors Planner",
				"Legacy Mobility Representative",
				"Senior Accounts Architect")
		})

		It("makes all kinds of email addresses", func() {
			Expect(func() string { return fake.Email() }, Returns,
				"leroy.buehler@spankrittweg.org",
				"hauke@stenzel.net",
				"lisanne.lauckner@ullmannpolizzi.name")
		})
		It("makes free email addresses", func() {
			Expect(func() string { return fake.FreeEmail() }, Returns,
				"leroy.buehler@yahoo.com",
				"aurora.wimmer@hotmail.com",
				"luk@hotmail.com")
		})
		It("makes safe email addreses", func() {
			Expect(func() string { return fake.SafeEmail() }, Returns,
				"leroy.buehler@example.com",
				"aurora.wimmer@example.net",
				"luk@example.net")
		})
		It("makes ", func() {
			Expect(func() string { return fake.IPv4Address().String() },
				Returns, "249.195.38.38", "45.21.213.246", "82.77.201.235")
		})
		It("makes ", func() {
			Expect(func() string { return fake.IPv6Address().String() }, Returns,
				"2001:cafe:64b1:b44b:f284:923e:d7df:7a61",
				"2001:cafe:5ba5:8588:9b70:75d3:19f9:826f",
				"2001:cafe:9ee7:badc:4f8c:8f6d:ce04:cf79")
		})
		It("makes ", func() {
			Expect(func() string { return fake.URL() }, Returns,
				"http://lang.de/sandy.eberhard",
				"http://kastenlewke.org/josua.rossberg",
				"http://ruckdeschel.com/kaan.herweg")
		})

		It("makes random words", func() {
			Expect(fake.Words(0, false), ToDeepEqual, []string{})
			Expect(fake.Words(1, false), ToDeepEqual, []string{"at"})
			Expect(fake.Words(2, false), ToDeepEqual, []string{"illum", "ut"})
			Expect(fake.Words(3, false), ToDeepEqual, []string{"est", "sit", "soluta"})
			Expect(fake.Words(0, true), ToDeepEqual, []string{})
			Expect(fake.Words(1, true), ToDeepEqual, []string{"stabilis"})
			Expect(fake.Words(2, true), ToDeepEqual, []string{"sunt", "theologus"})
			Expect(fake.Words(3, true), ToDeepEqual, []string{"facere", "super", "adipiscor"})
		})

		It("makes random characters", func() {
			Expect(fake.Characters(0), ToEqual, "")
			Expect(fake.Characters(50), ToEqual, "v7ottxf08ku5u1e8barnpqxfqhposxhm1ur5wc7jqy6ccjw8bx")
			for n := 0; n < 100; n++ {
				Expect(len(fake.Characters(n)), ToEqual, n)
			}
		})

		It("makes random sentences", func() {
			Expect(fake.Sentence(0, false), ToEqual, "Illum ut est sit soluta.")
			Expect(fake.Sentence(15, false), ToEqual, "Numquam nobis sunt quaerat ea dolores facere deleniti culpa numquam ut distinctio maxime consequatur est qui corporis sunt.")
			Expect(fake.Sentence(0, true), ToEqual, "Agnosco odit voluptas sumo ipsa.")
			Expect(fake.Sentence(15, true), ToEqual, "Tempus solutio umbra hic vulnero baiulus colo blanditiis circumvenio nostrum eius fugit cogo centum fuga.")
		})

		It("makes random paragraphs", func() {
			Expect(fake.Paragraph(0, false), ToEqual, "")
			Expect(fake.Paragraph(2, false), ToEqual, "Illum ut est sit soluta nulla numquam nobis. Quaerat ea dolores facere.")
			Expect(fake.Paragraph(4, false), ToEqual, "Culpa numquam ut distinctio maxime. Est qui corporis sunt officia odit et. Odit molestias voluptas porro. Magnam ipsa corporis.")
			Expect(fake.Paragraph(0, true), ToEqual, "")
			Expect(fake.Paragraph(2, true), ToEqual, "Non averto quisquam corpus. Baiulus colo blanditiis.")
			Expect(fake.Paragraph(4, true), ToEqual, "Tersus qui suscipit tenus et quod. Comprehendo coepi terminatio claudeo suscipio. Voluptas bis voluptatibus voluptatibus sol. Terebro arto autem canonicus stabilis defungo adnuo at.")
		})

		It("makes random German phone numbers", func() {
			Expect(func() string { return fake.PhoneNumber() }, Returns,
				"(02461) 0430723", "+49-4385-21431614", "(0804) 543021034")
		})
	})
}

func Returns(fun func() string, expected ...string) (string, bool) {
	for _, s := range expected {
		value := fun()
		if value != s {
			return fmt.Sprintf("\texpected: %#v\n\t   to be: %#v\n", s, value), false
		}
	}
	return "", true
}
