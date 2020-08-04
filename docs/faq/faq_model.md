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

-------------

_Similar Use-Case_:
>A new requirement is proposed that wants API access to that complex data structure, and we decide to use go-swagger for that
>implementation. It is a pure 'read' requirement, so no need for parameter validation by the API, just the ability to gain an XML
>or JSON form of the data structure by a client.

>Our organization decided to keep the API and subsystem layers separate, and to perform deep copies between them.
>The runtime performance costs are acceptable to us, and worth it to keep API-layer dependencies out of our core library.
>If someone can think of a better solution we would love to know!

>If we define a data model in the swagger specification, we end up having to do a
>deep copy of that data structure from subsystem to API if we want to avoid type injection.

*How do you use the swagger spec to define a raw JSON or XML transfer, defined by the subsystem's types?*

**Hint**: you may use the `x-go-type` model annotation that allows you to use pre-existing types as models, you annotate your spec like this:
https://github.com/go-swagger/go-swagger/blob/master/fixtures/codegen/existing-model.yml#L99-103

Originally from issue [#948](https://github.com/go-swagger/go-swagger/issues/948).

### Extra sections in POST body

_Use-case_: additional properties in object

If I have a swagger spec that expects

{"foo": 123}

and provide

{"foo": 123, "blah": 345}

it happily goes on about this way.

Two questions:
1. can I make it complain if extra stuff is included
2. can I access these extra sections within the go code/handler?

**Answer**: use `additionalProperties: false` or `additionalProperties: true` in your definition.
when it's set to true you'll have a `map[string]interface{}` added.

Originally from issue [#1337](https://github.com/go-swagger/go-swagger/issues/1337).

### How to support generate type int?

_Use-case_: generating `int` types

> I need to use swagger to generate my modes in go code.
> But I find I can hardly generate type `int`, always `int64`.
> Since I need to keep back compatibility for my project, I can hardly change the type.
> So in this case, does go-swagger meet this requirement?

**Answer**:  int is not a good option to support when it comes to contracts.

> Consider the following: you have an arm32 client on which int is int32, however your server is amd64.
> At this stage it's perfectly valid for the server to return int32 max value + 1, this will cause the client to overflow.
> So while go allows int as type I think for API contracts int is too ambiguous as definition leading to subtle but hard to debug failures.
> Similarly other languages may choose to default to int32 type instead of int64 type regardless of platform.

Originally from issue [#1205](https://github.com/go-swagger/go-swagger/issues/1205).

### Generate all models necessary for specified operation

_Use-case_: I'm specifying specific operations and I'd like to restrict the models to those needed for those operations. Is there a way to do that?

**Answer:** when using the generate server command, a repeatable --operation=xxx is available to restrict the scope of operations.

 > NOTE: this option is not available for `generate model`.

Originally from issue [#1427](https://github.com/go-swagger/go-swagger/issues/1427).

### Generated code changes the order of properties in struct

_Use-case_: the generated struct has attributes ordered differently than the original specification

Example:

```yaml
Product:
    type: "object"
    properties:
      product_id:
        type: "string"
      name:
        type: "string"
```
Generated by "swagger generate server":
```go
type Product struct {
	Name string `json:"name,omitempty"`
	ProductID string `json:"product_id,omitempty"`
}
```
I want product_id be the first property of Product struct.
Is there any way to keep the order of properties?

**Answer:** try x-order: n extension

Originally from issue [#1759](https://github.com/go-swagger/go-swagger/issues/1759).

### Fail to use swagger generate model -name

_Use-case_: I met a problem when I tried to rename the filename of the auto-generated model.

Example:

1. `swagger generate model -m ./models/vo  --name Person`
`unknown models: Person`
2. `swagger generate model -m ./models/vo  -name Person`
`unknown models: ame`
3. `swagger generate model -m ./models/vo  -name= Person`
`unknown models: ame=`

```json
{
  "swagger": "2.0",
   "info": {
        "version": "1.0.0",
        "title": "Simple API",
        "description": "A simple API to learn how to write OpenAPI Specification"
    },
   "schemes": [
        "http"
    ],
  "paths": {
      "/persons":{
	      "get":{
		      "summary":"获取一些目标person",
			  "description": "Returns a list containing all persons.",
			  "responses": {
                    "200": {
                        "description": "A list of Person",
                        "schema": {
                            "type": "array",
                            "items": {
                                "properties": {
                                    "firstName": {
                                        "type": "string"
                                    },
                                    "lastName": {
                                        "type": "string"
                                    },
                                    "username": {
                                        "type": "string"
                                    }
                                }
                            }
                        }
                    }
                }
		  }
	}
  }
}
```

**Answer:** you need to make it available with that name in the definitions section then it will know

Originally from issue [#1517](https://github.com/go-swagger/go-swagger/issues/1517).

-------------

Back to [all contributions](README.md#all-contributed-questions)
