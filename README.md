#Overview
This project is a yet another micro framework developed in Go. What
stands out in this case, is it's probably the only web framework for
Go that is both - lightweight and ready for writing production code,
right out of the box. It provides zero-configuration & hassle free
essentials that are absolutely necessary for any production app, but
which are missing from most of the web frameworks for Go. The goal
of the framework is to let you write your business logic since 1-st
line of code, without having to write all the usual boilerplate code
that is necessary to power a production application and its deployment.
Here is a list of some of the features supported out of the box.
  
* Switching app environments: dev, production, staging, etc
* Automatically load a config file based on the environment
* Declarative API for endpoint definition, with metadata
* Serving inventory of available APIs defined in the app
* Optional built-in UI for browsing and testing app APIs
* Structured logger is automatically available in the app
* Middleware support for the endpoints and endpoint groups
* Toolkit with essentials for authentication & access control
* High performance and efficient memory consumption
* Efficiently serving static resources (asset files)

#Code Samples

##Making the Imports
```go
import "flag"
import "net/http"
import "encoding/json"
import "gopkg.in/mgo.v2/bson"

import "github.com/ts33kr/boot"
import "github.com/ts33kr/boot/acl"
import "github.com/ts33kr/boot/auth"
import "github.com/ts33kr/boot/middle"
```

##Defining a GET Endpoint
```go
func ListBooks (ep *boot.Endpoint) {
    ep.Get("/filter/on-the-shelves/:id?")
    ep.Use(auth.HasAuthenticatedAccount)
    ep.Use(acl.Roles("personal", "admin"))
    ep.Use(acl.Grants("list-books", "get-book"))
    ep.Opt("id", "optional ID of the onewidget", nil)
    ep.Opt("skip", "offset integer when paginating", 0)
    ep.Opt("limit", "number of items to limit query", 0)
    ep.Available("production", "staging", "development")
    ep.Implement(func (c boot.Context, app *boot.App) {
        var responder http.ResponseWriter = c.responder
        account := auth.Account(c.Get("account", nil))
        id := bson.ObjectIdHex(c.Param("id", nil))
    })
}
```

##Defining a POST Endpoint
```go
func AddBook (ep *boot.Endpoint) {
    ep.Post("/filter/on-the-shelves")
    ep.Use(auth.HasAuthenticatedAccount)
    ep.Use(acl.Roles("personal", "admin"))
    ep.Use(acl.Grants("list-books", "add-book"))
    ep.Arg("title", "title of the book being added")
    ep.Arg("author", "name of the auththor who wrote it")
    ep.Arg("published", "year that the book was published")
    ep.Available("production", "staging", "development")
    ep.Implement(func (c boot.Context, app *boot.App) {
        var responder http.ResponseWriter = c.responder
        account := auth.Account(c.Get("account", nil))
    })
}
```

##Assemble and Run the App
```go
func main () {
    var app boot.App = boot.Make("book-store-app")
    env := flag.String("e", "development", "env to use")
    ver := flag.String("v", "0.0.1", "version of an app")
    defer flag.Parse() // parse command line arguments
    defer app.Bootload(*env, *ver, "./", "config/")
    books := app.Group("/books", "Manage Books")
    books.Use(auth.Strategy(auth.HttpBasic()))
    books.Use(middle.MustUseSSLEncryption)
    books.Use(middle.JournalRequests)
    books.Mount(ListBooks)
    books.Mount(AddBook)
}
```
