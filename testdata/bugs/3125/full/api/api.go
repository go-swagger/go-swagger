//go:build testintegration

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// FooBarRequest represents body of FooBar request.
type FooBarRequest struct {
	// Foo param
	Foo string `json:"foo"`
	// Bar params
	Bar []int `json:"bar"`
	// User param
	User User `json:"user"`
}

// FooBarResponse represents body of FooBar response.
type FooBarResponse struct {
	Baz struct {
		Prop string `json:"prop"`
	} `json:"baz"`
}

// FooBarHandler handles incoming foobar requests
func FooBarHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", http.StatusText(http.StatusBadRequest), err), http.StatusBadRequest)
		return
	}
	raw := req.FormValue("age")
	age, err := strconv.Atoi(raw)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", http.StatusText(http.StatusBadRequest), err), http.StatusBadRequest)
		return
	}

	r := FooBarRequest{
		Foo: req.FormValue("foo"),
		User: User{
			Name: req.FormValue("name"),
			Age:  age,
		},
	}

	resp := doSthWithRequest(r)

	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", http.StatusText(http.StatusInternalServerError), err), http.StatusInternalServerError)
		return
	}
}

func doSthWithRequest(req FooBarRequest) FooBarResponse {
	return FooBarResponse{}
}
