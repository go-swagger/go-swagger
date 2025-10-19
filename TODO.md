Badges:

* container version on github is wrong
* slack status link is down
* license scan with FOSSA is failing
* add 100% go badge (!)

![GitHub top language](https://img.shields.io/github/languages/top/go-swagger/go-swagger)


Badges organization:

Documentation
Quality
Releases
Compliance

![CII Best Practices](https://img.shields.io/cii/:metric/:projectId)

ClearlyDefined score?

CodeFactor?
![CodeFactor Grade](https://img.shields.io/codefactor/grade/github/go-swagger/go-swagger)

doc
* preview doc update
* upgrade hugo
* use theme relearn, that supports versioning

release:
* rewrite notes builder from github api (self: local dev2)
* use goreleaser
* changelog

docker:
* include trivy scan & sbom

compliance:
* find a replacement for FOSSA

openssf:
* signed releases -> goreleaser
* token-permissions
* vulns
* binary artifacts (github release) -> goreleaser

* fuzzing ?
* pinned dependencies: fixed
* security policy: fixed
* SAST: should be ok. Unclear result
* CII-Best-Practices: unclear
* Branch protection: scannr error


test:
* json output for integration tests
