package swagger

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func roundTripTest(t *testing.T, fixtureType, fileName string, schema interface{}) {
	specName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	Convey("verifying "+fixtureType+" fixture "+specName, t, func() {
		b, err := ioutil.ReadFile(fileName)
		So(err, ShouldBeNil)
		//Println()
		//Println("Reading file", fileName, "returned", string(b))
		var expected map[string]interface{}
		err = json.Unmarshal(b, &expected)
		So(err, ShouldBeNil)

		err = json.Unmarshal(b, schema)
		So(err, ShouldBeNil)

		//Println()
		//Println("unmarshalling from file resulted in: %#v", schema)
		cb, err := json.MarshalIndent(schema, "", "  ")
		So(err, ShouldBeNil)
		//Println()
		//Println("Marshalling to json returned", string(cb))

		var actual map[string]interface{}
		err = json.Unmarshal(cb, &actual)
		So(err, ShouldBeNil)
		//Println()
		//spew.Dump(expected)
		//spew.Dump(actual)
		//fmt.Printf("comparing %s\n\t%#v\nto\n\t%#+v\n", fileName, expected, actual)
		So(actual, ShouldResemble, expected)
	})
}

func TestPropertyFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "models", "properties")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		roundTripTest(t, "property", filepath.Join(path, f.Name()), &Schema{})
	}
}

func TestModelFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "models")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	specs := []string{"models", "modelWithComposition", "modelWithExamples", "multipleModels"}
FILES:
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		for _, spec := range specs {
			if strings.HasPrefix(f.Name(), spec) {
				roundTripTest(t, "model", filepath.Join(path, f.Name()), &Spec{})
				continue FILES
			}
		}
		//fmt.Println("trying", f.Name())
		roundTripTest(t, "model", filepath.Join(path, f.Name()), &Schema{})
	}
}

func TestParameterFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "resources", "parameters")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		roundTripTest(t, "parameter", filepath.Join(path, f.Name()), &Parameter{})
	}
}

func TestOperationFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "resources", "operations")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		roundTripTest(t, "operation", filepath.Join(path, f.Name()), &Operation{})
	}
}

func TestResponseFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "responses")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "multiple") {
			roundTripTest(t, "response", filepath.Join(path, f.Name()), &Response{})
		} else {
			roundTripTest(t, "responses", filepath.Join(path, f.Name()), &Responses{})
		}
	}
}

//func TestResourcesFixtures(t *testing.T) {
//path := filepath.Join("fixtures", "json", "resources")
//files, err := ioutil.ReadDir(path)
//if err != nil {
//t.Fatal(err)
//}
//for _, f := range files {
//if f.IsDir() {
//continue
//}
//roundTripTest(t, "resources", filepath.Join(path, f.Name()), &Spec{})
//}
//}
