---
title: About generating models
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about model generation -->

## Model generation

### Custom validation
_Use-Case_: is it possible to write my own validation code for custom types? If so, can someone give me an example?

>Usage example:
>There is bookstore, with info about books: author, title, price.
> And we want that books from some author will not cost more than some price.
> So I want to write and use a function ValidateBookCustom() like:
```golang
if book.author == "Some author" {
    if book.price > 1000 {
       return false
   }
    else return true
}
```


**Answer**: there are several ways to achieve that.
- reusing customized models
- reusing custom go types
- customizing code generation

You should know that models may be generated independently from server, then reused when generating a new server.

You might build on that: generating a first model, customizing the validation code, then reusing this model (possibly with some others) in your servers.

Another way is to use the `x-go-type extension`, to replace type generation with a custom type.

There is the opportunity to get go-swagger to reuse a predefined type to satisfy the definition in the swagger spec.
Imported package and type alias may be specified as options, as shown in this example:
https://github.com/go-swagger/go-swagger/blob/master/fixtures/codegen/existing-model.yml#L99-L103

That example reuses a type provided by a library with a package alias and type name. The code generator will respect this.

You might use both, preparing a customized model from an initially generated structure, then reusing type custom type in other declarations by hinting the generator with x-go-type.

Further, for repetitive customization, you might be willing to customize the generator's templates. Like in [here](https://github.com/go-swagger/go-swagger/blob/master/generator/templates/schemavalidator.gotmpl)
for models, or in [here](https://github.com/go-swagger/go-swagger/blob/master/generator/templates/server/parameter.gotmpl) for inline parameters.

Originally from issues [#997](https://github.com/go-swagger/go-swagger/issues/997) and [#1334](https://github.com/go-swagger/go-swagger/issues/1334)

### Non-required or nullable property?
_Use-Case_: when a definition has a property N, if N is a number and is not required,
the corresponding generated model has the struct flag `omitempty` for N.
This means that when N has been set to 0 the generated JSON omits N, despite it being validly set to 0.

I would still like to allow this variable to be unset, by setting it to null for example.
This will also apply for returning objects that return false and so on.

>The `"omitempty"` option specifies that the field should be omitted from the encoding if the field has an empty value,
>defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
>(from https://golang.org/pkg/encoding/json/#Marshal)

**Hint**: a workaround for this is to use the extension **x-nullable:true** on properties.

Originally from issue [#959](https://github.com/go-swagger/go-swagger/issues/959). (*more discussion on edge cases there*).

Related: [go-openapi/validate#19](https://github.com/go-openapi/validate/issues/19).

### String parameter in body and query
_Use-case_: I want to create an operation with string parameter in body, but go-swagger fails while generating.

When I change body to query, it works. How can I send my parameter in body with type string?

This *works* (param in query):
```YAML
post:
  description: post func
  operationId: postfunc
  parameters:
    - name: myparam
      in: query
      type: string
```
**But this fails (param in body):**
```YAML
post:
  description: post func
  operationId: postfunc
  parameters:
    - name: myparam
      in: body
      type: string
```

**Answer**: add the schema definition in body. This works:
```YAML
post:
  description: post func
  operationId: postfunc
  parameters:
    - name: myparam
      in: body
      required: true
      schema:
        type: string
```

**Hint**: more generally, you might want to check the validity of your spec re the OpenAPI 2.0 schema before trying generation, using the `swagger validate {spec}` command.

Originally from issue [#990](https://github.com/go-swagger/go-swagger/issues/990).

### Request response can have different objects returned based on query parameters
_Use-Case_: I have a POST request that returns different object models based on the query parameters.

*Is there any way to add multiple responses under the swagger route annotation?*

Like:

```YAML
Responses:
  200: response1
  200: response2
  ... etc
```

*Also is it possible to have different models for the request?*

**Answer**: **No**, as this is not supported in Openapi 2.0 specification

That being said, if you specify a wrapper class or base class, you can return multiple responses.

For example (in pseudo-swagger):

``` YAML
ResponseWrapper:
  type: object
  properties:
    response1:
      $ref: '#/definitions/response1'
    response2:
      $ref: '#/definitions/response2'
```

or perhaps more elegantly:
``` YAML
BaseObject:
  type: object
  properties:
    id:
      type: string
      format: uuid

Response1:
  allOf:
    - $ref: '#/definitions/BaseObject'
    - type: object
      properties:
        extendedAttributeForResponse1

Response2:
  allOf:
    - $ref: '#/definitions/BaseObject'
    - type: object
       properties:
         extendedAttribForResponse2
```

Allegedly, with OpenAPI 3.0 you'll be able to use the `anyOf:` operator, with the different response types.

Regarding parameters, you may also achieve this by putting in the path the query parameters that dictate the model.
``` YAML
paths:
  "/something?objectType=thisThing":
      get:
...
  "/something?objectType=otherThing":
      get:
...
```

Originally from issue [#932](https://github.com/go-swagger/go-swagger/issues/932).

### How to validate dates and times?
JSON schema and Swagger (aka OpenAPI 2.0) define ISO-8601 dates as a known format (e.g. date-time, or `yyyy-MM-dd'T'HH:mm:ss.SSS'Z`).

This format definition is used by go-swagger validators.
You just have to define the format as in:
```JSON
{
  "description": "The date and time that the device was registered.",
  "type":"string",
  "format": "date-time"
}
```

The `go-openapi/strfmt` package supports many additional string formats for validation.

Check out for more in [this repo](https://github.com/go-openapi/strfmt/tree/master/README.md). The full API
is documented [here](https://godoc.org/github.com/go-openapi/strfmt).

Regarding dates, this package extends validation to [RFC3339](https://tools.ietf.org/html/rfc3339) full-date format (e.g. "2006-01-02").

Originally from issue [#643](https://github.com/go-swagger/go-swagger/issues/643).

### Accessing the Default return value
_Use-Case_: I was wondering how I would get the default response from the client?

Note: see also [Access HTTP status code from client#597](https://github.com/go-swagger/go-swagger/issues/597).

I have a spec like this:
```YAML
/deploys/{deploy_id}:
  get:
    operationId: getDeploy
    parameters:
      - name: deploy_id
        type: string
        in: path
        required: true
    responses:
      '200':
        description: OK
        schema:
          $ref: "#/definitions/deploy"
      default:
        description: error
        schema:
          $ref: "#/definitions/error"
```
This spec generates two models: `GetDeployOK` and `GetDeployDefault`. The API generated will return the OK case.
```golang
func (a *Client) GetDeploy(params *GetDeployParams, authInfo runtime.ClientAuthInfoWriter) (*GetDeployOK, error) {
    // TODO: Validate the params before sending
    if params == nil {
        params = NewGetDeployParams()
    }

    result, err := a.transport.Submit(&runtime.ClientOperation{
        ID: "getDeploy",
        Method: "GET",
        PathPattern: "/deploys/{deploy_id}",
        ProducesMediaTypes: []string{"application/json"},
        ConsumesMediaTypes: []string{"application/json"},
        Schemes: []string{"https"},
        Params: params,
        Reader: &GetDeployReader{formats: a.formats},
        AuthInfo: authInfo,
    })
    if err != nil {
        return nil, err
    }
    return result.(*GetDeployOK), nil  TODO
}
```
*Does that mean that, if I get a non-2xx response, I should check the err to actually be a `GetDeployDefault` reference?*

Something like:
```golang
resp, err := c.Operations.GetDeploy(&params, authInfo)
if err != nil {
    if casted, ok := err.(models.GetDeployDefault); ok {
        // do something here....
    } else {
        false, err
    }
}
```

>I've been tracing through the code in `Runtime.Submit`: it delegates to the `GetDeployReader.ReadResponse` which makes the distinction.
>However, it remains unclear how that response is actually surfaced.

**Answer**: you can get pretty close to that with something like:
```golang
casted, ok := err.(*operations.GetDeployDefault)
```

Because it's a struct type it will be a pointer.

Originally from issue [#616](https://github.com/go-swagger/go-swagger/issues/616).

### How to avoid deep copies of complex data structures that need to be marshalled across the API?
_Use-Case_:
>An API that provides access to a complex data structure, defined and governed by a subsystem, should not have to spec the
>same data model to be marshalled, as this would require a deep copy of the data structure from the subsystem to the API layer's
>model universe.

*How do others deal with this problem?*

- If your question is "How do I write arbitrary response bodies from go-swagger generated server code?"
  (e.g. from subsystem structs that you have marshalled) then you may want to write your own `middleware.Responder`,
  which gives you direct access to the underlying `http.ResponseWriter`.
  At this point, though, why use go-swagger instead  of a lighter-weight framework?
- If your question is "how can I generate a swagger spec from my subsystem structs?", then you could check out the `swagger generate
  spec` CLI command.

>Further, a subsystem that builds a complex hierarchical data structure to support its own requirements for efficiency,
>access, and serialization does not want the types of the API data model to be injected into its namespace.

>Eventually, the subsystem can exist in many different contexts beyond the API, which is another reason it should not
>become dependent on any API type.


Back to [all contributions](/faq)
