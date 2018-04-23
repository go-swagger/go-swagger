# Approaches to API development

There are essentially two ways along which your API and its documentation tend to evolve.
This toolkit aims to support you along both and remove the repetitive nature of writing what is essentially boilerplate code for faster iteration times.

## Design first approach
The first one is an artifact of that first meeting where you hash out what your service will do.
Hopefully you've come up with some document that services as a kind of contract for what the affected people and teams will need to do.

At this stage you probably want to be able to generate a server, perhaps a client to talk to that server.

Most likely, you want to have your front-end team use mock data for that server, while the back-end team
is working on their part of the application.

This is the **design first** approach for swagger.

## Code first approach
Now we're moving on to the second iteration of the project and at this stage, there might be a design meeting for the
new features, or they might just be refinements of the previous APIs. When no changes are required to the contract,
all is great because people can keep doing what they were doing and everybody is happy.

However it might be that some change is required for the front-end, let's say they need a boolean added to some model
because they want to display an "accepted-terms-and-conditions" checkbox.

The back-end can decide to add this to the code, provide the necessary annotation and regenerate the swagger
specification document.

This is the **code first** approach for swagger.

## The toolkit
This toolkit supports both ways.

It tries to do so while staying as close as possible to the go stdlib interfaces, it tries to have no opinions
besides the fact that documentation is important.

It tries to integrate well with the rest of the go ecosystem as well as the swagger ecosystem.

A series of contrib projects will add higher level optional functionality so that you can use the generated code with
confidence whether your application is your personal blog or the next AWS.
