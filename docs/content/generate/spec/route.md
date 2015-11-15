+++
categories = ["spec", "generate"]
tags = []
date = "2015-11-14T20:10:39-08:00"
title = "swagger:route"

+++

A **swagger:route** annotation links a path to a method.
This operation gets a unique id, which is used in various places as method name.
One such usage is in method names for client generation for example.

Because there are many routers available, this tool does not try to parse the paths
you provided to your routing library of choice. So you have to specify your path pattern
yourself in valid swagger syntax.

```
swagger:route [method] [path pattern] [operation id] [?tag1 tag2 tag3]
```
