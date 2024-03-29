---
title: swagger serve
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 60
---
# Serve a documentation site

The toolkit has a command to serve a swaggr specificaion as a JSON document and optionally a UI to render this spec.

This server publishes UI widgets `redoc` or `swaggerUI`.
These assets are downloaded from public js repos, so you need an internet connection to use them.

### Usage

To serve a documentation site:

```
Usage:
  swagger [OPTIONS] serve [serve-OPTIONS] {specification file}

serve a spec and swagger or redoc documentation ui

Application Options:
  -q, --quiet                         silence logs
      --log-output=LOG-FILE           redirect logs to file

Help Options:
  -h, --help                          Show this help message

[serve command options]
          --base-path=                the base path to serve the spec and UI at
      -F, --flavor=[redoc|swagger]    the flavor of docs, can be swagger or redoc (default: redoc)
          --doc-url=                  override the url which takes a url query param to render the doc ui
          --no-open                   when present won't open the browser to show the url
          --no-ui                     when present, only the swagger spec will be served
      -p, --port=                     the port to serve this site [$PORT]
          --host=                     the interface to serve this site, defaults to 0.0.0.0 [$HOST]
```

This will start a server with CORS enabled so that sites on other domains can load your specification document.

> **Attention**: if you use external $ref to a local file, we recommend that you serve a flattened specification document,
> as the server won't serve external references.

Example:
```sh
cd examples/cli
swagger serve swagger.yml
```

### Flavors

At this moment the UI can be served into 2 flavors.

#### Redoc

The swagger source code has a middleware for embedding Redoc.
So for the redoc flavor we make use of that and use it with the spec you have on disk.

We use the redoc JS bundle hosted at https://cdn.jsdelivr.net/npm/redoc/bundles/redoc.standalone.js

#### Swagger UI

For the swagger flavor we use the UI bundle hostsed at https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js

#### Your own UI

You can use your own UI by pointing it to the spec served by this command.

When no UI is being served, the command will print the URL to the spec document to the terminal.

You can also use the `--doc-url` to provide another url as base, for example
the url to your documentation site, which would need to recognize the query param URL to load the swagger spec from,
through the browser.

### More

There are some more options for this command which you can view with:

```
swagger serve --help
```
