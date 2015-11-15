+++
categories = ["spec", "generate"]
tags = []
date = "2015-11-14T20:10:52-08:00"
title = "swagger:response"

+++

Reads a struct decorated with **swagger:response** and uses that information to fill up the headers and the schema for a response.
A swagger:route can specify a response name for a status code and then the matching response will be used for that operation in the swagger definition.

```
swagger:response [?response name]
```
