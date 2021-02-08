# petstore

This minimalist example demonstrates the use of the go-openapi runtime
as an "untyped" server, without any generated code.

Usage:

```bash
cd server
go build

./server &
2020/12/18 11:49:13 Serving petstore api on http://127.0.0.1:8344/api/

curl http://127.0.0.1:8344/api/pets/
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```
