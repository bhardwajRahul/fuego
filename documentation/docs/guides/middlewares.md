# Middlewares

Fuego supports `net/http` middlewares.

It means that all existing middlewares for `net/http`,
including the ones from `chi` and `gorilla` can be used with Fuego! :fire:

You can use them to add functionalities to your routes, such as logging,
authentication, etc.

## App-level middlewares

You can add middlewares to the whole server using the `Use` method:

```go title="main.go" showLineNumbers
package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	// Add a middleware to the whole server
	fuego.Use(s, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do something before the request
			next.ServeHTTP(w, r)
			// Do something after the request
		})
	})

	fuego.Use(s, myMiddleware)

	fuego.Get(s, "/", myController)

	s.Run()
}
```

## Group middlewares

You can also add middlewares to a group of routes using the `Group` method:

```go title="main.go" showLineNumbers
package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	// Add a middleware to a group of routes
	api := fuego.Group(s, "/api")

	fuego.Use(api, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do something before the request
			// Only affects routes in the group
			next.ServeHTTP(w, r)
			// Do something after the request
		})
	})

	// Requests to /api will go through the middleware
	fuego.Get(api, "/", myController)

	// Requests to / will NOT! go through the middleware
	fuego.Get(s, "/", myController)

	s.Run()
}
```

## Route middlewares

You can also add middlewares to a single route.
They are treated as an option to the route handler.
Simply add the middlewares as the last arguments of the route handler:

```go title="main.go" showLineNumbers {13-14}
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	// Declare the middlewares after the route handler
	fuego.Get(s, "/", myController,
		option.QueryInt("page", "The page number"),
		option.Middleware(middleware1),
		option.Middleware(middleware2, middleware3),
	)

	s.Run()
}
```
