package example

import (
	"store-management-be/baselib/network"
	"net/http"
	"store-management-be/application/example/example1"
)

func RouterRegisterMethods(r *network.Router) {
	r.LoadOnRoutes("/example", allRoute())
}

func allRoute() network.Routes {
	routes := network.Routes{
		{Method: http.MethodPost, Path: "/test1", HandlerFunc: example1.GetMsg},
	}
	return routes
}
