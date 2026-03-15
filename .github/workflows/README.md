# CI workflows

* [Tests](#tests)
* [Release](#release)
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

## Release

Triggered by pushing a `v*` tag (e.g. `git tag -a v0.31.0 -m "release message" && git push origin v0.31.0`).

### Jobs (sequential)

1. **update-doc** — builds and deploys the documentation site to GitHub Pages
   (reuses `update-doc.yaml`)
2. **docker-release** — builds and pushes Docker images tagged with the release version and `latest`
   (reuses `build-docker.yaml`, pushes to ghcr.io and Quay)
3. **publish-release** — builds binaries, packages, and creates the GitHub release:
   * sets up Go, UPX (binary compression), git-cliff (release notes), and GPG (artifact signing)
   * extracts the annotated tag message and passes it to git-cliff via `--with-tag-message`
   * generates release notes with git-cliff (uses `GITHUB_TOKEN` for PR/contributor enrichment)
   * runs goreleaser to build cross-platform binaries, deb/rpm packages, source archives,
     checksums, and GPG signatures — then publishes the GitHub release
   * optionally pushes deb and rpm packages to CloudSmith (when `CLOUDSMITH_API_KEY` is set)

### Configuration files

| File | Purpose |
|------|---------|
| `.goreleaser.yaml` | goreleaser v2 config: builds, archives, nfpm packages, signing, release |
| `.cliff.toml` | git-cliff config: commit grouping, changelog template, GitHub remote |

### Required secrets

| Secret | Used by |
|--------|---------|
| `GITHUB_TOKEN` | goreleaser (release creation), git-cliff (GitHub API) |
| `GPG_PRIVATE_KEY` | GPG key import for artifact and package signing |
| `GPG_PASSPHRASE` | passphrase for the GPG key |
| `CR_PAT` | push to ghcr.io |
| `QUAY_USERNAME` / `QUAY_PASS` | push to Quay registry |
| `CLOUDSMITH_API_KEY` | (optional) push deb/rpm to CloudSmith |
| `DISCORD_WEBHOOK_ID` / `DISCORD_WEBHOOK_TOKEN` | Discord release announcement |

### Artifacts produced

* signed binaries (linux, windows, macOS — amd64, arm, arm64, etc.)
* tar.gz/zip archive bundles with LICENSE
* standalone binaries (backward-compatible naming)
* signed deb and rpm packages
* signed source archive
* SHA-256 checksum file (signed)
* Docker images (ghcr.io, Quay)
* GitHub Pages documentation site

## CodeQL

## OpenSSF score card

## Update documentation

## Dependabot auto merge

TODO
