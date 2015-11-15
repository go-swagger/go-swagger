+++
categories = ["spec", "generate"]
tags = []
date = "2015-11-14T20:10:32-08:00"
title = "swagger:model"

+++

A **swagger:model** annotation optionally gets a model name as extra data on the line.
when this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

```
swagger:model [?model name]
```
