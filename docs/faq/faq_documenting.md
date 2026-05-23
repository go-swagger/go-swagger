---
title: About documentating your API
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about the serve UI use-case -->

### API Browser (SwaggerUI/ReDoc) on Generated Servers

Every go-swagger generated server comes with an **API Browser** built-in, so you can easily explore and test your API without any additional configuration.

#### Accessing the API Browser

When you run a generated go-swagger server, the API documentation is automatically available at:

| URL Path | Description |
|----------|-------------|
| `/docs` | ReDoc UI (default) |
| `/swagger.json` | Raw OpenAPI 2.0 specification |

For example, if your server is running on `http://localhost:8080`, you can access:

* **ReDoc UI**: Visit `http://localhost:8080/docs`
* **SwaggerUI**: Visit `http://localhost:8080/docs` (if configured) or `http://localhost:8080/swagger-ui/`
* **Raw Spec**: `http://localhost:8080/swagger.json`

#### Choosing Between ReDoc and SwaggerUI

By default, generated servers use **ReDoc** for API documentation. ReDoc provides a clean, three-panel layout optimized for readability with support for dark mode.

If you prefer SwaggerUI, you can configure it in your `configure_<appname>.go` file:

```go
func setupMiddlewares(handler http.Handler) http.Handler {
    return middleware.SwaggerUI(middleware.SwaggerUIOpts{
        BasePath: "/",
        SpecURL:  "/swagger.json",
        Path:     "swagger-ui",
    }, handler)
}
```

Then access SwaggerUI at `http://localhost:8080/swagger-ui/`.

For more details on serving options, see the [`swagger serve` command](../usage/serve_ui.md).

#### Customizing the API Browser Path

If you need to serve the API documentation at a different path, you can modify the `Path` option in the middleware configuration. Note that changing the path requires updating both the middleware setup and any references in your documentation.

#### How It Works

The generated server embeds the OpenAPI 2.0 specification as an embedded asset. When you visit `/docs`:

1. The server serves the embedded ReDoc or SwaggerUI assets
2. The UI loads the swagger specification from `/swagger.json`
3. The API documentation is rendered interactively

This approach ensures that:
- The API documentation is always in sync with the generated code
- No additional servers or build steps are required
- The documentation is available in any environment where the server runs

#### Common Use Cases

**Enable/disable API Browser in production:**
```go
func configureAPI(api *operations.MyAppAPI) http.Handler {
    // ... other configuration ...
    
    // Conditionally serve docs based on environment
    if os.Getenv("ENABLE_API_DOCS") != "false" {
        return setupMiddlewares(api.Serve())
    }
    return api.Serve()
}
```

**Require authentication for API docs:**
```go
func setupMiddlewares(handler http.Handler) http.Handler {
    return authMiddleware(handler) // your auth middleware
}
```

Originally from issue [#2401](https://github.com/go-swagger/go-swagger/issues/2401).

### Serving swagger-ui with the API Server

_Use-Case_: I was trying to serve swagger-ui from the generated API Server and
I didn't find a straightforward enough way in the docs,
so I've created my own swagger-ui middleware:

```golang
func setupGlobalMiddleware(handler http.Handler) http.Handler {
    return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Shortcut helpers for swagger-ui
        if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
            http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
            return
        }
        // Serving ./swagger-ui/
        if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
            http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
            return
        }
        handler.ServeHTTP(w, r)
    })
}
```

*Can this be achieved any easier?*

**Answer**: I think that's the way to do it.

At some point I included an embedded version of the swagger-ui in this toolkit but it became annoying to keep up to date
and severely bloated the size of my binary.

*What do you say if we add swagger-ui as git submodule, include this middleware in your binary and update the docs?*

I'm reluctant to do so at this point in time because a git submodule break go-gettability.

>I've had it included at one point but it's so much of a moving target that it would always be outdated.
>On top of it is a lot of javascript and html files and people haven't been over the moon when go-swagger gets
>vendored and they see all of that.

Originally from issue [#370](https://github.com/go-swagger/go-swagger/issues/370).

See also: How to serve Swagger UI from a preexisting web app? [#1029](https://github.com/go-swagger/go-swagger/issues/1029).

### How to serve Swagger UI from a preexisting web app?
_Use-Case_: Does go-swagger provide an `http.HandlerFunc` or other easy method for serving Swagger UI from a preexisting web app? 
I want my web app to expose `/swagger-ui`, without using code generation, and without hosting a separate server.

**Answer**: there are a few ways you can serve a UI.

Use the middleware provided in the go-openapi/runtime package: https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go

Originally from issues [#1029](https://github.com/go-swagger/go-swagger/issues/1029) and [#976](https://github.com/go-swagger/go-swagger/issues/976)

### How to use swagger-ui cors?

**Answer**: you can add a cors middleware.

Like: https://github.com/rs/cors

[Documentation on how to customize middleware](reference/middleware.md)

Working example (in `configure_name.go`):

```golang
import "github.com/rs/cors"

func setupGlobalMiddleware(handler http.Handler) http.Handler {
    handleCORS := cors.Default().Handler
    return handleCORS(handler)
}
```

Originally from issue [#481](https://github.com/go-swagger/go-swagger/issues/481).

### How to serve my UI files?
_Use-Case_: I generated server code using go-swagger with my swagger.yaml file like below.
```bash
$ swagger generate server --exclude-main -A myapp -t gen -f ./swagger.yaml
```
And I want to add new handler to serve my own UI files. 
In this case, is middleware only solution to serve UI files? Or can I add new handler to serve files without middleware?
```go
// Handler example which I want to add to swagger server
func pastaWorkspacePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/workspace.html")
}
```

I solved the problem using middleware.
```go
func FileServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
		} else {
			http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
		}
	})
}
```
But I'm not sure it is the best solution.

**Hint:** more info on using middlewares is found here: https://goswagger.io/use/middleware.html
That page also contains a link to a good explanation on how to create net/http middlewares.

> An implementation example is provided by the go-swagger serve UI command. It constructs a server with a redoc middleware:
> https://github.com/go-swagger/go-swagger/blob/f552963ac0dfdec0450f6749aeeeeb2d31cd5544/cmd/swagger/commands/serve.go#L35.

Originally from issue [#1375](https://github.com/go-swagger/go-swagger/issues/1375).
