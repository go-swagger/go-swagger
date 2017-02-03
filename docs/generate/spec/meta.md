# swagger:meta

The **swagger:meta** annotation flags a file as source for metadata about the API.
This is typically a doc.go file with your package documentation.

You can specify a Consumes and Produces key which has a new content type on each line
Schemes is a tag that is required and allows for a comma separated string composed of:
http, https, ws or wss

Host and BasePath can be specified but those values will be defaults,
they should get substituted when serving the swagger spec.

The description property uses the rest of the comment block as description for the api when not explictily provided

##### Syntax:

```
swagger:meta
```

##### Properties:

Annotation | Format
-----------|--------
**Terms Of Service** | allows for either a url or a free text definition describing the terms of services for the API
**Consumes** | a list of default (global) mime type values, one per line, for the content the API receives
**Produces** | a list of default (global) mime type values, one per line, for the content the API sends
**Schemes** | a list of default schemes the API accept (possible values: http, https, ws, wss) https is preferred as default when configured
**Version** | the current version of the API
**Host** | the host from where the spec is served
**Base path** | the default base path for this API
**Contact** | the name of for the person to contact concerning the API eg. John Doe&nbsp;&lt;john@blogs.com&gt;&nbsp;http://john.blogs.com
**License** | the name of the license followed by the URL of the license eg. MIT http://opensource.org/license/MIT

##### Example:

```go
// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v2
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> http://john.doe.com
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//
// swagger:meta
package classification
```

##### Result

```yaml
---
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
```
