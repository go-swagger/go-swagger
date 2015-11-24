+++
Categories = ["untyped", "server"]
Tags = []
date = "2015-11-22T23:21:52-08:00"
title = "Untyped API server"
weight = 5
series = ["home"]
+++

The toolkit supports serving a swagger spec with untyped data. This means that it uses mostly interface{} as params to
each operation and as result type. It does allow you to serve a spec up quickly. This is one of the building blocks
required to serve up stub API's and to generate a test server with predictable responses.

<!--more-->
