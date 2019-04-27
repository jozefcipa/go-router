# go-router
> Simple HTTP router wrapper for Go net/http package

## Install
`go get -u github.com/jozefcipa/go-router`

## Example
```go
func main() {
    // Initialize router
    r := router.Initialize()

    // Define routes
    r.Get("/", func(ctx *router.Context) *router.Context {
        ctx.Send("Hello world!")
        return ctx
    })

    // Start server at localhost:8000
    fmt.Println("Starting server at localhost:8000")
    if err := http.ListenAndServe(":8000", r.Handler()); err != nil {
        panic(err)
    }
}
```
## Documentation

### <a name="routing">Routing</a>
List of available functions that you can use to define HTTP endpoints

  - `func (r Router) Get(path string, handlers ...RouteHandler) *Context`
  - `func (r Router) Post(path string, handlers ...RouteHandler) *Context`
  - `func (r Router) Put(path string, handlers ...RouteHandler) *Context`
  - `func (r Router) Delete(path string, handlers ...RouteHandler) *Context`

#### <a name="url-params">URL Parameters</a>
You can define route which contains wildcard
```go
r.Get("/users/{userID}", func(ctx *router.Context) *router.Context {
    // you can access URL parameters via `ctx.Params` map
    fmt.Println("user id is", ctx.Params["userID"])
    ctx.Send("yay")
    return ctx
})
```

### <a name="middlewares">Middlewares</a>
Every endpoint function accepts one or multiple functions that will be called when endpoint is accessed. You can pass them as multiple subsequent parameters in route function or using `Middleware` array in `Group` function. They are all of type `RouteHandler` (<a href="router.go#L32">#see</a>).

### <a name="groups">Groups</a>
You can also use `Group` function to define common **prefix** or set of **middlewares**
```go
// Define routes group 
// Group function returns router instance
v1Router := router.Group(&router.RouteGroup{
  Prefix: "/v1",
  Middleware: []router.RouteHandler{
    func(ctx *router.Context) *router.Context {
      fmt.Println("v1 middleware")
      return ctx
    },
  },
})

// Define /foo route on `v1Router`
// this route will be resolved as /v1/foo 
// and will have all middlewares registered for group
v1Router.Get("/foo", func(ctx *router.Context) *router.Context {
  ctx.Send("Hello, this is /v1/foo")
  return ctx
})

// Add created group router to root router
r.Use(v1Router)
```

### <a name="response">Response</a>
This uses JSON responses as default. You can use `ctx.Send()` function to send response.

### <a name="errors">Errors</a>
If error ocurres, you can send one of predefined HTTP errors via `ctx.Error(httpError HTTPError)`.
#### HTTPErrors
- `func BadRequestError(payload ErrorPayload, msg ...string) *HTTPError // 400`
- `func UnauthorizedError(payload ErrorPayload, msg ...string) *HTTPError // 401`
- `func NotFoundError(msg ...string) *HTTPError // 404`
- `func func MethodNotAllowedError() *HTTPError // 405 - This is used by router itself primarily`
- `func InternalError(msg ...string) *HTTPError // 500`

All these functions return JSON response.
**Note**: `payload` field is **only** included if it was provided in error function
```json
{
    "error": "ERROR MESSAGE",
    "statusCode": 123,
    "payload": {}
}
```

### <a name="errors">Context</a>
Context is the only variable which is passed to every `RouterHandler` function and has following structure 
```go
type Context struct {
	Res          http.ResponseWriter
	Req          *http.Request
	Params       URLParams
	Query        interface{} // query parameters
	Method       string      // HTTP method
	URI          string
}
```

You can also store globally accessible data in context (e.g. Store authenticated user in middleware and access it later in route handler)
```go
ctx.Set("key", "value")
value, ok := ctx.Get("key") // ok returns false if key doesn't exist
```
