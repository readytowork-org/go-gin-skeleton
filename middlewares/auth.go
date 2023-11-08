package middlewares

import "boilerplate-api/infrastructure"

// AuthMiddleware struct
type AuthMiddleware struct {
	logger infrastructure.Logger
}

func AuthMiddlewareConstrctor(
	logger infrastructure.Logger,
) AuthMiddleware {
	return AuthMiddleware{
		logger: logger,
	}
}
