package swagger

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func roundTripTest(fileName string, schema interface{}) {
	b, err := ioutil.ReadFile(fileName)
	So(err, ShouldBeNil)
	var expected map[string]interface{}
	err = json.Unmarshal(b, &expected)
	So(err, ShouldBeNil)
	err = json.Unmarshal(b, &schema)
	So(err, ShouldBeNil)
	cb, err := json.Marshal(schema)
	So(err, ShouldBeNil)
	var actual map[string]interface{}
	err = json.Unmarshal(cb, &actual)
	So(err, ShouldBeNil)
	So(expected, ShouldResemble, actual)
}

func TestPropertyFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "models", "properties")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}

	Convey("the property fixtures should round trip", t, func() {
		for _, f := range files {
			Convey("for "+strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())), func() {
				roundTripTest(filepath.Join(path, f.Name()), Schema{})
			})
		}
	})
}
func TestModelFixtures(t *testing.T) {
	path := filepath.Join("fixtures", "json", "models")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Error(err)
	}
	var filepaths []os.FileInfo
	for _, f := range files {
		if !f.IsDir() {
			filepaths = append(filepaths, f)
		}
	}

	Convey("the spec should round trip for models", t, func() {

		for _, f := range filepaths {
			Convey("for "+strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())), func() {
				roundTripTest(filepath.Join(path, f.Name()), Spec{})
			})
		}
	})
}
