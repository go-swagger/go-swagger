---
paths:
  - ".github/workflows/**.yml"
  - ".github/workflows/**.yaml"
---

# GitHub Actions Workflows Formatting and Style Conventions

This rule captures YAML and bash formatting rules to provide a consistent maintainer's experience across CI workflows.

## File Structure

**REQUIRED**:  All github action workflows are organized as a flat structure beneath `.github/workflows/`.

> GitHub does not support a hierarchical organization for workflows yet.

**REQUIRED**:  YAML files are conventionally named `{workflow}.yml`, with the `.yml` extension.

## Code Style & Formatting

### Expression Spacing

**REQUIRED**: All GitHub Actions expressions must have spaces inside the braces:

```yaml
# ✅ CORRECT
env:
  PR_URL: ${{ github.event.pull_request.html_url }}
  TOKEN: ${{ secrets.GITHUB_TOKEN }}

# ❌ WRONG
env:
  PR_URL: ${{github.event.pull_request.html_url}}
  TOKEN: ${{secrets.GITHUB_TOKEN}}
```

> Provides a consistent formatting rule.

### Conditional Syntax

**REQUIRED**: Always use `${{ }}` in `if:` conditions:

```yaml
# ✅ CORRECT
if: ${{ inputs.enable-signing == 'true' }}
if: ${{ github.event.pull_request.user.login == 'dependabot[bot]' }}

# ❌ WRONG (works but inconsistent)
if: inputs.enable-signing == 'true'
```

> Provides a consistent formatting rule.

### GitHub Workflow Commands

**REQUIRED**: Use workflow commands for status messages that should appear as annotations, with **double colon separator**:

```bash
# ✅ CORRECT - Double colon (::) separator after title
echo "::notice title=build::Build completed successfully"
echo "::warning title=race-condition::Merge already in progress"
echo "::error title=deployment::Failed to deploy"

# ❌ WRONG - Single colon separator (won't render as annotation)
echo "::notice title=build:Build completed"  # Missing second ':'
echo "::warning title=x:message"             # Won't display correctly
```

**Syntax pattern:** `::LEVEL title=TITLE::MESSAGE`
- `LEVEL`: notice, warning, or error
- Double `::` separator is required between title and message

> Wrong syntax may raise untidy warnings and produce botched output.

### YAML arrays formatting

For steps, YAML arrays are formatted with the following indentation:

```yaml
# ✅ CORRECT - Clear spacing between steps
    steps:
      -
        name: Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a # v2.5.0
      -
        name: Checkout repository
        uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
        with:
          fetch-depth: 0

# ❌ WRONG - Dense format, more difficult to read
    steps:
      - name: Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a # v2.5.0
      - name: Checkout repository
        uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
        with:
          fetch-depth: 0

# ❌ WRONG - YAML comment or blank line could be avoided
    steps:
      #
      - name: Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@21025c705c08248db411dc16f3619e6b5f9ea21a # v2.5.0

      - name: Checkout repository
        uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2
        with:
          fetch-depth: 0
```

## Security Best Practices

### Version Pinning using SHAs

**REQUIRED**: Always pin action versions to commit SHAs:

> Runs must be repeatable with known pinned version. Automated updates are pushed frequently (e.g. daily or weekly)
> to keep pinned versions up-to-date.

```yaml
# ✅ CORRECT - Pinned to commit SHA with version comment
uses: actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8 # v6.0.1
uses: crazy-max/ghaction-import-gpg@e89d40939c28e39f97cf32126055eeae86ba74ec # v6.3.0

# ❌ WRONG - Mutable tag reference
uses: actions/checkout@v6
```

### Permission settings

**REQUIRED**: Always set minimal permissions at the workflow level.

```yaml
# ✅ CORRECT - Workflow level permissions set to minimum
permissions:
  contents: read

# ❌ WRONG - Workflow level permissions with undue privilege escalation
permissions:
  contents: write
  pull-requests: write
```

**REQUIRED**: Whenever a job needs elevated privileges, always raise required permissions at the job level.

```yaml
# ✅ CORRECT - Job level permissions set to the specific requirements for that job
jobs:
  dependabot:
    permissions:
      contents: write
      pull-requests: write
    uses: ./.github/workflows/auto-merge.yml
    secrets: inherit

# ❌ WRONG - Same permissions but set at workflow level instead of job level
permissions:
  contents: write
  pull-requests: write
```

> (Security best practice detected by CodeQL analysis)

### Undue secret exposure

**NEVER** use `secrets[inputs.name]` — always use explicit secret parameters.

> Using keyed access to secrets forces the runner to expose ALL secrets to the job, which causes a security risk
> (caught and reported by CodeQL security analysis).

```yaml
# ❌ SECURITY VULNERABILITY
# This exposes ALL organization and repository secrets to the runner
on:
  workflow_call:
    inputs:
      secret-name:
        type: string
jobs:
  my-job:
    steps:
      - uses: some-action@v1
        with:
          token: ${{ secrets[inputs.secret-name] }}  # ❌ DANGEROUS!
```

**SOLUTION**: Use explicit secret parameters with fallback for defaults:

```yaml
# ✅ SECURE
on:
  workflow_call:
    secrets:
      gpg-private-key:
        required: false
jobs:
  my-job:
    steps:
      - uses: go-openapi/gh-actions/ci-jobs/bot-credentials@master
        with:
          # Falls back to go-openapi default if not explicitly passed
          gpg-private-key: ${{ secrets.gpg-private-key || secrets.CI_BOT_GPG_PRIVATE_KEY }}
```

## Common Gotchas

### Description fields containing parsable expressions

**REQUIRED**: **DO NOT** use `${{ }}` expressions in description fields:

> They may be parsed by the runner, wrongly interpreted or causing failure (e.g. "not defined in this context").

```yaml
# ❌ WRONG - Can cause YAML parsing errors
description: |
  Pass it as: gpg-private-key: ${{ secrets.MY_KEY }}

# ✅ CORRECT
description: |
  Pass it as: secrets.MY_KEY
```

### Boolean inputs

**Boolean inputs are forbidden**: NEVER use `type: boolean` for workflow inputs due to unpredictable type coercion

> gh-action expressions using boolean job inputs are hard to predict and come with many quirks.

   ```yaml
   # ❌ FORBIDDEN - Boolean inputs have type coercion issues
   on:
     workflow_call:
       inputs:
         enable-feature:
           type: boolean        # ❌ NEVER USE THIS
           default: true

   # The pattern `x == 'true' || x == true` seems safe but fails when:
   # - x is not a boolean: `x == true` evaluates to true if x != null
   # - Type coercion is unpredictable and error-prone

   # ✅ CORRECT - Always use string type for boolean-like inputs
   on:
     workflow_call:
       inputs:
         enable-feature:
           type: string         # ✅ Use string instead
           default: 'true'      # String value

   jobs:
     my-job:
       # Simple, reliable comparison
       if: ${{ inputs.enable-feature == 'true' }}

   # ✅ In bash, this works perfectly (inputs are always strings in bash):
   if [[ '${{ inputs.enable-feature }}' == 'true' ]]; then
     echo "Feature enabled"
   fi
   ```

   **Rule**: Use `type: string` with values `'true'` or `'false'` for all boolean-like workflow inputs.

   **Note**: Step outputs and bash variables are always strings, so `x == 'true'` works fine for those.

### YAML fold scalars in action inputs

**NEVER** use `>` or `>-` (fold scalars) for `with:` input values:

> The YAML spec says fold scalars replace newlines with spaces, but the GitHub Actions runner
> does not reliably honor this for action inputs. The action receives the literal multi-line string
> instead of a single folded line, which breaks flag parsing.

```yaml
# ❌ BROKEN - Fold scalar, args received with embedded newlines
- uses: goreleaser/goreleaser-action@...
  with:
    args: >-
      release
        --clean
        --release-notes /tmp/notes.md

# ✅ CORRECT - Single line
- uses: goreleaser/goreleaser-action@...
  with:
    args: release --clean --release-notes /tmp/notes.md

# ✅ CORRECT - Literal block scalar (|) is fine for run: scripts
- run: |
    echo "line 1"
    echo "line 2"
```

**Rule**: Use single-line strings for `with:` inputs. Only use `|` (literal block scalar) for `run:` scripts where multi-line is intentional.
