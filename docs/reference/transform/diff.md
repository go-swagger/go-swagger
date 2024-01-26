---
title: Spec diff
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Generate a Difference Report between two versions of a spec

The toolkit has a command that will let you find the changes in a spec.

### Command usage

```
Usage:
  swagger [OPTIONS] diff [diff-OPTIONS] {original spec} {spec}

diff specs showing which changes will break existing clients

Application Options:
  -q, --quiet                    silence logs
      --log-output=LOG-FILE      redirect logs to file

Help Options:
  -h, --help                     Show this help message

[diff command options]
      -b, --break                When present, only shows incompatible changes
      -f, --format=[txt|json]    When present, writes output as json (default: txt)
      -i, --ignore=              Exception file of diffs to ignore (copy output from json diff format) (default: none specified)
      -d, --dest=                Output destination file or stdout (default: stdout)
```

### Diff output

The output is either a json array of diffs or a human-readable text change report split into breaking and non-breaking changes detected:

```
NON-BREAKING CHANGES:
=====================
/a/:put -  Added endpoint  
/a/{id}:post -  Deleted a deprecated endpoint  
/newpath/:post -  Added endpoint  

BREAKING CHANGES:
=================
/a/:post -  Deleted endpoint  
/b/:post -  Deleted endpoint  
```

### What does calculating Diffs enable me to do?

The challenge of managing changes to API's is one which the industry has wrestled
with for as long as we've had them.
The popularity of microservice architectures talking to each other via REST has multiplied the
complexity, particularly in a Continuous Integration Continuous Delivery (CI/CD) environment.

One popular solution is to use testing to ensure confidence in API compatibility as they change.
There are a few flavours of this but they are all variations on a theme.

  i) Humans imagine an API.
  ii) Humans code the server and one or more clients.
  iii) Humans write tests to ensure that all the humans are doing the right thing.

This is okay, but has challenges of its own.

 i) There's a lot of repetitive boilerplate that humans are forced to write
 ii) The tests can provide feedback which is late or inaccurate. Depending on the actual test
     they may or may not detect backwards incompatible changes. When they do detect it they could flag a failure at a point which is removed from the bug injection. ie I break the server API but
     only find out later when one of the client API tests fails (maybe)

### The Alternative - Spec Driven Development (SDD)

Instead of getting the primates to do all the repetitive typing and testing - we get computers to do it for us.
That's where the go-swagger generate (or swagger-codegen) come into play.
So now the workflow becomes:

  1) Humans imagine an API and express it in a swagger spec
  2) The code for the server endpoints is generated from that spec and wired in to the server app.
  3) The code for the client to access the API (eg RestTemplates/OkHttp etc) is also generated from the same spec and wired in.
  4) No tests need to be written to ensure that the client code matches the server code because they've both been auto-generated from the same spec.

### What About when something changes?

Great question! Glad you're still with us. Yes. The testing outlined in the human generated client/server workflow ensures that the API code matches not just when it's created, but also when it's updated. The most
popular incarnation of this type of testing is known as PACT testing.

#### Aren't all the cool kids using PACT for this?

Methodologies like Consumer-Driven Contract Testing(https://docs.pact.io/how_pact_works)[PACT] testing can give confidence that the API has not been broken... but it suffers from some challenges:

 1) Complexity in setting up PACT brokers and code to capture reproduce server states for interactions where state impacts what gets returned from a server.
 2) Depending on the sequence of execution a client change could cause the server CI/CD pipeline to asyncronously fail as a result of a breaking client change.

#### Using the spec and version control system to ensure backwards compatibility

 Instead of pinning our hopes on after-the-event tests to tell us we did something bad somewhere else,
 why not identify the problem at it's source: the spec itself. If the only thing that can influence the interaction is the API, and the spec defines the API - let's focus on that.

 Given that the spec is well defined, it should be possible to identify changes in a spec that will break. A simple example is deleting an endpoint in use by clients. In general, when both server and client are auto-generated from versions of the same specification, a change will break if:

  i) The server expects something new from the client which it doesn't provide
  ii) The server stops accepting data formats that the client provides.
  iii) The server stops accepting requests on a given endpoint that the client expects to call.
  iii) The server changes the format of returned data.
  iv) The server returns a different return code or different data.

In order to ferret these changes out at the source we add the spec to the server source code repository and then use go-swagger diff to compare any changes against the currently deployed version as a part of the server build.

If any breaking changes in the server are detected, the build fails immediately and the change is never built or deployed. This doesn't suffer from the "random failure" or delayed feedback which is present in relying on downstream testing to ensure backwards compatibility. Making changes to the spec can be done with confidence, because breaking changes are flagged immediately and never propagated to ANY client.

The question is which version do I compare against? The simplest approach would be to compare against the previously committed version but this is subject to flux. The next approach is to compare against the previously pushed version. This is closer to what's in production but may not be exactly what's in prod. The closest fit is for the CI/CD pipeline to tag the server repository with the currently deployed version.

If we ensure that each version of the spec is backwards compatible with it's predecessor then we can use semantic versioning  to ensure that clients on a given version of a spec in production will never be broken by an API change in production.

#### But I want to break stuff

Sometimes to make an omelette you have to break some eggs.

Sometimes you need to make a breaking change to evolve an API.

eg an endpoint needs a new required parameter, or an enum in a response needs a new value.

These cases need to be handled VERY CAREFULLY. Lets take adding the enum value.

 i) You add a new option to cooked_egg by adding "FRIED" to the existing enum ["POACHED","SCRAMBLED"]
 ii) You run swagger diff on your new spec and it says

 ```
  BREAKING CHANGES:
  =================
  /a/:get->200 - Response Body: Added possible enumeration(s) <FRIED>  - array[BreakfastOrder].cooked_egg : string
```
 iii) You now have some choices:
    a) revert, revert, revert - back off and not break the spec
    b) create a new version of the spec in the server and start to migrate clients to that
    c) do a managed migration
 iv)  Let's assume you've decided VEWWWY VEWWWY CAREFUWWY migrate your change. Here's how that would work...

 #### Migrating a breaking change

 As before, same scenario.
i) You add a new option to cooked_egg by adding "FRIED" to the existing enum ["POACHED","SCRAMBLED"]
ii) You run swagger diff on your new spec and it complains as above.
iii) You now generate a diff list using swagger diff -f json to produce a json formatted diff list.
iv) You copy the breaking change you wish to ignore into a swagger_diff_ignore_file.json
v) You ensure your build is configured to read this ignore file so now your build will no longer break.
vi) VERY IMPORTANT: Now you must specifically put a temporary conditional check plus tests in the server to ensure you don't return this new enum value... YET. This gives you a chance to migrate clients to the version with the new enum value.
vii) Deploy the change and then update the clients to use the new version (unless your clever build pipeline does that automatically because you use semantic versioning on your spec)
viii) Once all clients have been updated (**and deployed**) you can remove the item from your ignore file and the code that blocks the server returning the new value.

This type of change would be a nightmare if relying on after-the-event tests...

 #### This may... or may not break a client

What about a change that MAY break a client.

eg returning a new success or error code response

Will that break the client? It... depends. If the client has a firm set of responses it expects then perhaps, yes. New error repsonses are less likely to break the client as there are usually "Something unknown has gone sour" error handling to cope with that.

However, new success messages, say a 201-CREATED instead of a plain 200-SUCCESS might get missed by a client expecting a 200 and only a 200 to indicate a successful call. In this case, use of the ```deprecated``` tag in the swagger spec can indicate that the server has an implendign change and a similar API migration strategy to the enum migration outlined above can be used.

swagger diff is conservative in this regard, preferring to say something is breaking if it MAY block a client. You then have the option to use deprecated or the ignore file to temporarily manage the migration.

