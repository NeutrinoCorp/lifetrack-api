// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"github.com/google/wire"
	"github.com/neutrinocorp/lifetrack-api/internal/application/category"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/event"
	"github.com/neutrinocorp/lifetrack-api/internal/domain/repository"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/configuration"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/lifetrack-api/internal/infrastructure/persistence/inmemcategory"
)

// Injectors from wire.go:

func InjectAddCategoryHandler() (*category.AddCommandHandler, func(), error) {
	inMemory := inmemcategory.NewInMemory()
	configurationConfiguration := configuration.NewConfiguration()
	bus := provideEventBus(configurationConfiguration)
	addCommandHandler := category.NewAddCommandHandler(inMemory, bus)
	return addCommandHandler, func() {
	}, nil
}

func InjectGetCategoryQuery() (*category.GetQuery, func(), error) {
	inMemory := inmemcategory.NewInMemory()
	getQuery := category.NewGetQuery(inMemory)
	return getQuery, func() {
	}, nil
}

func InjectListCategoriesQuery() (*category.ListQuery, func(), error) {
	inMemory := inmemcategory.NewInMemory()
	listQuery := category.NewListQuery(inMemory)
	return listQuery, func() {
	}, nil
}

func InjectEditCategory() (*category.UpdateCommandHandler, func(), error) {
	inMemory := inmemcategory.NewInMemory()
	configurationConfiguration := configuration.NewConfiguration()
	bus := provideEventBus(configurationConfiguration)
	updateCommandHandler := category.NewUpdateCommandHandler(inMemory, bus)
	return updateCommandHandler, func() {
	}, nil
}

func InjectRemoveCategory() (*category.RemoveCommandHandler, func(), error) {
	inMemory := inmemcategory.NewInMemory()
	configurationConfiguration := configuration.NewConfiguration()
	bus := provideEventBus(configurationConfiguration)
	removeCommandHandler := category.NewRemoveCommandHandler(inMemory, bus)
	return removeCommandHandler, func() {
	}, nil
}

// wire.go:

var infraSet = wire.NewSet(configuration.NewConfiguration, wire.Bind(new(repository.Category), new(*inmemcategory.InMemory)), inmemcategory.NewInMemory, provideEventBus, eventbus.NewInMemory)

func provideEventBus(cfg configuration.Configuration) event.Bus {
	return eventbus.NewInMemory(cfg)
}
