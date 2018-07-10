# Use custom templates

When generating a server or client you can specify a directory to load custom templates from 
with `--template-dir`. It will recursively read all the `.gotmpl` files in the directory and 
load them as templates.

<!--more-->

Each file will be loaded and define a template named the same as the file without the suffix. If 
the file is in a subdirectory the directory name will be included in the template name and the
first character of the next path segment will be uppercased. e.g. 
 - template.gotmpl -> template
 - server/test.gotmpl -> serverTest

You can override the following templates. Check go-swagger/generator/templates for the default
definitions.
 
# Available Templates

# Client Templates

## clientFacade
Defined in `client/facade.gotmpl`

---
## clientParameter
Defined in `client/parameter.gotmpl`

---
## clientResponse
Defined in `client/response.gotmpl`
####requires 
 - clientresponse
 - schema
 - docstring

---
## clientresponse
Defined in `client/response.gotmpl`

---
## clientClient
Defined in `client/client.gotmpl`


# Server Templates

## serverParameter
Defined in `server/parameter.gotmpl`
####requires 
 - propertyparamvalidator
 - sliceparambinder

---
## sliceparamvalidator
Defined in `server/parameter.gotmpl`

---
## serverResponses
Defined in `server/responses.gotmpl`
####requires 
 - serverresponse

---
## sliceparambinder
Defined in `server/parameter.gotmpl`
####requires 
 - propertyparamvalidator
 - sliceparambinder

---
## serverresponse
Defined in `server/responses.gotmpl`


---
## serverOperation
Defined in `server/operation.gotmpl`
####requires 
 - schema
 - docstring

---
## propertyparamvalidator
Defined in `server/parameter.gotmpl`
####requires 
 - validationPrimitive
 - sliceparamvalidator

---
## serverMain
Defined in `server/main.gotmpl`

---
## bindprimitiveparam
Defined in `server/parameter.gotmpl`

---
## serverConfigureapi
Defined in `server/configureapi.gotmpl`

---
## serverBuilder
Defined in `server/builder.gotmpl`

---
## serverDoc
Defined in `server/doc.gotmpl`
