package implementation

// HandlerImpl implements all required configuration and api handling
// functionalities for todo list server backend
type HandlerImpl struct {
	TodosHandlerImpl
	ConfigureImpl
	AuthImpl
}

func New() *HandlerImpl {
	return &HandlerImpl{
		TodosHandlerImpl: TodosHandlerImpl{},
		ConfigureImpl: ConfigureImpl{
			flags: Flags{},
		},
		AuthImpl: AuthImpl{},
	}
}
