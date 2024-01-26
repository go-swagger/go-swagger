---
menu: main
title: About this project
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 10
description: What the go-swagger project is about
---
# Primer

The `go-swagger` and `go-openapi` projects started in 2015 to bring to the golang community
a toolkit to work with OpenAPI 2.0 (aka "Swagger").

This is a 100% community-driven open-source project.

The foundational philosophy of the project is to provide a _toolkit_, that is building blocks
you may assemble and customize for use in your own personal project.

There are too many use-cases with APIs, and too many ways to approach things, to cover them all,
or to cover them to the liking of everyone. We encourage users of `go-swagger` and `go-openapi` to fork,
reuse and copy what they find useful.

[Original presentation of the project in January 2016](../presentations/swagger-golangsf.html).

# Approaches to API development

There are essentially two ways along which your API and its documentation tend to evolve.

We all want to remove the repetitive nature of writing what is essentially boilerplate,
and achieve faster iteration times.

This toolkit helps you achieve these objectives with both approaches.

## Design-first approach

The first approach is an artifact of that first meeting where you hash out what your service will do.

Hopefully you've come up with some document that services as _a kind of contract_ for what the affected people and teams will need to do.

At this stage you probably want to be able to generate a server, and perhaps a client to talk to that server.

Most likely, you want to have your front-end team use mock data for that server, while the back-end team
is working on their part of the application.

This approach is centered on the notion of API contract: define the contract, and the rest would come easily.

{{<hint info>}}
This is the **design-first** approach to Swagger, also known as **contract-first**.
{{</hint>}}

For that approach, `go-swagger` may help in various ways:
* it can build a client or a server
* it can quickly serve a frontend for the contract specification
* for more involved use cases, you may merge, diff and rework specs

**Example**
```sh
swagger generate client --spec swagger.yaml
```

## Code-first approach

Now we're moving on to the second iteration of the project and at this stage, there might be a design meeting for the
new features, or they might just be refinements of the previous APIs.

When no changes are required to the contract, all is great because people can keep doing what they were doing and everybody is happy.

However it might be that some change is required for the front-end, let's say they need a boolean added to some model
because they want to display an "accepted-terms-and-conditions" checkbox.

The back-end can decide to add this to the code, provide the necessary annotation and regenerate the swagger
specification document.

{{<hint info>}}
This is the **code-first** approach to Swagger.
{{</hint>}}

Another variation is when an API server is already out in the wild and the team needs to produce a valid documentation,
following some standards, so the service becomes interoperable with new clients.

For that approach, `go-swagger` may help with a code scanner: your annoted code is scanned to produce a swagger specification.

**Example**
```sh
swagger generate spec ./...
```

# The `go-swagger` toolkit

This toolkit supports both ways. And boy, it's challenging...

It tries to do so while staying as close as possible to the go standard library interfaces. 

It tries to have no opinions besides the fact that documentation is important.

It tries to integrate well with the rest of the go ecosystem as well as the swagger ecosystem.
