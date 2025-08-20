- swagger:meta
  - Sets info (title, description) from package doc comments; version typically from build/flag.

- swagger:route METHOD /path tag operationId
  - Declares an operation under paths.
  - Example: // swagger:route GET /items items listItems // Lists items.

- swagger:parameters OperationId
  - A struct that defines parameters; field-level comments set in and other attributes.
  - Fields:
    - // in: path|query|header|cookie
    - // required: true
    - // description: ...
    - // style: form|simple|matrix|label|deepObject|spaceDelimited|pipeDelimited
    - // explode: true|false

- swagger:response Name
  - A struct that defines a response payload (schema goes in the “in: body” holder field).

- swagger:model
  - A struct exposed as a reusable schema.

OAS 3.x additive directives (new)
- oas:server
  - // oas:server url="[https://api.example.com/v1](https://api.example.com/v1)" description="prod"
  - Repeatable to add multiple servers.

- oas:serverVar
  - // oas:serverVar name="region" default="us" enum="us,eu"

- oas:requestBody
  - // oas:requestBody required=true description="Create item payload"
  - Applied near swagger:parameters, tells writer to produce requestBody.

- oas:requestContent
  - // oas:requestContent media="application/json,application/x-www-form-urlencoded"

- oas:responseContent
  - // oas:responseContent code=201 media="application/json"

- oas:header
  - // oas:header name="X-Rate-Limit" schema="int32" description="Calls per minute"

- oas:example
  - // oas:example name="success" media="application/json" file="examples/item.json"
  - Or inline: value="{"id":1}"

- oas:securityScheme
  - // oas:securityScheme name="BearerAuth" type="http" scheme="bearer" bearerFormat="JWT"

- oas:link
  - // oas:link name="GetItem" operationId="getItem" parameters.id="$response.body#/id"

- oas:callback
  - // oas:callback name="onEvent" expression="{$request.body#/callbackUrl}"

- oas:webhook
  - // oas:webhook name="orderStatusChanged" method=POST path="/webhooks/order/status"

- 3.1 schema hints (field-level on swagger:model structs)
  - // nullable: true → emit type: ["T","null"]
  - // oneOf: TypeA,TypeB
  - // anyOf: ...
  - // allOf: ...
  - // const: VALUE
  - // if: / then: / else: (document usage as advanced, optional)
