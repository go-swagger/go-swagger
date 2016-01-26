# Template Customization

The `swagger` command allows you to specify a directory to load custom templates from. It will recursively read all the `.gotmpl` files
in the directory and load them as templates.

Each file will define a template that is the name of the file without the .gotmpl. If the file is in a subdirectory the directory name 
will be included in the template name and the first character of the next path segment will be uppercased.

 e.g. 
 - template.gotmpl -> template
 - validation/test.gotmpl -> validationTest

Any templates defined in these files will also be included.
 
# Available Templates:

## objectvalidator
Defined in `schemavalidator.gotmpl`
####requires 
 - propertyvalidator


---
## structfieldIface
Defined in `structfield.gotmpl`
####requires 
 - propertyValidationDocString
 - docstring
 - schemaType

---
## schemavalidator
Defined in `schemavalidator.gotmpl`
####requires 
 - validationCustomformat
 - primitivefieldvalidator
 - slicevalidator
 - mapvalidator
 - propertyvalidator
 - dereffedSchemaType



---
## modelvalidator
Defined in `modelvalidator.gotmpl`
####requires 
 - schemavalidations



---
## schemaType
Defined in `schematype.gotmpl`
####requires 
 - schemaBody



---
## discriminatedSerializer
Defined in `tupleserializer.gotmpl`
####requires 
 - subTypeBody



---
## mapvalidator
Defined in `schemavalidator.gotmpl`
####requires 
 - propertyvalidator



---
## model
Defined in `model.gotmpl`
####requires 
 - header
 - docstring
 - schema


---
## withBaseTypeBody
Defined in `schemabody.gotmpl`
####requires 
 - schemaType
 - tuplefield
 - structfield
 - privtuplefield
 - privstructfield



---
## slicevalidator
Defined in `schemavalidator.gotmpl`
####requires 
 - propertyvalidator



---
## tuplefieldIface
Defined in `structfield.gotmpl`
####requires 
 - docstring
 - propertyValidationDocString
 - schemaType



---
## tupleSerializer
Defined in `tupleserializer.gotmpl`
####requires 
 - schemaType

---
## hasDiscriminatedSerializer
Defined in `tupleserializer.gotmpl`
####requires 
 - withoutBaseTypeBody
 - schemaType

---
## primitivefieldvalidator
Defined in `schemavalidator.gotmpl`

---
## schema
Defined in `schema.gotmpl`
####requires 
 - schemaBody
 - discriminatedSerializer
 - hasDiscriminatedSerializer
 - tuplefieldIface
 - structfieldIface
 - schemavalidator
 - docstring
 - tupleSerializer
 - propertyValidationDocString
 - schemaType
 - additionalPropertiesSerializer

---
## privstructfield
Defined in `structfield.gotmpl`
####requires 
 - schemaType

---
## structfield
Defined in `structfield.gotmpl`
####requires 
 - docstring
 - propertyValidationDocString
 - schemaType

---
## docstring
Defined in `docstring.gotmpl`

---
## header
Defined in `header.gotmpl`

---
## tupleserializer
Defined in `tupleserializer.gotmpl`

---
## schemaBody
Defined in `schemabody.gotmpl`
####requires 
 - docstring
 - propertyValidationDocString
 - schemaType
 - privtuplefield
 - privstructfield
 - tuplefield
 - structfield

---
## swaggerJsonEmbed
Defined in `swagger_json_embed.gotmpl`

---
## propertyvalidator
Defined in `schemavalidator.gotmpl`
####requires 
 - slicevalidator
 - validationCustomformat
 - mapvalidator
 - primitivefieldvalidator
 - objectvalidator

---
## additionalPropertiesSerializer
Defined in `additionalpropertiesserializer.gotmpl`
####requires 
 - schemaType

---
## subTypeBody
Defined in `schemabody.gotmpl`
####requires 
 - structfield
 - schemaType
 - tuplefield
 - privtuplefield
 - privstructfield

---
## tuplefield
Defined in `structfield.gotmpl`
####requires 
 - docstring
 - propertyValidationDocString
 - schemaType

---
## dereffedSchemaType
Defined in `schematype.gotmpl`
####requires 
 - schemaBody

---
## privtuplefield
Defined in `structfield.gotmpl`
####requires 
 - schemaType

---
## withoutBaseTypeBody
Defined in `schemabody.gotmpl`
####requires 
 - schemaType
 - tuplefield
 - structfield
 - privtuplefield
 - privstructfield


# Validation Templates

## validationStructfield
Defined in `validation/structfield.gotmpl`

---
## validationPrimitive
Defined in `validation/primitive.gotmpl`

---
## propertyValidationDocString
Defined in `validation/structfield.gotmpl`

---

## validationCustomformat
Defined in `validation/customformat.gotmpl`

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
