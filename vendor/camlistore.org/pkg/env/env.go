/*
Copyright 2015 The Camlistore Authors

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

// Package env detects what sort of environment Camlistore is running in.
package env // import "camlistore.org/pkg/env"

import (
	"os"
	"strconv"
	"sync"

	"google.golang.org/cloud/compute/metadata"
)

// IsDebug reports whether this is a debug environment.
func IsDebug() bool {
	return isDebug
}

// DebugUploads reports whether this is a debug environment for uploads.
func DebugUploads() bool {
	return isDebugUploads
}

// IsDev reports whether this is a development server environment (devcam server).
func IsDev() bool {
	return isDev
}

// OsGCE reports whether this process is running in a Google Compute
// Engine (GCE) environment.  This only returns true if the
// "camlistore-config-dir" instance metadata value is defined.
// Instances running in custom configs on GCE will be unaffected.
func OnGCE() bool {
	gceOnce.Do(detectGCE)
	return isGCE
}

var (
	gceOnce sync.Once
	isGCE   bool
)

func detectGCE() {
	if !metadata.OnGCE() {
		return
	}
	v, _ := metadata.InstanceAttributeValue("camlistore-config-dir")
	isGCE = v != ""
}

var (
	isDev             = os.Getenv("CAMLI_DEV_CAMLI_ROOT") != ""
	isDebug, _        = strconv.ParseBool(os.Getenv("CAMLI_DEBUG"))
	isDebugUploads, _ = strconv.ParseBool(os.Getenv("CAMLI_DEBUG_UPLOADS"))
)
