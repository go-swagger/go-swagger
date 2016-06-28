/*
Copyright 2014 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package gce configures hooks for running Camlistore for Google Compute Engine.
package gce // import "camlistore.org/pkg/osutil/gce"

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"camlistore.org/pkg/env"
	"camlistore.org/pkg/osutil"

	"go4.org/jsonconfig"
	_ "go4.org/wkfs/gcs"
	"golang.org/x/net/context"
	"google.golang.org/cloud/compute/metadata"
	"google.golang.org/cloud/logging"
)

func init() {
	if !env.OnGCE() {
		return
	}
	osutil.RegisterConfigDirFunc(func() string {
		v, _ := metadata.InstanceAttributeValue("camlistore-config-dir")
		if v == "" {
			return v
		}
		return path.Clean("/gcs/" + strings.TrimPrefix(v, "gs://"))
	})
	jsonconfig.RegisterFunc("_gce_instance_meta", func(c *jsonconfig.ConfigParser, v []interface{}) (interface{}, error) {
		if len(v) != 1 {
			return nil, errors.New("only 1 argument supported after _gce_instance_meta")
		}
		attr, ok := v[0].(string)
		if !ok {
			return nil, errors.New("expected argument after _gce_instance_meta to be a string")
		}
		val, err := metadata.InstanceAttributeValue(attr)
		if err != nil {
			return nil, fmt.Errorf("error reading GCE instance attribute %q: %v", attr, err)
		}
		return val, nil
	})
}

// LogWriter returns an environment-specific io.Writer suitable for passing
// to log.SetOutput. It will also include writing to os.Stderr as well.
func LogWriter() (w io.Writer) {
	w = os.Stderr
	if !env.OnGCE() {
		return
	}
	projID, err := metadata.ProjectID()
	if projID == "" {
		log.Printf("Error getting project ID: %v", err)
		return
	}
	scopes, _ := metadata.Scopes("default")
	haveScope := func(scope string) bool {
		for _, x := range scopes {
			if x == scope {
				return true
			}
		}
		return false
	}
	if !haveScope(logging.Scope) {
		log.Printf("when this Google Compute Engine VM instance was created, it wasn't granted enough access to use Google Cloud Logging (Scope URL: %v).", logging.Scope)
		return
	}

	logc, err := logging.NewClient(context.Background(), projID, "camlistored-stderr")
	if err != nil {
		log.Printf("Error creating Google logging client: %v", err)
		return
	}
	return io.MultiWriter(w, logc.Writer(logging.Debug))
}
