# CI workflows

* [Tests](#tests)
* [Examples regeneration](#examples-regeneration)
* [Release](#release)
* [Docker dev](#docker-dev)
* [Update documentation](#update-documentation)
* [Security & compliance](#security--compliance)
* [Dependabot auto-merge](#dependabot-auto-merge)

## Tests

`test.yaml` — triggered on pull requests and push to master.

* only if code changes or the way we test it
* linting
* build: smoke tests with build and basic commands
  * run on a matrix with 2 latest go versions and os: linux (ubuntu), macos, windows (6 runs)
* unit tests
  * run with gotestsum for summarized output
  * hack on windows to ensure that the TempDir lies on the same drive as the code
  * run on a matrix with 2 latest go versions and os: linux (ubuntu), macos, windows (6 runs)
  * collects code coverage
  * collects test reports
* coverage aggregation
  * upload to codecov in one single pass: too many parallel uploads often trigger failures, the retry action
    doesn't support the latest codecov action (composite)
  * we slightly degrade the reporting accuracy, as platform flags are no longer uploaded to codecov (nit)
* test reports aggregation
  * test report upload to codecov (evaluation purpose)
  * test report publishing on github

TODO: release rehearsal on prepare-release/* branch

## Examples regeneration

A two-phase pipeline that regenerates `go-swagger/examples` from PR changes,
plus a cleanup workflow for closed PRs. The two-phase design ensures untrusted
PR code never runs in a context with access to secrets.

### Phase 1 — `regen-examples.yaml`

Triggered on `pull_request`. Runs only when codegen-relevant files change.

1. Builds the swagger binary from the PR
2. Checks out `go-swagger/examples` (public, no token needed)
3. Runs `hack/regen-samples.sh` + `go build ./...`
4. If changes are detected, produces a `regen.patch` + `diffstat.txt`
5. Uploads the diff as an artifact (retention: 1 day)

**No secrets are used** — safe for fork PRs.

### Phase 2 — `regen-examples-pr.yaml`

Triggered by `workflow_run` after phase 1 completes successfully.

1. Downloads the patch artifact from phase 1
2. Applies it to `go-swagger/examples` using bot credentials (GitHub App + GPG signing)
3. Creates or updates a PR in `go-swagger/examples` via `peter-evans/create-pull-request`
4. Comments on the triggering go-swagger PR with a link and diffstat

**Security model:** this workflow never checks out or executes PR code.
It only applies a patch file (data, not code) produced by phase 1.

### Cleanup — `close-examples-pr.yaml`

Triggered by `pull_request_target` on `closed` events. When a go-swagger PR is
closed **without merge**, closes the corresponding examples PR and deletes its branch.

Uses `pull_request_target` for secret access — no PR code is checked out.

## Release

`release.yaml` — triggered by pushing a `v*` tag (e.g. `git tag -a v0.31.0 -m "release message" && git push origin v0.31.0`).

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

## Docker dev

`master.yaml` — triggered by `workflow_run` after `test` completes on master or `prepare-release/*`.

Builds and pushes Docker dev images (`:latest` tag from master) using `build-docker.yaml`.

`report-docker-vuln.yaml` — triggered by `workflow_run` after `docker-dev` completes.
Downloads the Trivy SARIF report produced during the Docker build and uploads it to
GitHub Advanced Security (code scanning dashboard).

## Update documentation

`doc-latest.yaml` — triggered on push to master or PR when `docs/`, `hack/doc-site/`,
or `update-doc.yaml` change.

Reuses `update-doc.yaml` to build the Hugo documentation site and deploy to GitHub Pages.
On PRs, the build is verified but not deployed.

## Security & compliance

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `codeql.yaml` | Push to master, PR to master, weekly schedule | GitHub CodeQL semantic analysis |
| `scorecard.yaml` | Push to master/tags, weekly schedule, branch protection rule changes | OpenSSF Scorecard supply-chain security |
| `scanner.yaml` | Push to master, daily schedule, branch protection rule changes | Trivy vulnerability + secret scan (repo-level) |

## Dependabot auto-merge

`auto-merge.yaml` — triggered on `pull_request`.

* Auto-approves all dependabot PRs
* Auto-merges (rebase) dependabot PRs for:
  * development dependencies (all updates)
  * `go-openapi` dependencies (minor + patch only)
  * `golang.org` dependencies (all updates)
