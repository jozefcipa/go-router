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
IN PROGRESS
