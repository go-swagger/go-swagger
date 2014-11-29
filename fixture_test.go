package swagger

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func propertyTest(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	So(err, ShouldBeNil)
	var expected map[string]interface{}
	err = json.Unmarshal(b, &expected)
	So(err, ShouldBeNil)
	schema := Schema{}
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
	Convey("the property fixtures should round trip", t, func() {
		path := filepath.Join("fixtures", "json", "models", "properties")
		files, err := ioutil.ReadDir(path)
		So(err, ShouldBeNil)

		for _, f := range files {
			propertyTest(filepath.Join(path, f.Name()))
		}
	})
}
