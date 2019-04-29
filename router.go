package router

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Router contains app defined list of routes
type Router struct {
	routes []*Route
	group  RouteGroup // if this router is group - this field contains prefix and middlewares
}

// Route entity
type Route struct {
	Method    string
	URI       string
	Handlers  []RouteHandler
	pattern   string
	variables []string
}

// RouteGroup type
type RouteGroup struct {
	Prefix     string
	Middleware []RouteHandler
}

// RouteHandler function type
type RouteHandler func(*Context) *Context

type httpHandler struct {
	router *Router
}

// Handler returns http handler instance
func (router *Router) Handler() http.Handler {
	return &httpHandler{
		router: router,
	}
}

func (h httpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// build context
	ctx := createContext(response, request)
	log.Printf("%s %s\n", ctx.Method, ctx.URI)

	// find handler
	handlers, err := getHandlerForRoute(ctx, h.router.routes)

	// check for route error - not found, method not allowed
	if err != nil {
		ctx.Error(err)
		return
	}

	// execute route handlers
	executeHandlers(handlers, ctx)
}

func getHandlerForRoute(ctx *Context, routes []*Route) ([]RouteHandler, *HTTPError) {
	for _, route := range routes {
		// check if URL matches some of predefined endpoints
		if len(matchURL(route, ctx.URI)) > 0 {
			// if we have a URL match, check if there is specified endpoint with given http method
			if match := getMatchForMethod(route, ctx.URI, ctx.Method, routes); len(match) > 0 {
				// bind URL variable placeholders with parsed values
				// and bind those params to ctx
				values := match[1:]
				ctx.Params = mapURLParams(route.variables, values)

				// return router handlers array
				return (*route).Handlers, nil
			} else {
				// If there is route match but with different method we mark it as HTTP 405
				return nil, MethodNotAllowedError()
			}
		}
	}
	return nil, NotFoundError()
}

func matchURL(route *Route, path string) []string {
	routeRegex := regexp.MustCompile(route.pattern)
	// if there is match it may contain variables from URL placeholders
	return routeRegex.FindStringSubmatch(path)
}

func getMatchForMethod(route *Route, path string, method string, routes []*Route) []string {
	for _, route := range routes {
		match := matchURL(route, path)
		if len(match) > 0 && route.Method == method {
			return match
		}
	}
	return []string{}
}

func executeHandlers(handlers []RouteHandler, ctx *Context) *Context {
	// If response has been already set, don't continue executing other middlewares
	if ctx.responseSent {
		return ctx
	}

	// Execute current middleware in pipeline
	handlers[0](ctx)

	if len(handlers) == 1 {
		return ctx
	}

	// shift current middleware and proceed to next in array
	return executeHandlers(handlers[1:], ctx)
}

// Initialize returns router
func Initialize() *Router {
	router := new(Router)
	router.routes = *new([]*Route)
	return router
}

// Group initializes new router with common properties like prefix or middlewares
func Group(group *RouteGroup) *Router {
	groupRouter := Initialize()
	groupRouter.group = *group
	return groupRouter
}

// Use takes grouped routes and adds it to root router
func (router *Router) Use(groupRouter *Router) {
	for _, route := range groupRouter.routes {
		// Add route from group
		router.addRoute(
			route.Method, // method remains the same
			groupRouter.group.Prefix+"/"+strings.TrimPrefix(route.URI, "/"), // e.g. /prefix/original-url
			append(groupRouter.group.Middleware, route.Handlers...),         // [group handlers, route handlers]
		)
	}
}

func (router *Router) addRoute(method string, path string, handlers []RouteHandler) {
	// parse URI to pattern
	// replace URI variable placeholders with regex
	// e.g. /users/{userID}/posts/{postID} -> ^/users/([^\/]+?)/posts/([^\/]+?)$
	var re = regexp.MustCompile(`\{(.*?)\}`)                      // e.g. {userId}
	pattern := "^" + re.ReplaceAllString(path, `([^\/]+?)`) + "$" // will be replaced with ([^\/]+?)

	// TODO: add support for optional parameters e.g. /users/{id?}
	// TODO: add support for specifying additional regex rules for variable e.g. /users/{id:[a-zA-Z]+}, or /users/{id:[0-9]+}

	// extract URI variables
	// we pick all these variable placeholders and save it for later data binding
	variablesMatch := re.FindAllStringSubmatch(path, -1)
	variables := []string{}
	for _, varMatch := range variablesMatch {
		variables = append(variables, varMatch[1])
	}

	// append route
	(*router).routes = append((*router).routes, &Route{
		Method:    method,
		URI:       normalizeURLSlashes(path),
		Handlers:  handlers,
		pattern:   pattern,
		variables: variables,
	})
}

// Get defines HTTP GET route
func (router *Router) Get(path string, handlers ...RouteHandler) {
	router.addRoute(http.MethodGet, path, handlers)
}

// Post defines HTTP POST route
func (router *Router) Post(path string, handlers ...RouteHandler) {
	router.addRoute(http.MethodPost, path, handlers)
}

// Put defines HTTP PUT route
func (router *Router) Put(path string, handlers ...RouteHandler) {
	router.addRoute(http.MethodPut, path, handlers)
}

// Delete defines HTTP DELETE route
func (router *Router) Delete(path string, handlers ...RouteHandler) {
	router.addRoute(http.MethodDelete, path, handlers)
}

func mapURLParams(names, values []string) URLParams {
	params := make(URLParams)
	for i, name := range names {
		params[name] = values[i]
	}
	return params
}

func normalizeURLSlashes(url string) string {
	url = strings.TrimSuffix(url, "/")       // trim trailing /
	url = "/" + strings.TrimPrefix(url, "/") // make sure every route starts with /
	return url
}
