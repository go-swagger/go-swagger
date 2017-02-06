+++
date = "2015-10-21T21:54:47-07:00"
title = "What is go-swagger"
series = "home"
weight = 9
+++

There are 2 axes along which your API and its documentation tend to evolve. This toolkit aims to support you along all
axes and remove the repetitive nature of writing what is essentially boilerplate code for faster iteration times.

<!--more-->

The first one is an artifact of that first meeting where you hash out what your service will do. Hopefully you've come
up with some document that services as a kind of contract for what the affected people/teams will need to do.
At this stage you want to be able to generate a server, perhaps a client to talk to that server.
It's not inconceivable you want to have your front-end team use mock data for that server and that the backend team
wants to be left in peace while they work on their part of the application.

This is the **design first** approach for swagger.

Now we're moving on to the second iteration of the project and at this stage, there might be a design meeting for the
new features, or they might just be refinements of the previous APIs. When no changes are required to the contract,
all is great because people can keep doing what they were doing and everybody is happy.
However it might be that there was a change that is required for the front-end, they need a boolean added to some model
because they want to display an-accepted-terms-and-conditions-checkbox.
The backend can decide to add this to the code, provide the necessary annotation and regenerate the swagger
specification document.

This is the **code first** approach for swagger.

This toolkit aims to support both these modes, remove the repetitive nature of writing what is essentially boilerplate
code. In doing so it tries to stay as close as possible to the go stdlib interfaces, it tries to have no opinions
besides the fact that documentation is important. And it tries to integrate well with the rest of the go ecosystem as
well as the swagger ecosystem.

A series of contrib projects will add higher level optional functionality so that you can use the generated code with
confidence whether your application is your personal blog or is the next AWS.
