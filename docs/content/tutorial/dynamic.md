+++
Categories = ["tutorial", "spec"]
Tags = []
date = "2015-11-22T23:21:52-08:00"
title = "Dynamic API definition"
weight = 8
series = ["tutorials", "home"]
+++

The toolkit supports building a swagger specification entirely with go code. It does allow you to serve a spec up quickly. This is one of the building blocks required to serve up stub API's and to generate a test server with predictable responses, however this is not as bad as it sounds...

<!--more-->

This tutorial uses the todo list application to serve a swagger based API defined entirely in go code.
Because we know what we want the spec to look like, first we'll just build the entire spec with the interal dsl.

```go
doc := spec.NewSwagger("")
doc
```

Now that we
