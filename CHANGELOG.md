# Change Log

## [0.2.7](https://github.com/go-swagger/go-swagger/tree/0.2.7) (2016-02-13)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.6...0.2.7)

## [0.2.6](https://github.com/go-swagger/go-swagger/tree/0.2.6) (2016-02-13)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.5...0.2.6)

## [0.2.5](https://github.com/go-swagger/go-swagger/tree/0.2.5) (2016-02-13)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.4...0.2.5)

## [0.2.4](https://github.com/go-swagger/go-swagger/tree/0.2.4) (2016-02-13)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.3...0.2.4)

**Closed issues:**

- spec generator strips out special characters in the beginning of lines [\#276](https://github.com/go-swagger/go-swagger/issues/276)

## [0.2.3](https://github.com/go-swagger/go-swagger/tree/0.2.3) (2016-02-09)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.2...0.2.3)

## [0.2.2](https://github.com/go-swagger/go-swagger/tree/0.2.2) (2016-02-08)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/0.2.1...0.2.2)

## [0.2.1](https://github.com/go-swagger/go-swagger/tree/0.2.1) (2016-02-07)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/v0.2.0...0.2.1)

**Implemented enhancements:**

- Add a command to initialize a swagger yaml spec [\#187](https://github.com/go-swagger/go-swagger/issues/187)

**Fixed bugs:**

- Kubernetes 2.0 Spec [\#239](https://github.com/go-swagger/go-swagger/issues/239)
- Validation dereferencing non-pointer [\#186](https://github.com/go-swagger/go-swagger/issues/186)
- Generated Validate method does not dereference pointer [\#182](https://github.com/go-swagger/go-swagger/issues/182)
- Generated Validatator on slice from interface method incorrect [\#181](https://github.com/go-swagger/go-swagger/issues/181)

**Closed issues:**

- Nested references in definitions cause failure [\#254](https://github.com/go-swagger/go-swagger/issues/254)
- array body parameters lead to uncompilable client code [\#249](https://github.com/go-swagger/go-swagger/issues/249)
- Optional query param enums are not validated [\#248](https://github.com/go-swagger/go-swagger/issues/248)
- Map as property created as pointer on objects [\#243](https://github.com/go-swagger/go-swagger/issues/243)
- Make generated client use `consumes` in schema for Accept headers, rather than transport consumers [\#235](https://github.com/go-swagger/go-swagger/issues/235)
- wrong identifier used  in generated code when validating parameter with not valid \(in golang\) identifier [\#223](https://github.com/go-swagger/go-swagger/issues/223)
- Delete requests with bodies cause a runtime error [\#219](https://github.com/go-swagger/go-swagger/issues/219)
- Generated server main is currently always overwritten [\#210](https://github.com/go-swagger/go-swagger/issues/210)
- Invalid code client generated in Default Parameter constructor [\#201](https://github.com/go-swagger/go-swagger/issues/201)
- Random model properties [\#180](https://github.com/go-swagger/go-swagger/issues/180)

## [v0.2.0](https://github.com/go-swagger/go-swagger/tree/v0.2.0) (2015-12-25)
[Full Changelog](https://github.com/go-swagger/go-swagger/compare/v0.1.0...v0.2.0)

**Implemented enhancements:**

- Add appveyor build [\#153](https://github.com/go-swagger/go-swagger/issues/153)
- Document supported vendor extensions [\#131](https://github.com/go-swagger/go-swagger/issues/131)
- Add documentation for generated server [\#130](https://github.com/go-swagger/go-swagger/issues/130)

**Fixed bugs:**

- Polymorphic/generic subtypes: discriminator getter method, and unmarshal function do not use definition name [\#175](https://github.com/go-swagger/go-swagger/issues/175)
- Spec generator fails for swagger:route that has no tags [\#171](https://github.com/go-swagger/go-swagger/issues/171)
- client: Polymorphic types as parameter generates pointer to interface [\#169](https://github.com/go-swagger/go-swagger/issues/169)
- client should respect default values [\#135](https://github.com/go-swagger/go-swagger/issues/135)
- models: missing optional fields must not be rejected by validators and must have a distinguishable zero value [\#132](https://github.com/go-swagger/go-swagger/issues/132)

**Closed issues:**

- Add server support for default header values in responses [\#172](https://github.com/go-swagger/go-swagger/issues/172)
- doesn't generate the instagram api server [\#170](https://github.com/go-swagger/go-swagger/issues/170)

## [v0.1.0](https://github.com/go-swagger/go-swagger/tree/v0.1.0) (2015-12-14)
**Implemented enhancements:**

- check licenses of dependencies  [\#154](https://github.com/go-swagger/go-swagger/issues/154)
- Empty or duplicate operation ids in codegen [\#134](https://github.com/go-swagger/go-swagger/issues/134)
- no empty names for path parameters [\#128](https://github.com/go-swagger/go-swagger/issues/128)
- Add validation for only body or formdata params [\#127](https://github.com/go-swagger/go-swagger/issues/127)
- Add support for security definitions to server codegen [\#113](https://github.com/go-swagger/go-swagger/issues/113)
- \[scanner\] security schemes [\#112](https://github.com/go-swagger/go-swagger/issues/112)

**Fixed bugs:**

- no code generated to handle unmarshalling a slice of a generic type [\#160](https://github.com/go-swagger/go-swagger/issues/160)
- swagger generate server with `-t` leads to non-compiliable generated code [\#155](https://github.com/go-swagger/go-swagger/issues/155)
- apiKey SecurityDefinitions work only if the header name=security definition name [\#152](https://github.com/go-swagger/go-swagger/issues/152)
- add allowEmptyValue support for a parameter [\#149](https://github.com/go-swagger/go-swagger/issues/149)
- Polymorphic validation code does not invoke generated Getter methods [\#146](https://github.com/go-swagger/go-swagger/issues/146)
- generate commands should work with urls too [\#145](https://github.com/go-swagger/go-swagger/issues/145)
- responses with a body of type interface{} don't render well [\#137](https://github.com/go-swagger/go-swagger/issues/137)
- responses with a schema render lots of extra schemas [\#136](https://github.com/go-swagger/go-swagger/issues/136)
- Validation fails with circular dependency [\#123](https://github.com/go-swagger/go-swagger/issues/123)
- server should have options for SSL when https scheme is present [\#115](https://github.com/go-swagger/go-swagger/issues/115)
- no enum detected for enum properties in combination with allOf [\#107](https://github.com/go-swagger/go-swagger/issues/107)
- Problem with query parameter with type array and collectionFormat: multi [\#106](https://github.com/go-swagger/go-swagger/issues/106)

**Closed issues:**

- server with valid schema and an extra slash \(/\) does not remove the extra [\#167](https://github.com/go-swagger/go-swagger/issues/167)
- divan/num2words causing `go get` failure [\#166](https://github.com/go-swagger/go-swagger/issues/166)
- Sample swagger.yml generating server fails for boolean, integer, number types in query params [\#163](https://github.com/go-swagger/go-swagger/issues/163)
- Does not support json keys that are numerical [\#162](https://github.com/go-swagger/go-swagger/issues/162)
- Support setting fields on interface/discriminated types [\#158](https://github.com/go-swagger/go-swagger/issues/158)
- Add HTTP/2 support [\#156](https://github.com/go-swagger/go-swagger/issues/156)
- Server does not compile if parameter description is missing [\#148](https://github.com/go-swagger/go-swagger/issues/148)
- Client GenCode tries to access field as method and `func\(\) httpkit.JSONConsumer` not being called. [\#147](https://github.com/go-swagger/go-swagger/issues/147)
- \[scanner\] support discriminators [\#142](https://github.com/go-swagger/go-swagger/issues/142)
- panic: assignment to entry in nil map [\#141](https://github.com/go-swagger/go-swagger/issues/141)
- main.go:XX: handler declared and not used [\#133](https://github.com/go-swagger/go-swagger/issues/133)
- codegen should account for reserved words [\#122](https://github.com/go-swagger/go-swagger/issues/122)
- look into using shippable [\#118](https://github.com/go-swagger/go-swagger/issues/118)
- inline schemas in responses fail to generate [\#116](https://github.com/go-swagger/go-swagger/issues/116)
- Consumers do not handle headers with charset in them [\#114](https://github.com/go-swagger/go-swagger/issues/114)
- Can't get models in the definitions [\#111](https://github.com/go-swagger/go-swagger/issues/111)
- untyped additional properties incorrectly flagged as having validations [\#108](https://github.com/go-swagger/go-swagger/issues/108)



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*
