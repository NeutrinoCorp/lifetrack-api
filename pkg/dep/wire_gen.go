// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/wire"
	"github.com/neutrinocorp/life-track-api/internal/application/command"
	"github.com/neutrinocorp/life-track-api/internal/application/query"
	"github.com/neutrinocorp/life-track-api/internal/domain/event"
	"github.com/neutrinocorp/life-track-api/internal/domain/repository"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InjectAddCategoryHandler() (*command.AddCategoryHandler, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	aws := eventbus.NewAWS(session, configuration)
	addCategoryHandler := command.NewAddCategoryHandler(category, aws)
	return addCategoryHandler, func() {
		cleanup()
	}, nil
}

func InjectGetCategoryQuery() (*query.GetCategory, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	getCategory := query.NewGetCategory(category)
	return getCategory, func() {
		cleanup()
	}, nil
}

func InjectListCategoriesQuery() (*query.ListCategories, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	listCategories := query.NewListCategories(category)
	return listCategories, func() {
		cleanup()
	}, nil
}

func InjectChangeCategoryState() (*command.ChangeCategoryStateHandler, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	aws := eventbus.NewAWS(session, configuration)
	changeCategoryStateHandler := command.NewChangeCategoryStateHandler(category, aws)
	return changeCategoryStateHandler, func() {
		cleanup()
	}, nil
}

func InjectEditCategory() (*command.EditCategoryHandler, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	aws := eventbus.NewAWS(session, configuration)
	editCategoryHandler := command.NewEditCategoryHandler(category, aws)
	return editCategoryHandler, func() {
		cleanup()
	}, nil
}

func InjectRemoveCategory() (*command.RemoveCategoryHandler, func(), error) {
	session := infrastructure.NewSession()
	configuration, err := infrastructure.NewConfiguration()
	if err != nil {
		return nil, nil, err
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		return nil, nil, err
	}
	category := provideCategoryRepository(session, configuration, logger)
	aws := eventbus.NewAWS(session, configuration)
	removeCategoryHandler := command.NewRemoveCategoryHandler(category, aws)
	return removeCategoryHandler, func() {
		cleanup()
	}, nil
}

// wire.go:

var infraSet = wire.NewSet(infrastructure.NewConfiguration, infrastructure.NewSession, logging.NewZapProd, provideCategoryRepository, wire.Bind(new(event.Bus), new(*eventbus.AWS)), eventbus.NewAWS)

func provideCategoryRepository(s *session.Session, cfg infrastructure.Configuration, logger *zap.Logger) repository.Category {
	return persistence.NewCategory(persistence.NewCategoryDynamoRepository(s, cfg), logger)
}
