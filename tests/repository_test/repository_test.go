package repository_test

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/internal/config"
	"boilerplate-api/tests"
	"go.uber.org/fx"
	"os"
	"reflect"
	"testing"
)

type RepoTests []RepoTest
type RepoTest interface {
	SetupRepoTest(*testing.T)
}

func (repoTest RepoTests) InternalTestSetup() (tests []testing.InternalTest) {
	for _, test := range repoTest {
		name := reflect.TypeOf(test).Name()
		tests = append(tests, testing.InternalTest{
			Name: name,
			F:    test.SetupRepoTest,
		})
	}
	return tests
}

func NewRepoTests(
	planRepoTest TestPlanRepository,
	storeRepoTest TestStoreRepository,
	aminRepoTest TestAdminRepository,
	userRepoTest TestUserRepository,
) RepoTests {
	return RepoTests{
		planRepoTest,
		storeRepoTest,
		aminRepoTest,
		userRepoTest,
	}
}

var RepoTestModules = fx.Options(
	config.TestENVModule,
	config.BaseModule,
	repository.Module,
	fx.Supply(config.EnvPath("../../.test.env")),
	fx.Provide(NewTestPlanRepository),
	fx.Provide(NewTestStoreRepository),
	fx.Provide(NewTestAdminRepository),
	fx.Provide(NewTestUserRepository),
	fx.Provide(NewRepoTests),
	fx.Invoke(bootstrapRepoTest),
)

func bootstrapRepoTest(
	repoTests RepoTests,
	migrations config.Migrations,
) {
	migrations.MigrateUp()

	os.Exit(testing.MainStart(
		tests.MatchStringOnly{},
		repoTests.InternalTestSetup(),
		nil,
		nil,
		nil,
	).Run())
}

func TestMain(m *testing.M) {
	fx.New(RepoTestModules).Run()
}
