# CI workflows

* [Tests](#tests)
* [GitHub codeQL](#codeql)
* [OpenSSF score card](#openssf-score-card)
* [Update documentation](#update-documentation)
* [Auto-merge](#dependabot-auto-merge)

## Tests

* on pull requests
  * only if code changes or the way we test it
  * linting
  * build : smoke tests with build and basic commands
    * run on a matrix with 2 latest go version and os: linux (ubuntu), macos, windows (6 runs)
  * unit tests
    * run with gotestsum for summarized output
    * hack on windows to ensure that the TempDir lies on the same drive as the code
    * run on a matrix with 2 latest go version and os: linux (ubuntu), macos, windows (6 runs)
    * collects code coverage
    * [x] collects test reports
  * coverage aggregation
    * [x] upload to codecov in one single pass: too many parallel uploads often trigger failures, the retry action 
      doesn't support the latest codecov acion (composite)
    * we slightly degrade the reporting accuracy, as platform flags are no longer uploaded to codecov (nit)
  * test reports aggregation
    * [x] test report upload to codecov (evaluation purpose)
    * [x] test report publishing on github

* on push to master

TODO: release rehearsal on prepare-release/* branch

## CodeQL

## OpenSSF score card

## Update documentation

## Dependabot auto merge

TODO
