<!-- Questions about the serve UI use-case -->
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
>On top of it it are a lot of javascript and html files and people haven't been over the moon when go-swagger gets
>vendored and they see all of that.

Originally from issue [#370](https://github.com/go-swagger/go-swagger/issues/370).

See also: How to serve Swagger UI from a preexisting web app? [#1029](https://github.com/go-swagger/go-swagger/issues/1029).

-------------------

Back to [all contributions](README.md#all-contributed-questions)