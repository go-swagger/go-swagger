/*Package swagger (2.0) provides a powerful interface to your API

Contains an implementation of Swagger 2.0.
It knows how to serialize, deserialize and validate swagger specifications.

Swagger is a simple yet powerful representation of your RESTful API.
With the largest ecosystem of API tooling on the planet, thousands of developers are supporting Swagger
in almost every modern programming language and deployment environment.

With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability.
We created Swagger to help fulfill the promise of APIs.

Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal
build the best possible services with RESTful APIs.Now in version 2.0, Swagger is more enabling than ever.
And it's 100% open source software.

Install:

		go get -u github.com/casualjim/go-swagger/cmd/swagger

The implementation also provides a number of command line tools to help working with swagger.

Currently there is a spec validator tool:

		swagger validate https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json

You can also serve a swagger document with the swagger UI

		swagger ui ./swagger.json

To generate a server for a swagger spec document:

		swagger generate all -f ./swagger.json -A [application-name] [--principal [principal-name]] --include-main --include-ui
*/
package swagger
