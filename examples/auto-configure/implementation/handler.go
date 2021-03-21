package implementation

import (
	"sync"

	"github.com/go-swagger/go-swagger/examples/auto-configure/models"
)

// HandlerImpl implements all required configuration and api handling
// functionalities for todo list server backend
type HandlerImpl struct {
	TodosHandlerImpl
	ConfigureImpl
	AuthImpl
}

func New() *HandlerImpl {
	return &HandlerImpl{
		TodosHandlerImpl: TodosHandlerImpl{
			lock:  sync.Mutex{},
			items: make(map[int64]*models.Item),
			idx:   0,
		},
		ConfigureImpl: ConfigureImpl{
			flags: Flags{},
		},
		AuthImpl: AuthImpl{},
	}
}
