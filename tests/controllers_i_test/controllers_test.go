package controllers_i_test

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/repository"
	"boilerplate-api/api/services"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	"boilerplate-api/tests"
	"go.uber.org/fx"
	"os"
	"reflect"
	"testing"
)

type ControllerTest interface {
	SetupControllerTest(*testing.T)
}

type ControllerTests []ControllerTest

func (controllerTests ControllerTests) InternalTestSetup() (tests []testing.InternalTest) {
	for _, test := range controllerTests {
		name := reflect.TypeOf(test).Name()
		tests = append(tests, testing.InternalTest{
			Name: name,
			F:    test.SetupControllerTest,
		})
	}
	return tests
}

func NewControllerTests() ControllerTests {
	return ControllerTests{}
}

var ControllerIntegrationTestModules = fx.Options(
	config.TestENVModule,
	config.BaseModule,
	repository.Module,
	services.Module,
	controllers.Module,
	fx.Supply(config.EnvPath("../../.test.env")),
	fx.Provide(router.NewRouter),
	fx.Provide(NewControllerTests),
	fx.Invoke(bootstrapRepoTest),
)

func bootstrapRepoTest(
	repoTests ControllerTests,
) {
	os.Exit(testing.MainStart(
		tests.MatchStringOnly{},
		repoTests.InternalTestSetup(),
		nil,
		nil,
		nil,
	).Run())
}

func TestMain(m *testing.M) {
	fx.New(ControllerIntegrationTestModules).Run()
}
