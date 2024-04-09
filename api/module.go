package api

import (
	"boilerplate-api/api/auth"
	"boilerplate-api/api/docs"
	"boilerplate-api/api/gcp_billing"
	"boilerplate-api/api/user"
	"boilerplate-api/api/utility"
	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Options(
		docs.Module,
		auth.Module,
		user.Module,
		gcp_billing.Module,
		utility.Module,
	),
)
