---
title: meta
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:meta

The **swagger:meta** annotation flags a file as source for metadata about the API.
This is typically a doc.go file with your package documentation.

You can specify a Consumes and Produces key which has a new content type on each line
Schemes is a tag that is required and allows for a comma separated string composed of:
http, https, ws or wss

Host and BasePath can be specified but those values will be defaults,
they should get substituted when serving the swagger spec.

The description property uses the rest of the comment block as description for the api when not explicitly provided

##### Syntax

```go
swagger:meta
```

##### Properties

Annotation | Format
swagger: '2.0'
consumes:
  - application/json
  - application/xml
produces:
  - application/json
  - application/xml
schemes:
  - http
  - https
info:
  description: "the purpose of this application is to provide an application\nthat is using plain go code to define an API\n\nThis should demonstrate all the possible comment annotations\nthat are available to turn go code into a fully compliant swagger 2.0 spec"
  title: 'Petstore API.'
  termsOfService: 'there are no TOS at this moment, use at your own risk we take no responsibility'
  contact: {name: 'John Doe', url: 'http://john.doe.com', email: john.doe@example.com}
  license: {name: MIT, url: 'http://opensource.org/licenses/MIT'}
  version: 0.0.1
host: localhost
basePath: /v2
x-meta-value: value
x-meta-array:
  - value1
  - value2
x-meta-array-obj:
  - name: obj
    value: field
```

##### Supported MIME types

Consumes      | Produces
