/*
Copyright 2014 the Camlistore authors.

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

// Package constants contains Camlistore constants.
//
// This is a leaf package, without dependencies.
package constants // import "camlistore.org/pkg/constants"

// MaxBlobSize is the max size of a single blob in Camlistore, in bytes.
const MaxBlobSize = 16 << 20

// DefaultMaxResizeMem is the default maximum number of bytes that
// will be allocated at peak for uncompressed pixel data while
// generating thumbnails or other resized images.
//
// If a single image is larger than the configured size for an
// ImageHandler, we'll never successfully resize it.  256M is a max
// image of ~9.5kx9.5k*3.
const DefaultMaxResizeMem = 256 << 20
