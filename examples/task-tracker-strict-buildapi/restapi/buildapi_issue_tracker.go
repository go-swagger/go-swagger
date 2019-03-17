// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"io"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/go-swagger/go-swagger/examples/task-tracker-strict-buildapi/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/task-tracker-strict-buildapi/restapi/operations/tasks"
)

func BuildAPI(
	ServeError func(http.ResponseWriter, *http.Request, error),
	Logger func(string, ...interface{}),

	JSONConsumer func(r io.Reader, target interface{}) error,
	MultipartformConsumer func(r io.Reader, target interface{}) error,

	JSONProducer func(w io.Writer, data interface{}) error,

	APIKeyAuth func(token string) (interface{}, error),
	TokenHeaderAuth func(token string) (interface{}, error),
	APIAuthorizer runtime.Authorizer,

	AddCommentToTask func(params tasks.AddCommentToTaskParams, principal interface{}) tasks.AddCommentToTaskResponder,
	CreateTask func(params tasks.CreateTaskParams, principal interface{}) tasks.CreateTaskResponder,
	DeleteTask func(params tasks.DeleteTaskParams, principal interface{}) tasks.DeleteTaskResponder,
	GetTaskComments func(params tasks.GetTaskCommentsParams) tasks.GetTaskCommentsResponder,
	GetTaskDetails func(params tasks.GetTaskDetailsParams) tasks.GetTaskDetailsResponder,
	ListTasks func(params tasks.ListTasksParams) tasks.ListTasksResponder,
	UpdateTask func(params tasks.UpdateTaskParams, principal interface{}) tasks.UpdateTaskResponder,
	UploadTaskFile func(params tasks.UploadTaskFileParams, principal interface{}) tasks.UploadTaskFileResponder,

	ServerShutdown func(),

) *operations.IssueTrackerAPI {
	api := &operations.IssueTrackerAPI{}

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "token" query is set
	api.APIKeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (api_key) token from query param [token] has not yet been implemented")
	}
	// Applies when the "X-Token" header is set
	api.TokenHeaderAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (token_header) X-Token from header param [X-Token] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.TasksAddCommentToTaskHandler = tasks.AddCommentToTaskHandlerFunc(func(params tasks.AddCommentToTaskParams, principal interface{}) tasks.AddCommentToTaskResponder {
		return tasks.AddCommentToTaskNotImplemented()
	})
	api.TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal interface{}) tasks.CreateTaskResponder {
		return tasks.CreateTaskNotImplemented()
	})
	api.TasksDeleteTaskHandler = tasks.DeleteTaskHandlerFunc(func(params tasks.DeleteTaskParams, principal interface{}) tasks.DeleteTaskResponder {
		return tasks.DeleteTaskNotImplemented()
	})
	api.TasksGetTaskCommentsHandler = tasks.GetTaskCommentsHandlerFunc(func(params tasks.GetTaskCommentsParams) tasks.GetTaskCommentsResponder {
		return tasks.GetTaskCommentsNotImplemented()
	})
	api.TasksGetTaskDetailsHandler = tasks.GetTaskDetailsHandlerFunc(func(params tasks.GetTaskDetailsParams) tasks.GetTaskDetailsResponder {
		return tasks.GetTaskDetailsNotImplemented()
	})
	api.TasksListTasksHandler = tasks.ListTasksHandlerFunc(func(params tasks.ListTasksParams) tasks.ListTasksResponder {
		return tasks.ListTasksNotImplemented()
	})
	api.TasksUpdateTaskHandler = tasks.UpdateTaskHandlerFunc(func(params tasks.UpdateTaskParams, principal interface{}) tasks.UpdateTaskResponder {
		return tasks.UpdateTaskNotImplemented()
	})
	api.TasksUploadTaskFileHandler = tasks.UploadTaskFileHandlerFunc(func(params tasks.UploadTaskFileParams, principal interface{}) tasks.UploadTaskFileResponder {
		return tasks.UploadTaskFileNotImplemented()
	})

	api.ServerShutdown = func() {}

	if ServeError != nil {
		api.ServeError = errors.ServeError
	}

	if Logger != nil {
		api.Logger = Logger
	}

	if JSONConsumer != nil {
		api.JSONConsumer = runtime.ConsumerFunc(JSONConsumer)
	}

	if MultipartformConsumer != nil {
		api.MultipartformConsumer = runtime.ConsumerFunc(MultipartformConsumer)
	}

	if JSONProducer != nil {
		api.JSONProducer = runtime.ProducerFunc(JSONProducer)
	}

	if APIKeyAuth != nil {
		api.APIKeyAuth = APIKeyAuth
	}
	if TokenHeaderAuth != nil {
		api.TokenHeaderAuth = TokenHeaderAuth
	}
	if APIAuthorizer != nil {
		api.APIAuthorizer = APIAuthorizer
	}

	if AddCommentToTask != nil {
		api.TasksAddCommentToTaskHandler = tasks.AddCommentToTaskHandlerFunc(AddCommentToTask)
	}

	if CreateTask != nil {
		api.TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(CreateTask)
	}

	if DeleteTask != nil {
		api.TasksDeleteTaskHandler = tasks.DeleteTaskHandlerFunc(DeleteTask)
	}

	if GetTaskComments != nil {
		api.TasksGetTaskCommentsHandler = tasks.GetTaskCommentsHandlerFunc(GetTaskComments)
	}

	if GetTaskDetails != nil {
		api.TasksGetTaskDetailsHandler = tasks.GetTaskDetailsHandlerFunc(GetTaskDetails)
	}

	if ListTasks != nil {
		api.TasksListTasksHandler = tasks.ListTasksHandlerFunc(ListTasks)
	}

	if UpdateTask != nil {
		api.TasksUpdateTaskHandler = tasks.UpdateTaskHandlerFunc(UpdateTask)
	}

	if UploadTaskFile != nil {
		api.TasksUploadTaskFileHandler = tasks.UploadTaskFileHandlerFunc(UploadTaskFile)
	}

	if ServerShutdown != nil {
		api.ServerShutdown = ServerShutdown
	}

	return api
}
