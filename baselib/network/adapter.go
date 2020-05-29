package network

var DefaultAdapter = NewHandlerBuilder().Adapt(RequesterHandler).Then(RecoveryHandler, JWTHandler, ForbidRequestTooOftenHandler)

var LoginHandler = NewHandlerBuilder().Adapt(RequesterHandler).Then(RecoveryHandler, NoJWTHandler, WebLoginTokenHandler, ForbidRequestTooOftenHandler)

