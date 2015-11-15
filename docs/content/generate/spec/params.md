+++
categories = ["spec"]
tags = []
date = "2015-11-14T20:10:43-08:00"
title = "swagger:params"

weight = 20
+++

Links a struct to one or more operations. The params in the resulting swagger spec can be composed of several structs.
There are no guarantees given on how property name overlaps are resolved when several structs apply to the same operation.
This tag works very similar to the swagger:model tag except that it produces valid parameter objects instead of schema
objects.

```
swagger:params [operationid1 operationid2]
```
