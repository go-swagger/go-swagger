# Go-Swagger: Custom Server Example

An example project, showcasing how one can create a custom OpenAPI-based Go server,
using go-swagger to generate its core.

You can regenerate the `./gen` directory using the [swagger][] cli:

```bash
$ rm -rf gen && swagger generate server --exclude-main -A greeter -t gen -f ./swagger/swagger.yml
```

Running the _greeter_ server on port `3000` is as simple as:

```bash
$ go run ./cmd/greeter/main.go --port 3000
```

You can test the server using [httpie][] as follows:

```bash
$ http get :3000/hello                  # returns 'Hello, World!'
$ http get :3000/hello name==Swagger    # returns 'Hello, Swagger!'
```

[swagger]: https://github.com/go-swagger/go-swagger
[httpie]:https://httpie.org
