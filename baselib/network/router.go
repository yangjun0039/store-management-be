package network

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"store-management-be/configer"
	"encoding/json"
	"time"
	"store-management-be/baselib/logger"
)

type NetProtocol int

const (
	HTTP  NetProtocol = iota
	HTTPS
)

var rootRouter *Router

type Router struct {
	*mux.Router
}

type Routes []Route

type Route struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

func NewRouter(routes Routes) *Router {
	aRouter := &Router{mux.NewRouter().StrictSlash(true)}
	return aRouter
}

func (r Route) Name() string {
	return fmt.Sprintf("[%v]%v", r.Method, r.Path)
}

func (r *Router) LoadOnRoutes(pathPrefix string, routes Routes) {
	aSubrouter := r.PathPrefix(pathPrefix).Subrouter()
	for _, route := range routes {
		aSubrouter.Methods(route.Method).Path(route.Path).Name(route.Name()).Handler(route.HandlerFunc)
	}
}

func (r *Router) timeoutHandler() http.Handler {
	const timeoutDuration = 60
	failure := NewHttpFailure(http.StatusInternalServerError, responseTimeout)
	bytes, _ := json.Marshal(*failure)
	return http.TimeoutHandler(r, timeoutDuration*time.Second, string(bytes))
}

type Handlers []FileHandler
type FileHandler struct {
	Path    string
	Handler http.Handler
}
// 静态文件路由
func (r *Router) FileLoadOnRoutes(hds Handlers) {
	for _, h := range hds {
		aSubrouter := r.PathPrefix(h.Path).Subrouter()
		aSubrouter.Methods(http.MethodGet).Handler(h.Handler)
	}
}

func (r *Router) Startup(protocol NetProtocol, port uint64) {
	rootRouter = r
	server := http.Server{Addr: fmt.Sprintf(":%v", port), Handler: r.timeoutHandler()}
	if protocol == HTTPS {
		certFile := configer.Conf.Key.Certificate
		keyFile := configer.Conf.Key.Private
		logger.LogSugar.Infof("https server start up, port:%d", port)
		server.ListenAndServeTLS(certFile, keyFile)
	} else if protocol == HTTP {
		logger.LogSugar.Infof("http server start up, port:%d", port)
		server.ListenAndServe()
	}
}