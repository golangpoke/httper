package httper

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"runtime"
)

type Middleware func(next http.Handler) http.Handler

type serveMux struct {
	mux *http.ServeMux
	// root handle
	IsBreakRootRouter bool
	NotFoundHandler   http.Handler
	// group
	routers     map[string]string
	groupPrefix string
	isRootGroup bool
	middlewares []Middleware
}

func NewServeMux() *serveMux {
	return &serveMux{
		mux: http.NewServeMux(),
		// root handle
		IsBreakRootRouter: true,
		NotFoundHandler:   http.NotFoundHandler(),
		// group
		routers:     make(map[string]string),
		isRootGroup: true,
		middlewares: make([]Middleware, 0),
	}
}

func (sm *serveMux) Use(middleware ...Middleware) {
	sm.middlewares = append(sm.middlewares, middleware...)
}

func (sm *serveMux) Group(prefix string) *serveMux {
	middlewares := make([]Middleware, 0)
	// add middlewares of parent serve mux,but not root serve mux
	if !sm.isRootGroup {
		middlewares = append(middlewares, sm.middlewares...)
	}
	newServeMux := &serveMux{
		mux:               sm.mux,
		IsBreakRootRouter: sm.IsBreakRootRouter,
		NotFoundHandler:   sm.NotFoundHandler,
		routers:           sm.routers,
		groupPrefix:       sm.groupPrefix + prefix,
		isRootGroup:       false,
		middlewares:       middlewares,
	}
	return newServeMux
}

func (sm *serveMux) GET(router string, handler http.Handler) {
	sm.registerRouter(http.MethodGet, router, handler)
}

func (sm *serveMux) POST(router string, handler http.Handler) {
	sm.registerRouter(http.MethodPost, router, handler)
}

func (sm *serveMux) PUT(router string, handler http.Handler) {
	sm.registerRouter(http.MethodPut, router, handler)
}

func (sm *serveMux) DELETE(router string, handler http.Handler) {
	sm.registerRouter(http.MethodDelete, router, handler)
}

func (sm *serveMux) registerRouter(method string, router string, handler http.Handler) {
	// use middlewares of not root serve mux
	if !sm.isRootGroup {
		handler = sm.useMiddlewares(handler, sm.middlewares...)
	}
	router = fmt.Sprintf("%s %s%s", method, sm.groupPrefix, router)
	sm.routers[router] = funcName(handler)
	sm.mux.Handle(router, handler)
}

func (sm *serveMux) Routes() map[string]string {
	return sm.routers
}

func (sm *serveMux) Start(addr string) error {
	var handler http.Handler
	handler = sm.mux
	if sm.IsBreakRootRouter {
		handler = sm.breakRootRouterHandler(sm.mux)
	}
	// use middlewares of the root serve mux
	handler = sm.useMiddlewares(handler, sm.middlewares...)
	return http.ListenAndServe(addr, handler)
}

func (sm *serveMux) breakRootRouterHandler(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		rootRouter := fmt.Sprintf("%s %s", r.Method, "/")
		// break root pattern
		if pattern == rootRouter && r.URL.Path != "/" {
			sm.NotFoundHandler.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (sm *serveMux) useMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	l := len(sm.middlewares) - 1
	// use middlewares reversely
	for i := l; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func funcName(fn any) string {
	n := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return path.Base(n)
}
