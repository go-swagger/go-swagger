---
title: swagger generate
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 35 
---
# Generating

The toolkit has a command to generate different targets.

### Usage

Available subcommands.

```
Usage:
  swagger [OPTIONS] generate <command>

generate go code for the swagger spec file

Application Options:
  -q, --quiet                  silence logs
      --log-output=LOG-FILE    redirect logs to file

Help Options:
  -h, --help                   Show this help message

Available commands:
  cli        generate a command line client tool from the swagger spec
  client     generate all the files for a client library
  markdown   generate a markdown representation from the swagger spec
  model      generate one or more models from the swagger spec
  operation  generate one or more server operations from the swagger spec
  server     generate all the files for a server application
  spec       generate a swagger spec document from a go application
  support    generate supporting files like the main function and the api builder
```

For code generation targets (`cli`, `client`, `model`, `operaion`, `server`, `support`), read more [here](../generate/).

For spec generation targets (`spec`), read more [there](../generate-spec/).

For markdown generation target (`markdown`), read [this](markdown.md).
