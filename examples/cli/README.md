# Generate CLI Client Tool
Generate a command line client tool for your server
(This is in alpha state and is under development)

* Generated CLI code is a wrapper of the generated client code, which reads command line options and args to construct appropriate parameters and send to the server.
* Based on [cobra framework](https://github.com/spf13/cobra).
* Support shell completions (not yet implemented), based on [cobra framework](https://github.com/spf13/cobra).
## General Command Layout
* Root command manages global flags
* Each open-api tag is a sub-command under root command. In go-swagger it is called operation group.
* Each open-api operationId is a sub-command under the tag it belongs to.
    * Each path and query parameter is a command line flag.
    * Body parameter corresponds to a command line flag, as a json string.
    * Each field in body parameter is a command line flag, and this extends to sub-definitions recursively.
        * body parameter json string will be taken as base payload, which flags for body fields will overwrite.

# Todo List Example
CLI tool in this folder is generated using the same swagger.yaml as `examples/auto-configure`. We will run that server to test this cli executable.

## Get Started
### Generate the code
Generate go-swagger command:
```
$ go run cmd/swagger/swagger.go generate cli --target=examples/cli --spec=examples/cli/swagger.yml
```
Or simply:
```
$ swagger generate cli --target=examples/cli --spec=examples/cli/swagger.yml
```
### Run the tool using go run
```
$ go run examples/cli/cmd/cli/main.go --debug --hostname localhost:12345 --x-todolist-token "example token" todos addOne --item.description "hi" --body "{}"
```
### Build the executable and then run
```
$ go build -o examples/cli/cmd/cli/todo-cli examples/cli/cmd/cli/main.go
$ ./examples/cli/cmd/cli/todo-cli
```
### Help Message Example
* Root help message will be of the following format:
```
$ ./examples/cli/cmd/cli/todo-cli --help
Usage:
  AToDoListApplication [command]

Available Commands:
  help        Help about any command
  todos       

Flags:
      --debug                     output debug logs
  -h, --help                      help for AToDoListApplication
      --hostname string           hostname of the service (default "localhost")
      --scheme string             Choose from: [http] (default "http")
      --x-todolist-token string    (default "none")

Use "AToDoListApplication [command] --help" for more information about a command.
```
* Command help message example:
```
$ ./examples/cli/cmd/cli/todo-cli todos updateOne --help
Usage:
  AToDoListApplication todos updateOne [flags]

Flags:
      --body string               Optional json string for body of form {"description":null}.
  -h, --help                      help for updateOne
      --id int                    Required. 
      --item.completed            
      --item.description string   Required. 
      --item.id int               ReadOnly.

Global Flags:
      --debug                     output debug logs
      --hostname string           hostname of the service (default "localhost")
      --scheme string             Choose from: [http] (default "http")
      --x-todolist-token string    (default "none")
```
### Running CLI against auto-configure server
* Start the server
```
$ go run examples/auto-configure/cmd/a-to-do-list-application-server/main.go --port=12345
```
* Make request using the CLI tool (using another shell)
```
$ go run examples/cli/cmd/cli/main.go --hostname localhost:12345 --x-todolist-token "example token" todos addOne --item.description "hi" --body "{}"
{"description":"hi"}

$ go run examples/cli/cmd/cli/main.go --hostname localhost:12345 --x-todolist-token "example token" todos findTodos
[{"description":"hi"}]

$ go run examples/cli/cmd/cli/main.go --hostname localhost:12345 --x-todolist-token "example token" todos updateOne --id 1 --item.completed true --item.description "done"
{"code":404,"message":"Item with id 1 is not found"}
```

### Missing Features
* `AllOf`
* Multi success response
* Array and maps in body (It is unclear how to support)
* Shell auto-complete
* Enums
    * In help message
    * In shell auto-complete
* Validate params before sending
* Read host, schemes and auth from config using [viper](https://github.com/spf13/viper) 