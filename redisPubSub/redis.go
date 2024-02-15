package redisPubSub

import (
	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewMessageConsumer),
	fx.Provide(NewMessagePublisher),
)
