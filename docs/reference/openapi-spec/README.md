# OpenAPI 3 (3.0.x and 3.1.0) – Practical Summary and Migration Guide from Swagger 2.0

This document summarizes the OpenAPI 3 Specification (OAS), with an emphasis on the differences from Swagger 2.0 and a concrete checklist to upgrade. It is designed to be a single, self-contained reference.

Note: OpenAPI 3.1.0 aligns with JSON Schema Draft 2020-12 and introduces a few structural changes compared to 3.0.x. If you’re moving from Swagger 2.0, most guidance below targets 3.0, with callouts for 3.1 specifics.

---

## 1) OpenAPI Document Overview

- An OpenAPI document (JSON or YAML) describes your API.
- Top-level fixed fields (3.0):
  - openapi: "3.0.x"
  - info
  - servers
  - paths
  - components
  - security
  - tags
  - externalDocs
- 3.1 additions/changes:
  - openapi: "3.1.0"
  - jsonSchemaDialect (optional; sets default schema dialect for Schema Objects)
  - webhooks (top-level)
  - components supports pathItems
  - Full alignment with JSON Schema Draft 2020-12 in Schema Objects

Key design goals in OAS 3:
- Replace Swagger 2.0’s host/basePath/schemes/consumes/produces with servers and per-media-type content definitions.
- First-class requestBody (removes Swagger 2.0 “in: body” parameter).
- Components as a centralized registry replaces “definitions/parameters/responses/securityDefinitions”.

---

## 2) Top-Level Fields

- openapi: The OAS version string: e.g., "3.0.3" or "3.1.0".
- info: Metadata (title, version, description, termsOfService, contact, license).
  - 3.1: license may include identifier per SPDX, URL optional.
- servers: Replaces host/basePath/schemes. An array of server objects (url, description, variables).
  - Example:
    ```yaml
    servers:
      - url: https://api.example.com/v1
      - url: https://{region}.example.com/{base}
        variables:
          region:
            default: us
            enum: [us, eu]
          base:
            default: api
    ```
- paths: Map of path templates to PathItem objects.
- webhooks (3.1): Map of webhook operation callbacks.
- components: Reusable parts (schemas, responses, parameters, examples, requestBodies, headers, securitySchemes, links, callbacks, pathItems [3.1]).
- security: Global security requirements (array of maps).
- tags: List of Tag objects (name, description, externalDocs).
- externalDocs: ExternalDocumentation object.
- jsonSchemaDialect (3.1): Default dialect URI for all Schema Objects (e.g., "https://json-schema.org/draft/2020-12/schema").

---

## 3) Servers vs Swagger 2.0 host/basePath/schemes

- Swagger 2.0:
  - host: api.example.com
  - basePath: /v1
  - schemes: [https]
  - consumes/produces at global or operation level (e.g., application/json)
- OAS 3:
  - servers replaces host/basePath/schemes. Paths are relative to servers URLs.
  - Negotiated content is specified per requestBody and responses via content: { mediaType: { schema, examples, encoding, … } }.
- Multiple servers are allowed. You can also define servers at the PathItem or Operation level for overrides.

---

## 4) Paths and Operations

- Path templating: /pets/{petId}
- PathItem common fields: summary, description, servers, parameters (applies to all operations under that path).
- Operations (get, put, post, delete, options, head, patch, trace) have:
  - tags, summary, description, externalDocs
  - operationId
  - parameters (no “in: body” in OAS 3)
  - requestBody (replaces “in: body”)
  - responses (required)
  - callbacks
  - deprecated
  - security (operation-level overrides/additions)
  - servers (operation-level overrides/additions)

---

## 5) Parameters

- Locations: path, query, header, cookie (cookie is new in OAS 3).
- No “in: body” anymore; use requestBody instead.
- Shape defined via schema (OpenAPI Schema Object aligns with JSON Schema features; see Section 11).
- Serialization:
  - style and explode controls:
    - Path/header default: simple; Query/cookie default: form.
    - Common styles:
      - simple (CSV for arrays in path/header): color=blue,red
      - form (query/cookie, supports explode): color=blue&color=red when explode: true; color=blue,red when explode: false
      - matrix (path): ;color=blue;size=large or ;color=blue,red
      - label (path): .blue.red
      - spaceDelimited (query arrays): color=blue red
      - pipeDelimited (query arrays): color=blue|red
      - deepObject (query objects): color[R]=100&color[G]=200
  - allowReserved (query): if true, reserved characters are not percent-encoded.
  - allowEmptyValue (query): present in 3.0; in 3.1 it is retained but has no effect on validation; prefer explicit schema to allow empty strings.
- required: true is mandatory for path parameters.

---

## 6) Request Body (new vs Swagger 2.0)

- OAS 3: requestBody object defines payloads (content negotiation):
  - content: map of media type → Media Type Object
  - Each media type can define schema, examples (or example), encoding (for multipart/form-data or x-www-form-urlencoded).
- Swagger 2.0’s body/formData parameters are replaced:
  - For JSON/XML/etc: use requestBody with content.application/json (or xml).
  - For form fields:
    - application/x-www-form-urlencoded or multipart/form-data with a schema describing fields.
    - encoding allows per-field contentType and other hints.
  - File uploads:
    - Use type: string, format: binary (or base64) within a schema under requestBody, and appropriate content type (e.g., multipart/form-data or application/octet-stream).

Example (outline to include):

Form-urlencoded request body:
```yaml
requestBody:
  required: true
  content:
    application/x-www-form-urlencoded:
      schema:
        type: object
        properties:
          username:
            type: string
          password:
            type: string
          tags:
            type: array
            items:
              type: string
        required:
          - username
          - password
      encoding:
        tags:
          style: form
          explode: true
```
Multipart form data with file upload:
```yaml
requestBody:
  required: true
  content:
    multipart/form-data:
      schema:
        type: object
        properties:
          name:
            type: string
          file:
            type: string
            format: binary
          metadata:
            type: string
        required:
          - name
          - file
      encoding:
        file:
          contentType: image/png, image/jpeg
          headers:
            X-Upload-Type:
              schema:
                type: string
```
Multiple content types (content negotiation):
```yaml
requestBody:
  required: true
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/Pet'
    application/xml:
      schema:
        $ref: '#/components/schemas/Pet'
    application/x-www-form-urlencoded:
      schema:
        $ref: '#/components/schemas/Pet'
```

---

## 7) Responses
- The Responses Object maps HTTP status codes (e.g., "200", "404") and an optional default to Response Objects.
- Each Response Object must include a description.
- A Response may define:
  - content: a map of media type → Media Type Object (with schema and example/examples)
  - headers: a map of header definitions
  - links: a map of Link Objects that describe follow-up operations
- Use default for non-specific error handling or catch-all responses.
- For no-body responses (e.g., 204, 304), omit content and provide a description.

Examples
200 JSON success response:
```yaml
responses:
  '200':
    description: Successful operation
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Pet'
        example:
          id: 1
          name: "Fluffy"
          status: "available"
        examples:
          cat:
            summary: "A cat example"
            value:
              id: 1
              name: "Fluffy"
              category:
                id: 1
                name: "Cats"
              status: "available"
          dog:
            summary: "A dog example"
            value:
              id: 2
              name: "Buddy"
              category:
                id: 2
                name: "Dogs"
              status: "pending"
    headers:
      X-Rate-Limit:
        description: Calls per hour allowed by the user
        schema:
          type: integer
          format: int32
```
Default error response:
```yaml
responses:
  default:
    description: Unexpected error
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
        example:
          code: 500
          message: "Internal server error"
```

No-content responses:
```yaml
responses:
  '204':
    description: No content - resource updated successfully
  '304':
    description: Not modified
    headers:
      Last-Modified:
        schema:
          type: string
          format: date-time
```
Response with links:
```yaml
responses:
  '201':
    description: Pet created successfully
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Pet'
    links:
      GetPetById:
        operationId: getPetById
        parameters:
          petId: '$response.body#/id'
        description: The `id` value returned in the response can be used as the `petId` parameter in `GET /pets/{petId}`.
```
## 8) Components

- components is a registry of reusable artifacts:
  - schemas: reusable Schema Objects for payloads and parameters
  - responses: reusable Response Objects (can include content, headers, links)
  - parameters: reusable Parameter Objects (path, query, header, cookie)
  - examples: reusable Example Objects
  - requestBodies: reusable Request Body definitions
  - headers: reusable Header Objects
  - securitySchemes: reusable security scheme definitions
  - links: reusable Link Objects
  - callbacks: reusable Callback Objects
  - pathItems (3.1): reusable Path Item templates
- Conventions and best practices:
  - Naming: use PascalCase or camelCase; avoid spaces and special characters.
  - DRY: prefer $ref to reuse definitions; keep schemas small and composable with allOf/oneOf/anyOf.
  - Stability: avoid breaking changes to widely referenced components; consider versioning component names if needed.
  - Consistency: centralize common headers (e.g., TraceId), pagination parameters, and standard error responses.

Example:

## 9) Examples

- example vs examples:
  - example: a single, anonymous example value (inline).
  - examples: a map of named Example Objects; supports summary, description, value or externalValue.
  - Precedence: when both example and examples are present on the same object, examples takes precedence.
- Where examples can appear:
  - Schema Object (example)
  - Media Type Object (example or examples)
  - Parameter Object (example or examples)
  - Header Object (example or examples)
- externalValue:
  - Use for large payloads or when reusing hosted example files.
  - Mutually exclusive with value (you cannot specify both).
- Referencing reusable examples:
  - Place under components.examples and $ref them at media type/parameter/header locations.

Examples:
Complete components section:
```yaml
components:
  schemas:
    Pet:
      type: object
      required:
        - name
        - photoUrls
      properties:
        id:
          type: integer
          format: int64
          example: 10
        name:
          type: string
          example: "doggie"
        category:
          $ref: '#/components/schemas/Category'
        photoUrls:
          type: array
          items:
            type: string
            format: uri
        tags:
          type: array
          items:
            $ref: '#/components/schemas/Tag'
        status:
          type: string
          enum: [available, pending, sold]
    
    Category:
      type: object
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
    
    Tag:
      type: object
      properties:
        id:
          type: integer
          format: int64
        name:
          type: string
    
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string
        details:
          type: string

  parameters:
    PetIdParam:
      name: petId
      in: path
      required: true
      description: ID of pet to return
      schema:
        type: integer
        format: int64
    
    LimitParam:
      name: limit
      in: query
      description: Maximum number of items to return
      schema:
        type: integer
        format: int32
        minimum: 1
        maximum: 100
        default: 20

  responses:
    NotFound:
      description: The specified resource was not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    
    Unauthorized:
      description: Authentication information is missing or invalid
      headers:
        WWW-Authenticate:
          schema:
            type: string

  requestBodies:
    PetBody:
      description: Pet object that needs to be added to the store
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Pet'
        application/xml:
          schema:
            $ref: '#/components/schemas/Pet'

  headers:
    X-Rate-Limit-Limit:
      description: The number of allowed requests in the current period
      schema:
        type: integer
    
    X-Rate-Limit-Remaining:
      description: The number of remaining requests in the current period
      schema:
        type: integer

  securitySchemes:
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
    
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    
    OAuth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/oauth/authorize
          tokenUrl: https://example.com/oauth/token
          scopes:
            read: Grants read access
            write: Grants write access
            admin: Grants admin access

  examples:
    PetExample:
      summary: A pet example
      description: An example of a pet
      value:
        id: 1
        name: "Fluffy"
        status: "available"
    
    ErrorExample:
      summary: Error response example
      value:
        code: 404
        message: "Pet not found"

  links:
    GetPetById:
      operationId: getPetById
      parameters:
        petId: '$response.body#/id'
    
    GetUserByName:
      operationRef: '#/paths/~1users~1{username}/get'
      parameters:
        username: '$response.body#/username'

  callbacks:
    PetStatusCallback:
      '{$request.body#/callbackUrl}':
        post:
          requestBody:
            required: true
            content:
              application/json:
                schema:
                  type: object
                  properties:
                    petId:
                      type: integer
                    status:
                      type: string
          responses:
            '200':
              description: Callback successfully processed
```
Schema-level example:
```yaml
schema:
  type: object
  properties:
    id:
      type: integer
    name:
      type: string
  example:
    id: 1
    name: "Example Name"
```

Parameter with examples:
```yaml
parameters:
  - name: status
    in: query
    schema:
      type: string
      enum: [available, pending, sold]
    examples:
      available:
        summary: "Available pets"
        value: "available"
      pending:
        summary: "Pending pets"
        value: "pending"
```
External example reference:
```yaml
content:
  application/json:
    schema:
      $ref: '#/components/schemas/LargeDataSet'
    examples:
      dataset1:
        summary: "Large dataset example"
        externalValue: "https://example.com/examples/dataset1.json"
```
Reusable examples in components:
```yaml
components:
  examples:
    UserExample:
      summary: "Standard user example"
      description: "A typical user with basic information"
      value:
        id: 123
        name: "Jane Smith"
        email: "jane@example.com"
        created_at: "2023-01-15T10:30:00Z"
    
    AdminExample:
      summary: "Admin user example"
      value:
        id: 456
        name: "Admin User"
        email: "admin@example.com"
        role: "administrator"
        permissions: ["read", "write", "delete"]

# Referenced in media types:
content:
  application/json:
    schema:
      $ref: '#/components/schemas/User'
    examples:
      user:
        $ref: '#/components/examples/UserExample'
      admin:
        $ref: '#/components/examples/AdminExample'
```
## 10) Links and Callbacks

- Links:
  - Purpose: describe relationships between a response and subsequent operations (enables discoverability).
  - Target an operation via operationId (preferred) or operationRef (absolute/relative ref to an operation path item).
  - Parameter mapping: use runtime expressions to pull values from the current request/response (e.g., $response.body#/id, $request.path.id, $response.header.Location).
  - Links live under a response's links object, keyed by a name.
- Callbacks:
  - Purpose: describe asynchronous, outbound requests that your API will send to a client-provided URL as part of an operation.
  - Structure: a map where each key is a runtime expression resolving to a callback URL; each value is a Path Item describing the callback operations (often POST).
  - Typical use: webhooks-like flows negotiated per request (distinct from 3.1 webhooks which are top-level and static).

Examples:
Links example:
```yaml
paths:
  /users:
    post:
      summary: Create a user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
          links:
            GetUserById:
              operationId: getUserById
              parameters:
                userId: '$response.body#/id'
              description: The `id` value from the created user can be used to fetch the user details.
            
            UpdateUser:
              operationRef: '#/paths/~1users~1{userId}/put'
              parameters:
                userId: '$response.body#/id'
              description: The `id` value can be used to update the user.
            
            GetUserAddress:
              operationId: getUserAddress
              parameters:
                userId: '$response.body#/id'
              description: Get the address for this user.

  /users/{userId}:
    get:
      operationId: getUserById
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
```
Callbacks example:
```yaml
paths:
  /webhooks/subscribe:
    post:
      summary: Subscribe to webhook notifications
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                callbackUrl:
                  type: string
                  format: uri
                  description: URL where notifications will be sent
                events:
                  type: array
                  items:
                    type: string
                  description: List of events to subscribe to
              required:
                - callbackUrl
                - events
      responses:
        '201':
          description: Subscription created successfully
      callbacks:
        orderStatusChanged:
          '{$request.body#/callbackUrl}':
            post:
              summary: Order status changed notification
              requestBody:
                required: true
                content:
                  application/json:
                    schema:
                      type: object
                      properties:
                        orderId:
                          type: integer
                        oldStatus:
                          type: string
                        newStatus:
                          type: string
                        timestamp:
                          type: string
                          format: date-time
              responses:
                '200':
                  description: Notification acknowledged
                '400':
                  description: Bad request
        
        paymentCompleted:
          '{$request.body#/callbackUrl}/payment':
            post:
              summary: Payment completed notification
              requestBody:
                required: true
                content:
                  application/json:
                    schema:
                      type: object
                      properties:
                        paymentId:
                          type: string
                        amount:
                          type: number
                        currency:
                          type: string
                        status:
                          type: string
                          enum: [completed, failed, pending]
              responses:
                '200':
                  description: Payment notification processed
```
Link with multiple parameters and request body:
```yaml
responses:
  '200':
    description: Order created
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Order'
    links:
      AddOrderItem:
        operationId: addOrderItem
        parameters:
          orderId: '$response.body#/id'
        requestBody:
          description: Item to add to the order
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/OrderItem'
        description: Add an item to the created order
```
## 11) Schema Objects

- Types and validation keywords:
  - types: object, array, string, number, integer, boolean; null via nullable (3.0) or union types (3.1).
  - objects: properties, required, additionalProperties (3.0), unevaluatedProperties (3.1).
  - arrays: items, minItems, maxItems, uniqueItems, prefixItems (3.1).
  - composition: allOf, anyOf, oneOf, not.
  - enumerations: enum; const (3.1).
  - strings: minLength, maxLength, pattern, format.
  - numbers: minimum/maximum (and exclusiveMinimum/exclusiveMaximum), multipleOf.
- Nullability:
  - 3.0: nullable: true (OAS extension to JSON Schema).
  - 3.1: use union types (type: [string, 'null']) for nullability.
- Discriminator:
  - Works with oneOf/anyOf to support polymorphic payloads based on a property value.
  - mapping allows explicit schema targets for discriminator values.
- Directional flags:
  - readOnly: present in responses only; writeOnly: present in requests only.
  - deprecated: signals clients that a property is deprecated (non-breaking).
- JSON Schema alignment (3.1):
  - Supports $id, $anchor, $defs, if/then/else, dependentSchemas/required, contentSchema, patternProperties, unevaluatedProperties, etc.
- $ref behavior:
  - 3.0: $ref follows JSON Reference; avoid siblings next to $ref.
  - 3.1: $ref follows JSON Schema 2020-12; still avoid siblings next to $ref for portability.
- Best practices:
  - Keep schemas modular; prefer composition over deep inheritance chains.
  - Use descriptions and examples to communicate intent.
  - Validate formats and constraints against real data.

Examples:
Basic types and validation:
```yaml
components:
  schemas:
    # String with validation
    Username:
      type: string
      minLength: 3
      maxLength: 50
      pattern: '^[a-zA-Z0-9_]+$'
      example: "user123"
    
    # Number with constraints
    Price:
      type: number
      minimum: 0
      maximum: 999999.99
      multipleOf: 0.01
      example: 29.99
    
    # Integer with format
    UserId:
      type: integer
      format: int64
      minimum: 1
      example: 12345
    
    # Boolean
    IsActive:
      type: boolean
      example: true
    
    # Date and time
    CreatedAt:
      type: string
      format: date-time
      example: "2023-12-15T10:30:00Z"
```
Array schemas:
```yaml
components:
  schemas:
    # Simple array
    Tags:
      type: array
      items:
        type: string
      minItems: 1
      maxItems: 10
      uniqueItems: true
      example: ["tag1", "tag2", "tag3"]
    
    # Array of objects
    Users:
      type: array
      items:
        $ref: '#/components/schemas/User'
      example:
        - id: 1
          name: "Alice"
        - id: 2
          name: "Bob"
    
    # Array with prefixItems (3.1)
    Coordinates:
      type: array
      prefixItems:
        - type: number
          description: Latitude
        - type: number
          description: Longitude
      items: false
      minItems: 2
      maxItems: 2
      example: [40.7128, -74.0060]
```
Object schemas with composition:
```yaml
components:
  schemas:
    # Basic object
    Address:
      type: object
      required:
        - street
        - city
        - country
      properties:
        street:
          type: string
          example: "123 Main St"
        city:
          type: string
          example: "New York"
        state:
          type: string
          example: "NY"
        country:
          type: string
          example: "USA"
        zipCode:
          type: string
          pattern: '^\d{5}(-\d{4})?$'
          example: "10001"
      additionalProperties: false
    
    # Inheritance with allOf
    Person:
      type: object
      required:
        - name
        - email
      properties:
        id:
          type: integer
          format: int64
          readOnly: true
        name:
          type: string
          minLength: 1
          maxLength: 100
        email:
          type: string
          format: email
        createdAt:
          type: string
          format: date-time
          readOnly: true
    
    Employee:
      allOf:
        - $ref: '#/components/schemas/Person'
        - type: object
          required:
            - employeeId
            - department
          properties:
            employeeId:
              type: string
              writeOnly: true
            department:
              type: string
              enum: [engineering, sales, marketing, hr]
            salary:
              type: number
              minimum: 0
              writeOnly: true
            startDate:
              type: string
              format: date
```

Polymorphism with discriminator:
```yaml
components:
  schemas:
    Animal:
      type: object
      required:
        - name
        - animalType
      properties:
        name:
          type: string
        animalType:
          type: string
      discriminator:
        propertyName: animalType
        mapping:
          cat: '#/components/schemas/Cat'
          dog: '#/components/schemas/Dog'
    
    Cat:
      allOf:
        - $ref: '#/components/schemas/Animal'
        - type: object
          properties:
            animalType:
              type: string
              const: cat
            livesRemaining:
              type: integer
              minimum: 0
              maximum: 9
    
    Dog:
      allOf:
        - $ref: '#/components/schemas/Animal'
        - type: object
          properties:
            animalType:
              type: string
              const: dog
            breed:
              type: string
            goodBoy:
              type: boolean
              default: true
```
Nullability (3.1 style):
```yaml
components:
  schemas:
    # Union types for nullability in 3.1
    NullableString:
      type: [string, 'null']
      example: null
    
    OptionalUser:
      type: object
      properties:
        name:
          type: string
        email:
          type: [string, 'null']
        age:
          anyOf:
            - type: integer
              minimum: 0
            - type: 'null'
```
Advanced JSON Schema features (3.1):
```yaml
components:
  schemas:
    # Conditional schemas
    UserAccount:
      type: object
      properties:
        accountType:
          type: string
          enum: [basic, premium]
        features:
          type: array
          items:
            type: string
      if:
        properties:
          accountType:
            const: premium
      then:
        properties:
          features:
            minItems: 5
      else:
        properties:
          features:
            maxItems: 3
    
    # Dependent schemas
    BillingAddress:
      type: object
      properties:
        name:
          type: string
        address:
          type: string
        country:
          type: string
      dependentSchemas:
        country:
          if:
            properties:
              country:
                const: "USA"
          then:
            properties:
              state:
                type: string
              zipCode:
                type: string
                pattern: '^\d{5}(-\d{4})?$'
            required: [state, zipCode]
```
## 12) Content and Media Types

- OAS 3 replaces consumes/produces with per-media-type content objects on requestBody and responses.
- Each media type entry can define:
  - schema: the payload shape
  - example/examples: sample payloads
  - encoding: for application/x-www-form-urlencoded and multipart/form-data, fine-tune per-property contentType, style, explode, and headers.
- Common media types and handling:
  - application/json: structured data; leverage JSON Schema features.
  - text/plain: schema type: string.
  - application/octet-stream: binary data via schema type: string, format: binary.
  - application/x-www-form-urlencoded: form fields described by an object schema; arrays may need encoding hints.
  - multipart/form-data: mixed text and files; file parts are schema type: string, format: binary.
  - Vendor-specific (e.g., application/vnd.acme+json) are supported.

Examples:
JSON content:
```yaml
content:
  application/json:
    schema:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        tags:
          type: array
          items:
            type: string
    example:
      id: 1
      name: "Example"
      tags: ["tag1", "tag2"]
    examples:
      minimal:
        summary: "Minimal example"
        value:
          id: 1
          name: "Min"
      full:
        summary: "Full example"
        value:
          id: 1
          name: "Complete Example"
          tags: ["important", "featured"]
```

Multiple media types:
```yaml
content:
  application/json:
    schema:
      $ref: '#/components/schemas/Pet'
    examples:
      dog:
        value:
          name: "Rex"
          type: "dog"
  application/xml:
    schema:
      $ref: '#/components/schemas/Pet'
    example: |
      <Pet>
        <name>Rex</name>
        <type>dog</type>
      </Pet>
  text/plain:
    schema:
      type: string
    example: "Pet name: Rex, Type: dog"
```
Form data with encoding:
```yaml
content:
  application/x-www-form-urlencoded:
    schema:
      type: object
      properties:
        name:
          type: string
        tags:
          type: array
          items:
            type: string
        coordinates:
          type: object
          properties:
            lat:
              type: number
            lng:
              type: number
    encoding:
      tags:
        style: form
        explode: true
      coordinates:
        style: deepObject
        explode: true
    example:
      name: "Location"
      tags: ["outdoor", "scenic"]
      coordinates:
        lat: 40.7128
        lng: -74.0060
```
Multipart form data:
```yaml
content:
  multipart/form-data:
    schema:
      type: object
      properties:
        metadata:
          type: object
          properties:
            title:
              type: string
            description:
              type: string
        profileImage:
          type: string
          format: binary
        attachments:
          type: array
          items:
            type: string
            format: binary
    encoding:
      metadata:
        contentType: application/json
      profileImage:
        contentType: image/png, image/jpeg
        headers:
          X-Custom-Header:
            schema:
              type: string
      attachments:
        contentType: application/octet-stream
        headers:
          X-File-Type:
            schema:
              type: string
```
Binary content:
```yaml
content:
  application/octet-stream:
    schema:
      type: string
      format: binary
    example: "[Binary data]"
  
  image/png:
    schema:
      type: string
      format: binary
  
  application/pdf:
    schema:
      type: string
      format: binary
```

Vendor-specific media types:
```yaml
content:
  application/vnd.api+json:
    schema:
      type: object
      properties:
        data:
          type: object
          properties:
            type:
              type: string
            id:
              type: string
            attributes:
              type: object
    example:
      data:
        type: "users"
        id: "123"
        attributes:
          name: "John Doe"
          email: "john@example.com"
  
  application/vnd.company.myformat+json:
    schema:
      $ref: '#/components/schemas/CustomFormat'
```
Content with different schemas per media type:
```yaml
responses:
  '200':
    description: Success
    content:
      application/json:
        schema:
          type: object
          properties:
            data:
              $ref: '#/components/schemas/User'
            meta:
              type: object
      application/hal+json:
        schema:
          type: object
          properties:
            _embedded:
              type: object
            _links:
              type: object
      text/csv:
        schema:
          type: string
        example: "id,name,email\n1,John,john@example.com"
```
## 13) Security

- Define reusable schemes under components.securitySchemes:
  - http:
    - basic: scheme: basic
    - bearer: scheme: bearer; bearerFormat is an optional hint (e.g., JWT)
  - apiKey: identify where the key is passed (in: header | query | cookie) and its name.
  - oauth2: describe available flows (authorizationCode, clientCredentials, password, implicit [legacy]); each flow lists authorizationUrl, tokenUrl (as applicable), refreshUrl (optional), scopes.
  - openIdConnect: specify openIdConnectUrl (OpenID Connect Discovery endpoint).
- Apply security using security arrays at top-level (global) or per operation:
  - The array is an OR-list of requirement objects.
  - Each requirement object is an AND of schemes; for oauth2 the value is the list of required scopes.

Example:
Security schemes in components:
```yaml
components:
  securitySchemes:
    # API Key in header
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      description: API key for authentication
    
    # API Key in query parameter
    ApiKeyQuery:
      type: apiKey
      in: query
      name: api_key
    
    # API Key in cookie
    ApiKeyCookie:
      type: apiKey
      in: cookie
      name: auth_token
    
    # Basic authentication
    BasicAuth:
      type: http
      scheme: basic
    
    # Bearer token (JWT)
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT token for authentication
    
    # OAuth2 with multiple flows
    OAuth2:
      type: oauth2
      description: OAuth2 authentication
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/oauth/authorize
          tokenUrl: https://example.com/oauth/token
          refreshUrl: https://example.com/oauth/refresh
          scopes:
            read: Read access to resources
            write: Write access to resources
            admin: Administrative access
        
        clientCredentials:
          tokenUrl: https://example.com/oauth/token
          scopes:
            service: Service-to-service access
        
        password:
          tokenUrl: https://example.com/oauth/token
          scopes:
            read: Read access
            write: Write access
    
    # OpenID Connect
    OpenIdConnect:
      type: openIdConnect
      openIdConnectUrl: https://example.com/.well-known/openid_configuration
      description: OpenID Connect authentication
```
Global security:
```yaml
security:
  - ApiKeyAuth: []
  - BearerAuth: []
  - OAuth2:
      - read
      - write
```      
Operation level security (overrides global)
```yaml
paths:
  /public-info:
    get:
      summary: Get public information
      security: []  # No authentication required
      responses:
        '200':
          description: Public information
  
  /admin/users:
    get:
      summary: Get all users (admin only)
      security:
        - OAuth2:
            - admin
        - BearerAuth: []
      responses:
        '200':
          description: List of users
  
  /user/profile:
    get:
      summary: Get user profile
      security:
        - ApiKeyAuth: []
        - BearerAuth: []
        - OAuth2:
            - read
      responses:
        '200':
          description: User profile
    
    put:
      summary: Update user profile
      security:
        - BearerAuth: []
        - OAuth2:
            - write
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserProfile'
      responses:
        '200':
          description: Profile updated

```
Mulitple authentication options (OR relationship):
```yaml
security:
  - ApiKeyAuth: []
  - BearerAuth: []
  - BasicAuth: []
```
Combined authentication requirements:
```yaml
security:
  - ApiKeyAuth: []
    OAuth2: [read]
```
OAuth2 scopes in different operations:
```yaml
paths:
  /posts:
    get:
      security:
        - OAuth2: [read]
      responses:
        '200':
          description: List of posts

    post:
      security:
        - OAuth2: [write]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Post'
      responses:
        '201':
          description: Post created

  /admin/posts:
    delete:
      security:
        - OAuth2: [admin]
      responses:
        '204':
          description: Post deleted
```
## 14) Headers

- Header Object mirrors Parameter Object (without name and in fields).
- Serialization:
  - default style for headers is simple; arrays are comma-separated (e.g., "a,b,c").
  - Use schema to describe type/format; references to components.headers are encouraged for reuse.
- Common patterns:
  - Pagination: X-Total-Count, Link (RFC 5988/8288)
  - Correlation: TraceId, Request-Id

Example:

Response headers:
```yaml
responses:
  '200':
    description: Successful response
    headers:
      X-Rate-Limit-Limit:
        description: The number of allowed requests in the current period
        schema:
          type: integer
        example: 1000

      X-Rate-Limit-Remaining:
        description: The number of remaining requests in the current period
        schema:
          type: integer
        example: 999

      X-Rate-Limit-Reset:
        description: The time when the rate limit resets (Unix timestamp)
        schema:
          type: integer
          format: int64
        example: 1703001600

      ETag:
        description: Entity tag for caching
        schema:
          type: string
        example: '"abc123"'

      Last-Modified:
        description: Last modification date
        schema:
          type: string
          format: date-time
        example: "2023-12-15T10:30:00Z"

      Location:
        description: URL of the created resource
        schema:
          type: string
          format: uri
        example: "https://api.example.com/users/123"

      Content-Range:
        description: Range of returned data
        schema:
          type: string
        example: "bytes 200-1023/2048"

    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
```
Reusable headers in components:
```yaml
components:
  headers:
    X-Rate-Limit-Limit:
      description: The number of allowed requests in the current period
      schema:
        type: integer
        minimum: 1
      example: 1000
    
    X-Rate-Limit-Remaining:
      description: The number of remaining requests in the current period
      schema:
        type: integer
        minimum: 0
    
    X-Request-ID:
      description: Unique request identifier for tracing
      schema:
        type: string
        format: uuid
      example: "123e4567-e89b-12d3-a456-426614174000"
    
    X-Correlation-ID:
      description: Correlation identifier for request tracking
      schema:
        type: string
      example: "abc-123-def-456"
    
    X-API-Version:
      description: API version used for this request
      schema:
        type: string
      example: "v1.2.3"
    
    X-Total-Count:
      description: Total number of items (for pagination)
      schema:
        type: integer
        minimum: 0
      example: 42

# Using reusable headers
responses:
  '200':
    description: Success
    headers:
      X-Rate-Limit-Limit:
        $ref: '#/components/headers/X-Rate-Limit-Limit'
      X-Rate-Limit-Remaining:
        $ref: '#/components/headers/X-Rate-Limit-Remaining'
      X-Request-ID:
        $ref: '#/components/headers/X-Request-ID'
```
Headers with complex schemas:
```yaml
responses:
  '200':
    description: Success
    headers:
      X-Custom-Data:
        description: Custom metadata
        schema:
          type: object
          properties:
            version:
              type: string
            timestamp:
              type: string
              format: date-time
            region:
              type: string
              enum: [us-east, us-west, eu-west, ap-southeast]
        example:
          version: "1.0"
          timestamp: "2023-12-15T10:30:00Z"
          region: "us-east"
      
      X-Feature-Flags:
        description: Enabled feature flags
        schema:
          type: array
          items:
            type: string
        style: simple
        explode: false
        example: ["feature1", "feature2", "beta-feature"]
```
Headers with different serialization styles:
```yaml
responses:
  '200':
    description: Success
    headers:
      X-Tags-Simple:
        description: Tags using simple style (default for headers)
        schema:
          type: array
          items:
            type: string
        style: simple
        explode: false
        # Results in: tag1,tag2,tag3
        example: ["tag1", "tag2", "tag3"]
      
      X-Coordinates:
        description: Coordinates as comma-separated values
        schema:
          type: array
          items:
            type: number
        style: simple
        explode: false
        example: [40.7128, -74.0060]
        # Results in: 40.7128,-74.0060
```

## 15) Tags and External Docs

- Tags:
  - Fields: name (required), description, externalDocs.
  - Use tags to group operations for navigation in UIs and for modular documentation.
  - Define common tags at the top-level and reference them by name in operations.
- externalDocs:
  - Can be applied at top-level, tag, path item, and operation levels.
  - Useful to link to guides, tutorials, or deeper reference material.

Example:
```yaml
tags:
- name: users
  description: User management operations
  externalDocs:
  description: User API documentation
  url: https://docs.example.com/users

- name: pets
  description: Pet store operations
  externalDocs:
  description: Pet store guide
  url: https://docs.example.com/pets

- name: orders
  description: Order management
  externalDocs:
  description: Order processing workflow
  url: https://docs.example.com/orders

- name: admin
  description: Administrative operations (requires admin privileges)
  externalDocs:
  description: Admin guide
  url: https://docs.example.com/admin

- name: webhooks
  description: Webhook management and callbacks
```
```yaml
externalDocs:
  description: Complete API documentation and guides
  url: https://docs.example.com/api

info:
  title: Pet Store API
  version: 1.0.0
  description: |
    This is a sample Pet Store Server based on the OpenAPI 3.1 specification.

    You can find out more about the Pet Store at [https://petstore.example.com](https://petstore.example.com).
  termsOfService: https://petstore.example.com/terms
  contact:
    name: API Support
    url: https://petstore.example.com/support
    email: support@petstore.example.com
  license:
    name: MIT
    identifier: MIT
```
```yaml
paths:
  /users/{userId}/profile:
    get:
      summary: Get user profile
      description: Retrieve detailed user profile information
      externalDocs:
        description: User profile fields documentation
        url: https://docs.example.com/user-profile-fields
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User profile
```
```yaml
  paths:
  /orders/{orderId}/process:
    post:
      tags:
        - orders
      summary: Process an order
      description: |
        Process a pending order through the fulfillment pipeline.
        This operation may take several minutes to complete.
      externalDocs:
        description: Order processing workflow documentation
        url: https://docs.example.com/order-processing-workflow
      parameters:
        - name: orderId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                priority:
                  type: string
                  enum: [low, normal, high, urgent]
                notes:
                  type: string
      responses:
        '202':
          description: Order processing started
        '409':
          description: Order cannot be processed in current state
```
## 16) Reusable Path Items (3.1)

- Use components.pathItems to define reusable Path Item templates containing operations, parameters, and responses.
- Reference these from paths via $ref to avoid duplication across similar endpoints.
- Use cases:
  - Standardized health/readiness endpoints
  - Shared CRUD patterns across resources
  - Common parameter sets and response envelopes

Example:

## 17) Webhooks (3.1)

- Top-level webhooks describe inbound calls your service receives from external systems (reverse of callbacks).
- Use webhooks when events are initiated externally and delivered to your API; use callbacks when the client provides a URL during an operation and your server calls back as part of that flow.
- Structure: webhooks is a map of names to Path Item Objects (operations, requestBody, responses).

Example:
Basic webhooks definition:
```yaml
openapi: 3.1.0
info:
  title: E-commerce API with Webhooks
  version: 1.0.0

webhooks:
  # Order status change webhook
  orderStatusChanged:
    post:
      summary: Order status changed notification
      description: |
        Sent when an order status changes (e.g., from pending to shipped).
        Configure webhook URL in your account settings.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - eventType
                - orderId
                - oldStatus
                - newStatus
                - timestamp
              properties:
                eventType:
                  type: string
                  const: order.status_changed
                orderId:
                  type: string
                  example: "order_123456"
                oldStatus:
                  type: string
                  enum: [pending, processing, shipped, delivered, cancelled]
                  example: "processing"
                newStatus:
                  type: string
                  enum: [pending, processing, shipped, delivered, cancelled]
                  example: "shipped"
                timestamp:
                  type: string
                  format: date-time
                  example: "2023-12-15T10:30:00Z"
                order:
                  $ref: '#/components/schemas/Order'
                metadata:
                  type: object
                  properties:
                    trackingNumber:
                      type: string
                    carrier:
                      type: string
            examples:
              shipped:
                summary: Order shipped notification
                value:
                  eventType: order.status_changed
                  orderId: "order_123456"
                  oldStatus: "processing"
                  newStatus: "shipped"
                  timestamp: "2023-12-15T10:30:00Z"
                  metadata:
                    trackingNumber: "1Z999AA1234567890"
                    carrier: "UPS"
      responses:
        '200':
          description: Webhook received and processed successfully
        '400':
          description: Invalid webhook payload
        '401':
          description: Webhook signature verification failed
        '500':
          description: Webhook processing failed (will be retried)

  # Payment completed webhook
  paymentCompleted:
    post:
      summary: Payment completed notification
      description: Sent when a payment is successfully processed
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - eventType
                - paymentId
                - orderId
                - amount
                - currency
                - status
                - timestamp
              properties:
                eventType:
                  type: string
                  const: payment.completed
                paymentId:
                  type: string
                  example: "pay_abcd1234"
                orderId:
                  type: string
                  example: "order_123456"
                amount:
                  type: number
                  format: decimal
                  example: 99.99
                currency:
                  type: string
                  pattern: '^[A-Z]{3}$'
                  example: "USD"
                status:
                  type: string
                  const: completed
                timestamp:
                  type: string
                  format: date-time
                paymentMethod:
                  type: object
                  properties:
                    type:
                      type: string
                      enum: [card, bank_transfer, digital_wallet]
                    last4:
                      type: string
                      description: Last 4 digits for card payments
                customer:
                  $ref: '#/components/schemas/Customer'
      responses:
        '200':
          description: Webhook acknowledged
        '422':
          description: Duplicate payment notification

  # User registration webhook
  userRegistered:
    post:
      summary: New user registration notification
      description: Sent when a new user registers on the platform
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - eventType
                - userId
                - timestamp
              properties:
                eventType:
                  type: string
                  const: user.registered
                userId:
                  type: string
                  example: "user_789"
                timestamp:
                  type: string
                  format: date-time
                user:
                  type: object
                  properties:
                    id:
                      type: string
                    email:
                      type: string
                      format: email
                    name:
                      type: string
                    registrationSource:
                      type: string
                      enum: [web, mobile, api]
                    emailVerified:
                      type: boolean
            examples:
              webRegistration:
                summary: Web registration
                value:
                  eventType: user.registered
                  userId: "user_789"
                  timestamp: "2023-12-15T10:30:00Z"
                  user:
                    id: "user_789"
                    email: "john@example.com"
                    name: "John Doe"
                    registrationSource: "web"
                    emailVerified: false
      responses:
        '200':
          description: Webhook processed
        '400':
          description: Invalid user data

  # Subscription updated webhook
  subscriptionUpdated:
    post:
      summary: Subscription updated notification
      description: |
        Sent when a user's subscription changes (upgrade, downgrade, cancellation, renewal).
        Includes both the old and new subscription details.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - eventType
                - subscriptionId
                - userId
                - changeType
                - timestamp
              properties:
                eventType:
                  type: string
                  const: subscription.updated
                subscriptionId:
                  type: string
                userId:
                  type: string
                changeType:
                  type: string
                  enum: [created, upgraded, downgraded, cancelled, renewed, reactivated]
                timestamp:
                  type: string
                  format: date-time
                oldSubscription:
                  $ref: '#/components/schemas/Subscription'
                newSubscription:
                  $ref: '#/components/schemas/Subscription'
                effectiveDate:
                  type: string
                  format: date-time
                  description: When the change takes effect
      responses:
        '200':
          description: Subscription change processed
```

Webhook security and headers:
```yaml
webhooks:
  secureWebhook:
    post:
      summary: Secure webhook with signature verification
      description: |
        This webhook includes signature verification for security.
        The signature is provided in the X-Webhook-Signature header.
      security:
        - WebhookSignature: []
      parameters:
        - name: X-Webhook-Signature
          in: header
          required: true
          schema:
            type: string
          description: HMAC-SHA256 signature of the request body
        - name: X-Webhook-ID
          in: header
          required: true
          schema:
            type: string
          description: Unique webhook delivery ID
        - name: X-Webhook-Timestamp
          in: header
          required: true
          schema:
            type: integer
          description: Unix timestamp of when the webhook was sent
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                eventType:
                  type: string
                data:
                  type: object
      responses:
        '200':
          description: Webhook verified and processed
        '401':
          description: Invalid signature

components:
  securitySchemes:
    WebhookSignature:
      type: apiKey
      in: header
      name: X-Webhook-Signature
      description: HMAC-SHA256 signature for webhook verification
```

Webhook with retry logic documentation:
```yaml
webhooks:
  reliableWebhook:
    post:
      summary: Reliable webhook with retry policy
      description: |
        This webhook follows a retry policy:
        - Initial delivery attempt
        - Retry after 1 minute if failed
        - Retry after 5 minutes if failed
        - Retry after 30 minutes if failed
        - Final retry after 2 hours
        - Webhook is marked as failed after all retries
        
        A webhook is considered failed if:
        - HTTP status code is not 2xx
        - Request times out (30 seconds)
        - Connection cannot be established
        
        Successful delivery requires:
        - HTTP 200-299 status code
        - Response within 30 seconds
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                eventType:
                  type: string
                attemptNumber:
                  type: integer
                  description: Current delivery attempt (1-5)
                  minimum: 1
                  maximum: 5
                maxAttempts:
                  type: integer
                  description: Total number of delivery attempts
                  const: 5
                data:
                  type: object
      responses:
        '200':
          description: Webhook processed successfully
        '202':
          description: Webhook accepted for processing
        '400':
          description: Bad request (will not be retried)
        '401':
          description: Unauthorized (will not be retried)
        '422':
          description: Unprocessable entity (will not be retried)
        '500':
          description: Server error (will be retried)
        '503':
          description: Service unavailable (will be retried)
```

## 18) Specification Extensions

- Extensions use x- prefixed fields and are allowed on most objects.
- Governance:
  - Namespace consistently (e.g., x-org- or x-company-).
  - Document purpose, data types, and consumers; maintain a registry of supported extensions.
  - Avoid leaking internal-only extensions into public specs.
- Tooling:
  - Not all tools preserve or act on extensions; verify behavior in codegen, validators, and documentation generators.
  - Keep extensions optional for clients; do not rely on them for protocol correctness.

Examples:

API-level extensions:
```yaml
openapi: 3.1.0
info:
  title: Example API
  version: 1.0.0
  x-api-id: "api-12345"
  x-business-unit: "e-commerce"
  x-owner-team: "platform-team"
  x-support-contact: "platform-support@example.com"

# Custom server extensions
servers:
  - url: https://api.example.com/v1
    description: Production server
    x-environment: production
    x-region: us-east-1
    x-load-balancer: alb-prod-001
  
  - url: https://staging-api.example.com/v1
    description: Staging server
    x-environment: staging
    x-region: us-east-1
    x-auto-deploy: true

# Global extensions
x-rate-limiting:
  default:
    requests_per_minute: 1000
    burst_limit: 100
  premium:
    requests_per_minute: 5000
    burst_limit: 500

x-response-time-sla:
  p95: "200ms"
  p99: "500ms"

x-monitoring:
  alerts:
    - name: "high_error_rate"
      threshold: "5%"
    - name: "slow_response"
      threshold: "1s"
```

Operation-level extensions:
```yaml
paths:
  /users:
    get:
      summary: List users
      x-internal-only: false
      x-rate-limit: 100  # requests per minute for this endpoint
      x-cache-ttl: 300   # cache for 5 minutes
      x-cost-category: "read"
      x-monitoring:
        critical: true
        alert_on_errors: true
      responses:
        '200':
          description: Users list
    
    post:
      summary: Create user
      x-internal-only: false
      x-rate-limit: 10   # more restrictive for write operations
      x-cost-category: "write"
      x-audit-log: true  # log all create operations
      x-requires-approval: false
      responses:
        '201':
          description: User created

  /admin/users/{userId}/suspend:
    post:
      summary: Suspend user (admin only)
      x-internal-only: true
      x-admin-only: true
      x-audit-log: true
      x-requires-approval: true
      x-approval-workflow: "user-suspension"
      responses:
        '204':
          description: User suspended
```

Schema extensions:
```yaml
components:
  schemas:
    User:
      type: object
      x-table-name: "users"
      x-primary-key: "id"
      x-timestamps: true
      x-soft-deletes: true
      properties:
        id:
          type: integer
          format: int64
          x-auto-increment: true
          x-database-index: true
          readOnly: true
        
        email:
          type: string
          format: email
          x-unique: true
          x-database-index: true
          x-validation:
            - rule: "email"
            - rule: "max:255"
        
        name:
          type: string
          maxLength: 100
          x-searchable: true
          x-display-name: "Full Name"
          x-validation:
            - rule: "required"
            - rule: "string"
            - rule: "max:100"
        
        role:
          type: string
          enum: [user, admin, moderator]
          x-enum-descriptions:
            user: "Regular user with basic permissions"
            admin: "Administrator with full access"
            moderator: "Moderator with content management permissions"
        
        preferences:
          type: object
          x-json-column: true  # stored as JSON in database
          x-default: {}
          properties:
            theme:
              type: string
              enum: [light, dark]
              default: light
            notifications:
              type: boolean
              default: true
    
    # Database relationship extensions
    Order:
      type: object
      x-table-name: "orders"
      x-relationships:
        - type: "belongsTo"
          model: "User"
          foreign_key: "user_id"
        - type: "hasMany"
          model: "OrderItem"
          foreign_key: "order_id"
      properties:
        id:
          type: integer
          x-primary-key: true
        user_id:
          type: integer
          x-foreign-key:
            table: "users"
            column: "id"
        total:
          type: number
          format: decimal
          x-currency: true
          x-precision: 2
```

Component extensions:
```yaml
components:
  parameters:
    PageParam:
      name: page
      in: query
      schema:
        type: integer
        minimum: 1
        default: 1
      x-ui-hint: "pagination"
      x-description-long: |
        Page number for pagination. Pages are 1-indexed.
        Maximum page size is determined by the 'limit' parameter.
  
  responses:
    Success:
      description: Operation successful
      x-business-meaning: "The requested operation completed without errors"
      content:
        application/json:
          schema:
            type: object
            properties:
              status:
                type: string
                x-always-present: true
              data:
                type: object
                x-dynamic-schema: true

  securitySchemes:
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      x-key-management:
        rotation_period: "90 days"
        key_length: 32
        allowed_characters: "alphanumeric"
      x-rate-limiting:
        tier_1: 1000  # requests per hour for tier 1 keys
        tier_2: 5000  # requests per hour for tier 2 keys
        tier_3: 10000 # requests per hour for tier 3 keys

# Tooling-specific extensions
x-codegen:
  package_name: "example_api_client"
  client_name: "ExampleAPIClient"
  generate_models: true
  generate_operations: true
  target_languages: ["python", "javascript", "go"]

x-documentation:
  generator: "custom"
  theme: "company-branded"
  include_examples: true
  group_by_tags: true
  custom_css: "https://cdn.example.com/api-docs.css"

x-deployment:
  docker:
    base_image: "node:18-alpine"
    port: 3000
  kubernetes:
    replicas: 3
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "512Mi"
        cpu: "500m"

# Custom validation extensions
x-validation-rules:
  email_domains:
    allowed: ["example.com", "partner.com"]
    blocked: ["tempmail.com", "10minutemail.com"]
  
  rate_limits:
    global:
      requests_per_second: 100
      burst: 200
    per_user:
      requests_per_minute: 1000
      requests_per_hour: 10000

# Business logic extensions
x-business-rules:
  user_creation:
    - rule: "email_verification_required"
      condition: "always"
    - rule: "admin_approval_required"
      condition: "domain not in allowed_domains"
  
  order_processing:
    - rule: "fraud_check_required"
      condition: "order_total > 1000 OR new_customer"
    - rule: "manual_review_required"
      condition: "high_risk_country OR suspicious_activity"
```

Vendor-specific extensions:
```yaml
# AWS API Gateway extensions
paths:
  /users:
    get:
      x-amazon-apigateway-integration:
        type: aws_proxy
        httpMethod: POST
        uri: arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:GetUsers/invocations
        payloadFormatVersion: "2.0"
      x-amazon-apigateway-request-validator: "all"
      x-amazon-apigateway-cors:
        allowOrigins:
          - "https://example.com"
        allowMethods:
          - GET
          - POST
        allowHeaders:
          - Content-Type
          - Authorization

# Kong API Gateway extensions
x-kong-plugin-rate-limiting:
  minute: 100
  hour: 1000
  policy: "redis"
  redis_host: "redis.example.com"

x-kong-plugin-cors:
  origins:
    - "https://example.com"
  methods:
    - GET
    - POST
  headers:
    - Accept
    - Content-Type
    - Authorization

# Custom company extensions (namespaced)
x-acme-company:
  service_tier: "premium"
  data_classification: "internal"
  compliance_requirements:
    - "PCI-DSS"
    - "SOX"
  business_owner: "product-team"
  technical_owner: "platform-team"
```


## 19) Differences Summary: Swagger 2.0 → OpenAPI 3

- host/basePath/schemes → servers
- consumes/produces → requestBody.content / responses.[code].content
- parameters:
  - body/formData → requestBody
  - cookie location added (in: cookie)
- definitions/parameters/responses/securityDefinitions → components.{schemas,parameters,responses,securitySchemes,...}
- securityDefinitions → components.securitySchemes
- file uploads: type: file → schema: { type: string, format: binary } under appropriate content type
- examples: moved to media type level and components.examples (named); schema-level example remains available
- nullable: vendor extension in 2.0 → nullable (3.0) or union types (3.1)
- callbacks (3.0) and webhooks (3.1) introduced
- jsonSchemaDialect and full JSON Schema 2020-12 alignment in 3.1

## 20) Migration Checklist

- Bump version to openapi: "3.0.3" or "3.1.0" (choose based on tooling).
- Replace host/basePath/schemes with servers.
- Convert consumes/produces to requestBody.content and responses.[code].content.
- Move in: body/in: formData parameters to requestBody; define media types and schemas.
- Review parameter serialization (style, explode, allowReserved) and cookie params.
- Create components.* registries; move reusable definitions and update $ref paths.
- Update file uploads to { type: string, format: binary } with appropriate media types.
- Migrate securityDefinitions to components.securitySchemes; update security requirements.
- Ensure every response has a description; add content where a body is returned; omit content for 204/304.
- Update nullability:
  - 3.0: use nullable; 3.1: use union types
  - Adopt JSON Schema features (e.g., if/then/else) if moving to 3.1.
- Optionally add links, callbacks (3.0), and webhooks (3.1).
- Validate with linters/validators; run contract tests; verify documentation and codegen outputs.

## 21) Common Migration Pitfalls

- Missing response content for responses that return a body; contrast with proper no-body responses (204/304) where content must be omitted.
- Leaving consumes/produces in place (ignored by OAS 3 tooling).
- Retaining in: body or in: formData parameters (invalid in OAS 3); must use requestBody.
- Incorrect file uploads (must be type: string, format: binary or base64; correct media type).
- $ref alongside sibling keywords (generally ignored by tools); keep $ref as the only key (except description in some tools—avoid for portability).
- Nullability confusion:
  - Using nullable in 3.1 (prefer union types) or forgetting to allow null where needed.
- Misplaced examples (schema vs media type) leading to tools not displaying them.
- Omitting cookie parameters migration and serialization rules.
- Forgetting to quote numeric status codes in YAML (e.g., '200').
- Not updating securityDefinitions to components.securitySchemes and security requirements.

## 22) Minimal Skeleton Examples

3.0.3 minimal:
```yaml
openapi: 3.0.3
info:
  title: Minimal API
  version: 1.0.0
  description: A minimal OpenAPI 3.0.3 specification example
servers:
  - url: https://api.example.com/v1
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "healthy"
        default:
          description: Unexpected error
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: string
components:
  schemas:
    Error:
      type: object
      required:
        - error
      properties:
        error:
          type: string
```

OpenAPI 3.1.0 minimal skeleton:
```yaml
openapi: 3.1.0
info:
  title: Minimal API
  version: 1.0.0
  description: A minimal OpenAPI 3.1.0 specification example
  license:
    name: MIT
    identifier: MIT
jsonSchemaDialect: https://json-schema.org/draft/2020-12/schema
servers:
  - url: https://api.example.com/v1
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    const: "healthy"
        default:
          description: Unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    Error:
      type: object
      required:
        - error
      properties:
        error:
          type: string
          description: Error message
        code:
          type: ["integer", "null"]
          description: Error code (if applicable)
```

Complete example with common patterns (3.1.0):
```yaml
openapi: 3.1.0
info:
  title: Pet Store API
  version: 1.0.0
  description: |
    A sample Pet Store Server demonstrating OpenAPI 3.1 features.
    
    This API showcases:
    - JSON Schema 2020-12 alignment
    - Modern authentication patterns
    - Comprehensive error handling
    - Webhook support
  contact:
    name: API Support
    email: support@petstore.example.com
    url: https://petstore.example.com/support
  license:
    name: MIT
    identifier: MIT

jsonSchemaDialect: https://json-schema.org/draft/2020-12/schema

servers:
  - url: https://api.petstore.example.com/v1
    description: Production server
  - url: https://staging.petstore.example.com/v1
    description: Staging server

security:
  - BearerAuth: []

paths:
  /pets:
    get:
      summary: List all pets
      tags:
        - pets
      parameters:
        - $ref: '#/components/parameters/LimitParam'
        - $ref: '#/components/parameters/OffsetParam'
      responses:
        '200':
          description: A list of pets
          headers:
            X-Total-Count:
              $ref: '#/components/headers/X-Total-Count'
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Pet'
        default:
          $ref: '#/components/responses/Error'
    
    post:
      summary: Create a pet
      tags:
        - pets
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/NewPet'
      responses:
        '201':
          description: Pet created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'
        '400':
          $ref: '#/components/responses/BadRequest'
        default:
          $ref: '#/components/responses/Error'

  /pets/{petId}:
    parameters:
      - $ref: '#/components/parameters/PetIdParam'
    
    get:
      summary: Get a pet by ID
      tags:
        - pets
      responses:
        '200':
          description: Pet details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'
        '404':
          $ref: '#/components/responses/NotFound'
        default:
          $ref: '#/components/responses/Error'

webhooks:
  petStatusChanged:
    post:
      summary: Pet status changed
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - petId
                - oldStatus
                - newStatus
                - timestamp
              properties:
                petId:
                  type: integer
                oldStatus:
                  $ref: '#/components/schemas/PetStatus'
                newStatus:
                  $ref: '#/components/schemas/PetStatus'
                timestamp:
                  type: string
                  format: date-time
      responses:
        '200':
          description: Webhook processed

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  parameters:
    PetIdParam:
      name: petId
      in: path
      required: true
      schema:
        type: integer
        format: int64
        minimum: 1
      description: The ID of the pet
    
    LimitParam:
      name: limit
      in: query
      schema:
        type: integer
        minimum: 1
        maximum: 100
        default: 20
      description: Maximum number of results
    
    OffsetParam:
      name: offset
      in: query
      schema:
        type: integer
        minimum: 0
        default: 0
      description: Number of results to skip

  headers:
    X-Total-Count:
      description: Total number of items
      schema:
        type: integer
        minimum: 0

  responses:
    Error:
      description: Unexpected error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    
    BadRequest:
      description: Bad request
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ValidationError'
    
    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'

  schemas:
    Pet:
      type: object
      required:
        - id
        - name
        - status
      properties:
        id:
          type: integer
          format: int64
          readOnly: true
          examples: [1, 2, 3]
        name:
          type: string
          minLength: 1
          maxLength: 100
          examples: ["Fluffy", "Rex", "Whiskers"]
        status:
          $ref: '#/components/schemas/PetStatus'
        tags:
          type: array
          items:
            type: string
          examples: [["friendly", "playful"], ["guard-dog"], []]
        createdAt:
          type: string
          format: date-time
          readOnly: true
    
    NewPet:
      type: object
      required:
        - name
        - status
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 100
        status:
          $ref: '#/components/schemas/PetStatus'
        tags:
          type: array
          items:
            type: string
    
    PetStatus:
      type: string
      enum:
        - available
        - pending
        - sold
      examples: ["available", "pending"]
    
    Error:
      type: object
      required:
        - message
      properties:
        message:
          type: string
          description: Human-readable error message
        code:
          type: ["string", "null"]
          description: Error code for programmatic handling
        details:
          type: ["object", "null"]
          description: Additional error details
    
    ValidationError:
      allOf:
        - $ref: '#/components/schemas/Error'
        - type: object
          properties:
            errors:
              type: array
              items:
                type: object
                properties:
                  field:
                    type: string
                  message:
                    type: string
                  code:
                    type: string

tags:
  - name: pets
    description: Pet operations
    externalDocs:
      description: Pet care guide
      url: https://petstore.example.com/pet-care
```

## 23) Quick Reference: Where Did It Go?

- host/basePath/schemes → servers
- definitions → components.schemas
- parameters (global) → components.parameters
- responses (global) → components.responses
- securityDefinitions → components.securitySchemes
- consumes/produces → requestBody.content / responses.[code].content
- body/formData parameter → requestBody
- file (type) → string with format: binary/base64
- examples (2.0) → example/examples at media type, parameter, or header; components.examples for reuse
- schemes+consumes/produces → replaced by servers and content per media type

## 24) Validation Tips

- Structural checks:
  - Every response has a description; 204/304 have no content.
  - Request/response content types match actual payloads; schemas align.
  - $ref integrity: no broken references; avoid siblings with $ref.
- Parameters:
  - Serialization style/explode aligned with server expectations.
  - Path parameters are required: true and appear in the template.
  - Cookie parameters accounted for where applicable.
- Examples:
  - Place examples at the correct level (media type vs schema vs parameter/header).
  - Keep examples valid against schemas.
- Security:
  - securitySchemes defined and referenced correctly; oauth2 scopes present where required.
- 3.1-specific:
  - JSON Schema 2020-12 compliance (keywords, formats).
  - Set jsonSchemaDialect if using alternatives.
  - Prefer union types for nullability.

Tip:
- Use linters and CI validation to catch regressions; test with documentation generators and codegens you depend on.

## 25) Final Notes

- Version selection:
  - 3.0.3: widest tooling support; good default for conservative environments.
  - 3.1.0: best alignment with modern JSON Schema and adds webhooks/pathItems; verify tool compatibility.
- Tooling:
  - Validate across your toolchain (linters, servers, codegens, docs) before committing to 3.1.
  - Watch for differences in $ref handling and example rendering among tools.
- Authoring best practices:
  - Embrace DRY with components and $ref.
  - Provide rich descriptions and examples for better DX and self-serve onboarding.
  - Keep schemas small, composable, and versioned thoughtfully.
