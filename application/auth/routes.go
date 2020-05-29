package auth

import (
	"net/http"
	"store-management-be/baselib/network"
	"store-management-be/application/auth/validation"
)

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/authentication", allRoute())
}

func allRoute() network.Routes {
	routes := network.Routes{
		{http.MethodPost, "/login", network.LoginHandler.Run(validation.Login)},
	}
	return routes
}


//var loginRputer = network.Routes{
//	{http.MethodPost, "/login", network.LoginHandler.Run(validation.Login)},
//}

