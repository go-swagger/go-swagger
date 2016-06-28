/*
Copyright 2011 Google Inc.

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

// Package googlestorage is simple Google Cloud Storage client.
//
// It does not include any Camlistore-specific logic.
package googlestorage // import "camlistore.org/pkg/googlestorage"

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"camlistore.org/pkg/blob"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	api "google.golang.org/api/storage/v1"
	"google.golang.org/cloud/compute/metadata"
)

const (
	gsAccessURL = "https://storage.googleapis.com"
	// Scope is the OAuth2 scope used for Google Cloud Storage.
	Scope = "https://www.googleapis.com/auth/devstorage.read_write"
)

// A Client provides access to Google Cloud Storage.
type Client struct {
	client  *http.Client
	service *api.Service
}

// An Object holds the name of an object (its bucket and key) within
// Google Cloud Storage.
type Object struct {
	Bucket string
	Key    string
}

func (o *Object) valid() error {
	if o == nil {
		return errors.New("invalid nil Object")
	}
	if o.Bucket == "" {
		return errors.New("missing required Bucket field in Object")
	}
	if o.Key == "" {
		return errors.New("missing required Key field in Object")
	}
	return nil
}

// A SizedObject holds the bucket, key, and size of an object.
type SizedObject struct {
	Object
	Size int64
}

// NewServiceClient returns a Client for use when running on Google
// Compute Engine.  This client can access buckets owned by the same
// project ID as the VM.
func NewServiceClient() (*Client, error) {
	if !metadata.OnGCE() {
		return nil, errors.New("not running on Google Compute Engine")
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
	if !haveScope("https://www.googleapis.com/auth/devstorage.full_control") &&
		!haveScope("https://www.googleapis.com/auth/devstorage.read_write") {
		return nil, errors.New("when this Google Compute Engine VM instance was created, it wasn't granted access to Cloud Storage")
	}
	client := oauth2.NewClient(context.Background(), google.ComputeTokenSource(""))
	service, _ := api.New(client)
	return &Client{client: client, service: service}, nil
}

func NewClient(oauthClient *http.Client) *Client {
	service, _ := api.New(oauthClient)
	return &Client{
		client:  oauthClient,
		service: service,
	}
}

func (o *Object) String() string {
	if o == nil {
		return "<nil *Object>"
	}
	return fmt.Sprintf("%v/%v", o.Bucket, o.Key)
}

func (so SizedObject) String() string {
	return fmt.Sprintf("%v/%v (%vB)", so.Bucket, so.Key, so.Size)
}

// Makes a simple body-less google storage request
func (gsa *Client) simpleRequest(method, url_ string) (resp *http.Response, err error) {
	// Construct the request
	req, err := http.NewRequest(method, url_, nil)
	if err != nil {
		return
	}
	req.Header.Set("x-goog-api-version", "2")

	return gsa.client.Do(req)
}

// GetObject fetches a Google Cloud Storage object.
// The caller must close rc.
func (c *Client) GetObject(obj *Object) (rc io.ReadCloser, size int64, err error) {
	if err = obj.valid(); err != nil {
		return
	}
	resp, err := c.simpleRequest("GET", gsAccessURL+"/"+obj.Bucket+"/"+obj.Key)
	if err != nil {
		return nil, 0, fmt.Errorf("GS GET request failed: %v\n", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, 0, os.ErrNotExist
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("GS GET request failed status: %v\n", resp.Status)
	}

	return resp.Body, resp.ContentLength, nil
}

// GetPartialObject fetches part of a Google Cloud Storage object.
// If length is negative, the rest of the object is returned.
// The caller must close rc.
func (c *Client) GetPartialObject(obj Object, offset, length int64) (rc io.ReadCloser, err error) {
	if offset < 0 || length < 0 {
		return nil, blob.ErrNegativeSubFetch
	}
	if err = obj.valid(); err != nil {
		return
	}

	req, err := http.NewRequest("GET", gsAccessURL+"/"+obj.Bucket+"/"+obj.Key, nil)
	if err != nil {
		return
	}
	req.Header.Set("x-goog-api-version", "2")
	if length >= 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
	} else {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GS GET request failed: %v\n", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, os.ErrNotExist
	}
	if !(resp.StatusCode == http.StatusPartialContent || (offset == 0 && resp.StatusCode == http.StatusOK)) {
		resp.Body.Close()
		if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
			return nil, blob.ErrOutOfRangeOffsetSubFetch
		}
		return nil, fmt.Errorf("GS GET request failed status: %v\n", resp.Status)
	}

	return resp.Body, nil
}

// StatObject checks for the size & existence of a Google Cloud Storage object.
// Non-existence of a file is not an error.
func (gsa *Client) StatObject(obj *Object) (size int64, exists bool, err error) {
	if err = obj.valid(); err != nil {
		return
	}
	res, err := gsa.simpleRequest("HEAD", gsAccessURL+"/"+obj.Bucket+"/"+obj.Key)
	if err != nil {
		return
	}
	res.Body.Close() // per contract but unnecessary for most RoundTrippers

	switch res.StatusCode {
	case http.StatusNotFound:
		return 0, false, nil
	case http.StatusOK:
		if size, err = strconv.ParseInt(res.Header["Content-Length"][0], 10, 64); err != nil {
			return
		}
		return size, true, nil
	default:
		return 0, false, fmt.Errorf("Bad head response code: %v", res.Status)
	}
}

// PutObject uploads a Google Cloud Storage object.
// shouldRetry will be true if the put failed due to authorization, but
// credentials have been refreshed and another attempt is likely to succeed.
// In this case, content will have been consumed.
func (gsa *Client) PutObject(obj *Object, content io.Reader) error {
	if err := obj.valid(); err != nil {
		return err
	}
	const maxSlurp = 2 << 20
	var buf bytes.Buffer
	n, err := io.CopyN(&buf, content, maxSlurp)
	if err != nil && err != io.EOF {
		return err
	}
	contentType := http.DetectContentType(buf.Bytes())
	if contentType == "application/octet-stream" && n < maxSlurp && utf8.Valid(buf.Bytes()) {
		contentType = "text/plain; charset=utf-8"
	}

	objURL := gsAccessURL + "/" + obj.Bucket + "/" + obj.Key
	var req *http.Request
	if req, err = http.NewRequest("PUT", objURL, ioutil.NopCloser(io.MultiReader(&buf, content))); err != nil {
		return err
	}
	req.Header.Set("x-goog-api-version", "2")
	req.Header.Set("Content-Type", contentType)

	var resp *http.Response
	if resp, err = gsa.client.Do(req); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad put response code: %v", resp.Status)
	}
	return nil
}

// DeleteObject removes an object.
func (gsa *Client) DeleteObject(obj *Object) error {
	if err := obj.valid(); err != nil {
		return err
	}
	resp, err := gsa.simpleRequest("DELETE", gsAccessURL+"/"+obj.Bucket+"/"+obj.Key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error deleting %v: bad delete response code: %v", obj, resp.Status)
	}
	return nil
}

// EnumerateObjects lists the objects in a bucket.
// If after is non-empty, listing will begin with lexically greater object names.
// If limit is non-zero, the length of the list will be limited to that number.
func (gsa *Client) EnumerateObjects(bucket, after string, limit int) ([]SizedObject, error) {
	// Build url, with query params
	var params []string
	if after != "" {
		params = append(params, "marker="+url.QueryEscape(after))
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("max-keys=%v", limit))
	}
	query := ""
	if len(params) > 0 {
		query = "?" + strings.Join(params, "&")
	}

	resp, err := gsa.simpleRequest("GET", gsAccessURL+"/"+bucket+"/"+query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad enumerate response code: %v", resp.Status)
	}

	var xres struct {
		Contents []SizedObject
	}
	defer resp.Body.Close()
	if err = xml.NewDecoder(resp.Body).Decode(&xres); err != nil {
		return nil, err
	}

	// Fill in the Bucket on all the SizedObjects
	for _, o := range xres.Contents {
		o.Bucket = bucket
	}

	return xres.Contents, nil
}

// BucketInfo returns information about a bucket.
func (c *Client) BucketInfo(bucket string) (*api.Bucket, error) {
	return c.service.Buckets.Get(bucket).Do()
}
