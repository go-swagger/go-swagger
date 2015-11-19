// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/go-swagger/go-swagger/examples/generated/models"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/pet"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/store"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/user"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.PetStoreAPI) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.XMLConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		return errors.NotImplemented("xml consumer has not yet been implemented")
	})
	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.XMLProducer = httpkit.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("xml producer has not yet been implemented")
	})

	api.APIKeyAuth = func(token string) (*models.User, error) {
		return nil, errors.NotImplemented("api key auth api_key from header has not yet been implemented")
	}

	api.GetPetByIDHandler = pet.GetPetByIDHandlerFunc(func(params pet.GetPetByIDParams, principal *models.User) (*models.Pet, error) {
		return nil, errors.NotImplemented("operation getPetById has not yet been implemented")
	})
	api.LogoutUserHandler = user.LogoutUserHandlerFunc(func() error {
		return errors.NotImplemented("operation logoutUser has not yet been implemented")
	})
	api.DeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) error {
		return errors.NotImplemented("operation deleteUser has not yet been implemented")
	})
	api.DeleteOrderHandler = store.DeleteOrderHandlerFunc(func(params store.DeleteOrderParams) error {
		return errors.NotImplemented("operation deleteOrder has not yet been implemented")
	})
	api.UpdateUserHandler = user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) error {
		return errors.NotImplemented("operation updateUser has not yet been implemented")
	})
	api.LoginUserHandler = user.LoginUserHandlerFunc(func(params user.LoginUserParams) (string, error) {
		return "", errors.NotImplemented("operation loginUser has not yet been implemented")
	})
	api.DeletePetHandler = pet.DeletePetHandlerFunc(func(params pet.DeletePetParams, principal *models.User) error {
		return errors.NotImplemented("operation deletePet has not yet been implemented")
	})
	api.UpdatePetHandler = pet.UpdatePetHandlerFunc(func(params pet.UpdatePetParams, principal *models.User) error {
		return errors.NotImplemented("operation updatePet has not yet been implemented")
	})
	api.CreateUsersWithListInputHandler = user.CreateUsersWithListInputHandlerFunc(func(params user.CreateUsersWithListInputParams) error {
		return errors.NotImplemented("operation createUsersWithListInput has not yet been implemented")
	})
	api.AddPetHandler = pet.AddPetHandlerFunc(func(params pet.AddPetParams, principal *models.User) error {
		return errors.NotImplemented("operation addPet has not yet been implemented")
	})
	api.CreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) error {
		return errors.NotImplemented("operation createUser has not yet been implemented")
	})
	api.FindPetsByStatusHandler = pet.FindPetsByStatusHandlerFunc(func(params pet.FindPetsByStatusParams, principal *models.User) ([]models.Pet, error) {
		return nil, errors.NotImplemented("operation findPetsByStatus has not yet been implemented")
	})
	api.GetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) (*models.User, error) {
		return nil, errors.NotImplemented("operation getUserByName has not yet been implemented")
	})
	api.PlaceOrderHandler = store.PlaceOrderHandlerFunc(func(params store.PlaceOrderParams) (*models.Order, error) {
		return nil, errors.NotImplemented("operation placeOrder has not yet been implemented")
	})
	api.UpdatePetWithFormHandler = pet.UpdatePetWithFormHandlerFunc(func(params pet.UpdatePetWithFormParams, principal *models.User) error {
		return errors.NotImplemented("operation updatePetWithForm has not yet been implemented")
	})
	api.GetOrderByIDHandler = store.GetOrderByIDHandlerFunc(func(params store.GetOrderByIDParams) (*models.Order, error) {
		return nil, errors.NotImplemented("operation getOrderById has not yet been implemented")
	})
	api.FindPetsByTagsHandler = pet.FindPetsByTagsHandlerFunc(func(params pet.FindPetsByTagsParams, principal *models.User) ([]models.Pet, error) {
		return nil, errors.NotImplemented("operation findPetsByTags has not yet been implemented")
	})
	api.CreateUsersWithArrayInputHandler = user.CreateUsersWithArrayInputHandlerFunc(func(params user.CreateUsersWithArrayInputParams) error {
		return errors.NotImplemented("operation createUsersWithArrayInput has not yet been implemented")
	})

}
