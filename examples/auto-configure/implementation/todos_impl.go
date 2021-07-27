package implementation

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/examples/auto-configure/models"
	"github.com/go-swagger/go-swagger/examples/auto-configure/restapi/operations/todos"
)

type TodosHandlerImpl struct {
	// locks the entire struct
	lock  sync.Mutex
	items map[int64]*models.Item
	// counter for newest item
	idx int64
}

func (i *TodosHandlerImpl) AddOne(params todos.AddOneParams, principal interface{}) middleware.Responder {
	i.lock.Lock()
	defer i.lock.Unlock()
	newItem := params.Body
	if newItem == nil {
		return todos.NewAddOneDefault(http.StatusBadRequest).
			WithPayload(&models.Error{
				Code:    http.StatusBadRequest,
				Message: &[]string{"Item Body is nil"}[0],
			})
	}
	// assign new id
	newItem.ID = i.idx
	i.idx++

	i.items[newItem.ID] = newItem
	return todos.NewAddOneCreated().WithPayload(newItem)
}

func (i *TodosHandlerImpl) DestroyOne(params todos.DestroyOneParams, principal interface{}) middleware.Responder {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.items[params.ID]; !ok {
		return todos.NewDestroyOneDefault(http.StatusNotFound)
	}
	delete(i.items, params.ID)
	return todos.NewDestroyOneNoContent()
}

func (i *TodosHandlerImpl) FindTodos(params todos.FindTodosParams, principal interface{}) middleware.Responder {
	i.lock.Lock()
	defer i.lock.Unlock()
	mergedParams := todos.NewFindTodosParams()
	mergedParams.Since = swag.Int64(0)
	if params.Since != nil {
		mergedParams.Since = params.Since
	}
	if params.Limit != nil {
		mergedParams.Limit = params.Limit
	}
	limit := *mergedParams.Limit
	since := *mergedParams.Since
	// copy all items and return
	result := make([]*models.Item, 0)
	for id, item := range i.items {
		if len(result) >= int(limit) {
			break
		}
		if since == 0 || id > since {
			result = append(result, item)
		}
	}
	return todos.NewFindTodosOK().WithPayload(result)
}

func (i *TodosHandlerImpl) UpdateOne(params todos.UpdateOneParams, principal interface{}) middleware.Responder {
	i.lock.Lock()
	defer i.lock.Unlock()

	if _, ok := i.items[params.ID]; !ok {
		errStr := fmt.Sprintf("Item with id %v is not found", params.ID)
		return todos.NewUpdateOneDefault(http.StatusNotFound).
			WithPayload(&models.Error{
				Code:    http.StatusNotFound,
				Message: &errStr,
			})
	}
	params.Body.ID = params.ID
	i.items[params.ID] = params.Body

	return todos.NewUpdateOneOK().WithPayload(params.Body)
}
