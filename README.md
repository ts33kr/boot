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
