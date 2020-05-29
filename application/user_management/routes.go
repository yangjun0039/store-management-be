package user_management

import (
	"net/http"
	"store-management-be/baselib/network"
	"store-management-be/application/user_management/member"
)

func MountSubrouterOn(r *network.Router) {
	r.LoadOnRoutes("/user", allRoute())
}

func allRoute() network.Routes {
	routes := network.Routes{
		{Method: http.MethodPost, Path: "/member-info", HandlerFunc: network.DefaultAdapter.Run(member.MemberInfo)},
		{Method: http.MethodPost, Path: "/add-member", HandlerFunc: network.DefaultAdapter.Run(member.AddMember)},
	}
	return routes
}
