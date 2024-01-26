---
title: About testing
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about testing -->
## API testing

### Any suggestions how to leverage the generated client for benchmarking the API?
_Use-Case_: I want to benchmark the throughput and latency of different API calls.

>I can think of using the go testing/benchmarking framework and hand role the API calls, or I can leverage the generated client and augment it with a benchmarking framework.

*What have others done?*

At the moment, the toolkit does not generate testing tools. You may be interested in this contribution: <https://github.com/go-openapi/stubs>
(generates random JSON to fill in testcases).

We acknowledge that  API testing is an important use-case. However, it is not yet supported. Pull requests to move on forward in that direction are welcome.

Originally from issue [#787](https://github.com/go-swagger/go-swagger/issues/787).

### Using httptest
_Use-Case_: I would like to use httptest for testing my handlers.
Go-swagger provides a Server, but not a configured handler. 

**Hint**: I use this hack : in a file `test.go` in the restapi folder, I steal the private `configureAPI` function. It works.

```golang
package restapi

import (
    loads "github.com/go-openapi/loads"
    "github.com/pim/pam/poum/restapi/operations"
    "net/http"
)

func getAPI() (*operations.ThefactoryAPI, error) {
    swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
    if err != nil {
        return nil, err
    }
    api := operations.NewThefactoryAPI(swaggerSpec)
    return api, nil
}

func GetAPIHandler() (http.Handler, error) {
    api, err := getAPI()
    if err != nil {
        return nil, err
    }
    h := configureAPI(api)
    err = api.Validate()
    if err != nil {
        return nil, err
    }
    return h, nil
}
```

And I can use this in tests like this:
```golang
handler, err := restapi.GetAPIHandler()
if err != nil {
    t.Fatal("get api handler", err)
}
ts := httptest.NewServer(handler)
defer ts.Close()
res, err := http.Get(ts.URL + "/api/v1/boxes")
```

But, hacking restapi, which use my handlers is cyclic, I can't drop my test near my handlers, and this is still a hack.

*What is the official way to manage handler testing?*

**Hint**: you don't actually need httptest to test the handlers.
A handler is essentially a function of parameters to result.
The result knows how to write itself to a http.ResponseWriter, and you already know that that part works.
So to test a handler what you require is to test just your code.

So to test the AddOne operation from the todo list this, there are 2 functions involved in the implementation.

The first function uses the data from the request to actually write the todo item to a store, this can be tested separately.
```golang
func addItem(item *models.Item) error {
    if item == nil {
        return errors.New(500, "item must be present")
    }

    itemsLock.Lock()
    defer itemsLock.Unlock()

    newID := newItemID()
    item.ID = newID
    items[newID] = item

    return nil
}
```

Then there is the actual handler:

```golang
todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
    if err := addItem(params.Body); err != nil {
        return todos.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
    }
    return todos.NewAddOneCreated().WithPayload(params.Body)
})
```
To test this second function we don't need to use the httptest package, you can assume that that part of the code works. So all you have to test is whether or not you get the right return types for a given set of parameters.

Originally from issue [#719](https://github.com/go-swagger/go-swagger/issues/719).

