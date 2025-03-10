# Serialization

Serialization is the process of converting Go data into a format that can be stored or transmitted. Deserialization is the process of converting serialized data back into its original Go form.

The classic example is transforming **Go data into JSON** and back.

Fuego automatically serializes and deserializes inputs and outputs with standard `encoding/json` package.

## Serialize data

To serialize data, just return the data you want to serialize from your controller. It will be automatically serialized into JSON, XML, YAML, or HTML, depending on the `Accept` header in the request.

- JSON `Accept: application/json` (default) (this default can be changed as an option in the `fuego.Server` struct)
- XML `Accept: application/xml`
- YAML `Accept: application/yaml`
- HTML `Accept: text/html`
- Plain text `Accept: text/plain`

```go
type MyReturnType struct {
	Message string `json:"message"`
}

func helloWorld(c fuego.ContextNoBody) (MyReturnType, error) {
	return MyReturnType{Message: "Hello, World!"}, nil
}

// curl request: curl -X GET http://localhost:8080/ -H "Accept: application/json"
// response: {"message":"Hello, World!"}

// curl request: curl -X GET http://localhost:8080/ -H "Accept: application/xml"
// response: <MyReturnType><Message>Hello, World!</Message></MyReturnType>
```

## Custom response - Bypass return type

If you want to bypass the automatic serialization, you can directly write to the response writer.

```go
func helloWorld(c fuego.ContextNoBody) (any, error) {
	w := c.Response()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello, World!"))                                      // Write directly to the response writer.
	_ = json.NewEncoder(w).Encode(MyReturnType{Message: "Hello, World!"}) // You can also use json.NewEncoder(w).Encode to serialize data directly into JSON
	fuego.SendJSON(w, MyReturnType{Message: "Hello, World!"})             // Or use fuego.SendJSON to serialize data directly into JSON

	return nil, nil // If you return nil, nil fuego will not try to serialize a response
}
```

## Custom serialization

But you can also use the `Serialize` and `Deserialize` functions to manually serialize and deserialize data.

See [the documentation](https://pkg.go.dev/github.com/go-fuego/fuego#Server) for the `fuego.Server` struct for more information.

```go
package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	s := fuego.NewServer()

	s.Serialize = func(w http.ResponseWriter, ans any) {
		w.Header().Set("Content-Type", "text/plain")
		json.NewEncoder(w).Encode(ans)
	}

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```
