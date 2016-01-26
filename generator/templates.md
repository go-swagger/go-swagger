# Templates

## propertyValidationDocString

Defined in validation/structfield.gotmpl

requires []


## schemaType
Defined in schematype.gotmpl
requires [schemaBody]


## sliceparamvalidator
Defined in server/parameter.gotmpl
requires []


## modelvalidator
Defined in modelvalidator.gotmpl
requires [schemavalidations]


## discriminatedSerializer
Defined in tupleserializer.gotmpl
requires [subTypeBody]


## additionalpropertiesserializer
Defined in additionalpropertiesserializer.gotmpl
requires []


## additionalPropertiesSerializer
Defined in additionalpropertiesserializer.gotmpl
requires [schemaType]


## mapvalidator
Defined in schemavalidator.gotmpl
requires [propertyvalidator]


## primitivefieldvalidator
Defined in schemavalidator.gotmpl
requires []


## slicevalidator
Defined in schemavalidator.gotmpl
requires [propertyvalidator]


## hasDiscriminatedSerializer
Defined in tupleserializer.gotmpl
requires [schemaType withoutBaseTypeBody]


## tupleSerializer
Defined in tupleserializer.gotmpl
requires [schemaType]


## clientFacade
Defined in client/facade.gotmpl
requires []


## clientResponse
Defined in client/response.gotmpl
requires [clientresponse docstring schema]


## withoutBaseTypeBody
Defined in schemabody.gotmpl
requires [privstructfield schemaType structfield tuplefield privtuplefield]


## schemaBody
Defined in schemabody.gotmpl
requires [docstring tuplefield privtuplefield privstructfield structfield propertyValidationDocString schemaType]


## swaggerJsonEmbed
Defined in swagger_json_embed.gotmpl
requires []


## serverMain
Defined in server/main.gotmpl
requires []


## serverParameter
Defined in server/parameter.gotmpl
requires [sliceparambinder propertyparamvalidator]


## serverResponses
Defined in server/responses.gotmpl
requires [serverresponse]


## bindprimitiveparam
Defined in server/parameter.gotmpl
requires []


## clientParameter
Defined in client/parameter.gotmpl
requires []


## validationCustomformat
Defined in validation/customformat.gotmpl
requires []


## propertyvalidator
Defined in schemavalidator.gotmpl
requires [primitivefieldvalidator objectvalidator slicevalidator validationCustomformat mapvalidator]


## schemavalidator
Defined in schemavalidator.gotmpl
requires [validationCustomformat propertyvalidator dereffedSchemaType slicevalidator mapvalidator primitivefieldvalidator]


## validationPrimitive
Defined in validation/primitive.gotmpl
requires []


## propertyparamvalidator
Defined in server/parameter.gotmpl
requires [sliceparamvalidator validationPrimitive]


## privtuplefield
Defined in structfield.gotmpl
requires [schemaType]


## validationStructfield
Defined in validation/structfield.gotmpl
requires []


## schemabody
Defined in schemabody.gotmpl
requires []


## sliceparambinder
Defined in server/parameter.gotmpl
requires [propertyparamvalidator sliceparambinder]


## docstring
Defined in docstring.gotmpl
requires []


## objectvalidator
Defined in schemavalidator.gotmpl
requires [propertyvalidator]


## withBaseTypeBody
Defined in schemabody.gotmpl
requires [privstructfield tuplefield structfield schemaType privtuplefield]


## model
Defined in model.gotmpl
requires [docstring schema header]


## tuplefieldIface
Defined in structfield.gotmpl
requires [docstring propertyValidationDocString schemaType]


## structfieldIface
Defined in structfield.gotmpl
requires [docstring schemaType propertyValidationDocString]


## clientresponse
Defined in client/response.gotmpl
requires []


## subTypeBody
Defined in schemabody.gotmpl
requires [schemaType privstructfield privtuplefield structfield tuplefield]


## privstructfield
Defined in structfield.gotmpl
requires [schemaType]


## tuplefield
Defined in structfield.gotmpl
requires [docstring propertyValidationDocString schemaType]


## dereffedSchemaType
Defined in schematype.gotmpl
requires [schemaBody]


## structfield
Defined in structfield.gotmpl
requires [propertyValidationDocString docstring schemaType]


## serverDoc
Defined in server/doc.gotmpl
requires []


## clientClient
Defined in client/client.gotmpl
requires []


## serverBuilder
Defined in server/builder.gotmpl
requires []


## tupleserializer
Defined in tupleserializer.gotmpl
requires []


## serverresponse
Defined in server/responses.gotmpl
requires []


## schema
Defined in schema.gotmpl
requires [tupleSerializer discriminatedSerializer hasDiscriminatedSerializer tuplefieldIface structfieldIface schemavalidator schemaBody schemaType propertyValidationDocString additionalPropertiesSerializer docstring]


## schematype
Defined in schematype.gotmpl
requires []


## serverOperation
Defined in server/operation.gotmpl
requires [schema docstring]


## serverConfigureapi
Defined in server/configureapi.gotmpl
requires []


## header
Defined in header.gotmpl
requires []


