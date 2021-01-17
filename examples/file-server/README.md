# File upload server

This example demonstrates how to build a simple file upload endpoint
with swagger and go-swagger.

## Try it

1. Build the server

```
cd restapi/cmd/file-upload-server
go build

./file-upload-server --port 8000
2021/01/17 18:54:09 Serving file upload at http://127.0.0.1:8000
```

2. Run the client

From another terminal:

```
go run upload_file.go swagger.yml
```

Logs on the server:
```
2021/01/17 18:54:15 received file name: swagger.yml
2021/01/17 18:54:15 received file size: 512
2021/01/17 18:54:15 copied bytes 512
2021/01/17 18:54:15 file uploaded copied as upload427417421/uploaded_file_0.dat
```

The file has been copied in a temporary folder `cmd/file-upload-server/upload*/`


## Specification

We use the swagger type `file` in a multipart form, like so:

```yaml
paths:
  /upload:
    post:
      consumes:
      - multipart/form-data
      parameters:
      - name: file
        in: formData
        type: file
```

## Server side

The handler receives a `io.ReadCloser` as the file to consume.

Under the hood, the runtime builds this with a `*runtime.File`, which provides access to some header information, such as:

```go
		if namedFile, ok := params.File.(*runtime.File); ok {
			log.Printf("received file name: %s", namedFile.Header.Filename)
			log.Printf("received file size: %d", namedFile.Header.Size)
		}
```

## Client side

The local file is handled as a `runtime.NamedReadCloser` (that is, a `io.ReadCloser` plus the `Name() string` method).
A regular `os.File` satisfies this.

The file can be passed directly to the client method, like so:

```go
	params := uploads.NewUploadFileParams().WithFile(reader)

	_, err := uploader.Uploads.UploadFile(params)
```
