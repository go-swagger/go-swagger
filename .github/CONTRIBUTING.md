You'll find here general guidelines to contribute to this project.
They mostly correspond to standard practices for open source repositories.

We have tried to keep things as simple as possible.

> [!NOTE]
> If you're an experienced go developer on github, then you should just feel at home with us
> and you may well skip the rest of this document.
>
> You'll essentially apply the usual guidelines for a go library project on github.

These guidelines are common to all libraries published on github by the `go-openapi` and `go-swagger` organizations,
so you'll feel at home with any of our projects.

You'll find more detailed (or repo-specific) instructions in the [maintainer's docs][maintainers-doc].

[maintainers-doc]: https://goswagger.io/go-swagger/contributing/ci/

## How can I contribute

There are many ways in which you can contribute, not just code. Here are a few ideas:

- Reporting issues or bugs
- Suggesting improvements
- Documentation
- Art work that makes the project look great
- Code
    - proposing bug fixes and new features that are within the main project scope
    - improving test coverage
    - addressing code quality issues

## Questions & issues

### Asking a question

You may inquire anything about this library by reporting a "Question" issue on github.

You may also join our discord server where you may discuss issues or requests.

[![Discord Server][discord-badge]][discord-url]

[discord-badge]: https://img.shields.io/discord/1446918742398341256?logo=discord&label=discord&color=blue
[discord-url]: https://discord.gg/FfnFYaC3k5

### Reporting issues

Reporting a problem with our libraries _is_ a valuable contribution.
You can do this on the github issues page of this repository.

Please be as specific as possible when describing your issue.

Whenever relevant, please provide information about your environment (go version, OS).

Adding a code snippet to reproduce the issue is great, and a big time saver for maintainers.

### Triaging issues

You can help triage issues which may include:

* reproducing bug reports
* asking for important information, such as version numbers or reproduction instructions
* answering questions and sharing your insight in issue comments

## Code contributions

### Pull requests are always welcome

We are always thrilled to receive pull requests, and we do our best to
process them as fast as possible.

Not sure if that typo is worth a pull request? Do it! We will appreciate it.

If your pull request is not accepted on the first try, don't be discouraged!
If there's a problem with the implementation, hopefully you've received feedback on what to improve.

If you have a lot of ideas or a lot of issues to solve, try to refrain a bit and post focused
pull requests.
Think that they must be reviewed by a maintainer and it is easy to lose track of things on big PRs.

We're trying very hard to keep the go-openapi packages lean and focused.

Together, these packages constitute a toolkit for go developers:
it won't do everything for everybody out of the box,
but everybody can use it to do just about everything related to OpenAPI.

This means that we might decide against incorporating a new feature.

However, there might be a way to implement that feature *on top of* our libraries.

### Environment

You just need a `go` compiler to be installed. No special tools are needed to work with our libraries.

The minimal go compiler version required is always the old stable (latest minor go version - 1).

Our libraries are designed and tested to work on `Linux`, `MacOS` and `Windows`.

If you're used to work with `go` you should already have everything in place.

Although not required, you'll be certainly more productive with a local installation of `golangci-lint`,
the meta-linter our CI uses.

If you don't have it, you may install it like so:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

### Conventions

#### Git flow

Fork the repo and make changes to your fork in a feature branch.

To submit a pull request, push your branch to your fork (e.g. `upstream` remote):
github will propose to open a pull request on the original repository.

Typically you'd follow some common naming conventions:

- if it's a bug fixing branch, name it `fix/XXX-something` where XXX is the number of the
  issue on github
- if it's a feature branch, create an enhancement issue to announce your
  intentions, and name it `feature/XXX-something` where XXX is the number of the issue.

NOTE: we don't enforce naming conventions on branches: it's your fork after all.

#### Tests

Submit unit tests for your changes.

Go has a great built-in test framework ; use it!

Take a look at existing tests for inspiration, and run the full test suite on your branch
before submitting a pull request.

Our CI measures test coverage and the test coverage of every patch.

Although not a blocking step - because there are so many special cases -
this is an indicator that maintainers consider when approving a PR.
Please try your best to cover at least 80% of your patch.

#### Code style

You may read our stance on code style [there](./STYLE.md).

#### Documentation

Don't forget to update the documentation when creating or modifying a feature.

Most documentation for this library is directly found in code as comments for godoc.

The documentation for the `go-swagger` package is published on [the public go docs site][go-doc].

---

Check your documentation changes for clarity, concision, and correctness.

If you want to assess the rendering of your changes when published to `pkg.go.dev`, you may
want to install the `pkgsite` tool proposed by `golang.org`.

```sh
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```

Then run on the repository folder:

```sh
pkgsite .
```

This will run a godoc server locally where you may see the documentation generated from your local repository.

[go-doc]: https://pkg.go.dev/github.com/go-swagger/go-swagger

#### Commit messages

Pull requests descriptions should be as clear as possible and include a
reference to all the issues that they address.

Pull requests must not contain commits from other users or branches.

Commit messages are not required to follow the "conventional commit" rule, but it's certainly a good
thing to follow that convention (e.g. "fix: fixed panic in XYZ", "ci: did this", "feat: did that" ...).

The title in your commit message is used directly to produce our release notes: try to keep them neat.

The commit message body should detail your changes.

If an issue should be closed by a commit, please add this reference in the commit body:

```
* fixes #{issue number}
```

#### Code review

Code review comments may be added to your pull request.

Discuss, then make the suggested modifications and push additional commits to your feature branch.

Be sure to post a comment after pushing. The new commits will show up in the pull
request automatically, but the reviewers will not be notified unless you comment.

Before the pull request is merged,
**make sure that you've squashed your commits into logical units of work**
using `git rebase -i` and `git push -f`.

After every commit the test suite should be passing.

Include documentation changes in the same commit so that a revert would remove all traces of the feature or fix.

#### Sign your work

Software is developed by real people.

The sign-off is a simple line at the end of your commit message,
which certifies that you wrote it or otherwise have the right to
pass it on as an open-source patch.

We require the simple DCO below with an email signing your commit.
PGP-signed commit are greatly appreciated but not required.

The rules are pretty simple:

- read our [DCO][dco-doc] (from [developercertificate.org][dco-source])
- if you agree with these terms, then you just add a line to every git commit message

```
Signed-off-by: Joe Smith <joe@gmail.com>
```

using your real name (sorry, no pseudonyms or anonymous contributions.)

You can add the sign-off when creating the git commit via `git commit -s`.

[dco-doc]: ./DCO.md
[dco-source]: https://developercertificate.org

## Code contributions by AI agents

Our agentic friends are welcome to contribute!

We only have a few demands to keep-up with human maintainers.

1. Issues and PRs written or posted by agents should always mention the original (human) poster for reference
2. We don't accept PRs attributed to agents. We don't want commits signed like "author: @claude.code".
   Agents or bots may coauthor commits, though.
3. Security vulnerability reports by agents should always be reported privately and mention the original (human) poster
   (see also [Security Policy][security-doc]).

[security-doc]: ../SECURITY.md
