package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/middleware/untyped"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// our spec doesn't have application/json in the consumes or produces
	// so we need to clear those settings out
	api := untyped.NewAPI(specDoc).WithoutJSONDefaults()

	// register serializers
	mediaType := "application/io.goswagger.examples.todo-list.v1+json"
	api.DefaultConsumes = mediaType
	api.DefaultProduces = mediaType
	api.RegisterConsumer(mediaType, runtime.JSONConsumer())
	api.RegisterProducer(mediaType, runtime.JSONProducer())

	// register the operation handlers
	api.RegisterOperation("GET", "/", findTodos)
	api.RegisterOperation("POST", "/", addOne)
	api.RegisterOperation("PUT", "/{id}", updateOne)
	api.RegisterOperation("DELETE", "/{id}", destroyOne)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}

	// construct the application context for this server
	// use the loaded spec document and the api descriptor with the default router
	app := middleware.NewContext(specDoc, api, nil)

	log.Println("serving", specDoc.Spec().Info.Title, "at http://localhost:8000")

	// serve the api with spec and UI
	if err := http.ListenAndServe(":8000", app.APIHandler(nil)); err != nil {
		log.Fatalln(err)
	}
}

var findTodos = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'findTodos'")
	log.Printf("%#v\n", params)

	return items, nil
})

var addOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'addOne'")
	log.Printf("%#v\n", params)

	body := params.(map[string]interface{})["body"].(map[string]interface{})
	addItem(body)
	return body, nil
})

var updateOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'updateOne'")
	log.Printf("%#v\n", params)

	data := params.(map[string]interface{})
	id := data["id"].(int64)
	body := data["body"].(map[string]interface{})
	return updateItem(id, body)
})

var destroyOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'destroyOne'")
	log.Printf("%#v\n", params)

	removeItem(params.(map[string]interface{})["id"].(int64))
	return nil, nil
})

var items = []map[string]interface{}{
	{"id": int64(1), "description": "feed dog", "completed": true},
	{"id": int64(2), "description": "feed cat"},
}

var itemsLock = &sync.Mutex{}
var lastItemID int64 = 2

func newItemID() int64 {
	return atomic.AddInt64(&lastItemID, 1)
}

func addItem(item map[string]interface{}) {
	itemsLock.Lock()
	defer itemsLock.Unlock()
	item["id"] = newItemID()
	items = append(items, item)
}

func updateItem(id int64, body map[string]interface{}) (map[string]interface{}, error) {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	item, err := itemByID(id)
	if err != nil {
		return nil, err
	}
	delete(body, "id")
	for k, v := range body {
		item[k] = v
	}
	return item, nil
}

func removeItem(id int64) {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	var newItems []map[string]interface{}
	for _, item := range items {
		if item["id"].(int64) != id {
			newItems = append(newItems, item)
		}
	}
	items = newItems
}

func itemByID(id int64) (map[string]interface{}, error) {
	for _, item := range items {
		if item["id"].(int64) == id {
			return item, nil
		}
	}
	return nil, errors.NotFound("not found: item %d", id)
}
