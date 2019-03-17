# To do list example

Shows a fully loaded server that by default listens on a unix, http and https socket.

You can run just the http listener for quick testing:

```shellsession
go run ./cmd/todo-list-server/main.go --scheme http
```

## Run full server

To run the full server you need to build the binary and run it with sudo enabled.

```shellsession
go build ./cmd/todo-list-server
sudo ./todo-list-server --tls-certificate mycert1.crt --tls-certificate-key mycert1.key
```
