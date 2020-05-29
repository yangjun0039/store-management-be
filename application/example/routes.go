package example

import (
	"store-management-be/baselib/network"
	"net/http"
	"store-management-be/application/example/example1"
)

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/example", allRoute())
}

func allRoute() network.Routes {
	routes := network.Routes{
		{Method: http.MethodPost, Path: "/test1", HandlerFunc: network.DefaultAdapter.Run(example1.GetMsg)},
	}
	return routes
}
