# Serve a documentation site

The toolkit has a command to serve a spec json document and optionally a UI for a given spec. 
It embeds redoc, so you can use that documentation site without an internet connection.
Or it can load your local spec into the swagger docs viewer at http://petstore.swagger.io

<!--more-->

### Usage 

To serve a documentation site:

```
swagger serve [http-url|filepath]
```

This will start a server with cors enabled so that sites on other domains can load your specification document. 

### Flavors

At this moment the UI can be served into 2 flavors.

#### Redoc

The swagger source code has a middleware for embedding Redoc.
So for the redoc flavor we make use of that and use it with the spec you have on disk.

#### Swagger UI

For the swagger flavor we use the UI hosted at http://petstore.swagger.io.
The server has CORS enabled and appends the url for the spec JSON to the petstore url as a query string. 

#### Your own UI

You can use your own UI by pointing it to the spec served by this command.
When no ui is being served, the terminal will print the url to the spec document.
You can also use the `--doc-url` to provide another url as base. 
The url to your documentation site for example, which would need to recognize the query param url to load the swagger spec from, through the browser.

### More

There are some more options for this command which you can view with:

```
swagger serve --help
```
