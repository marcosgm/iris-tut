package main

import (
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
)

// Book example.
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

var books []Book
var bookstore_title string

func init() {
	books = []Book{
		{"IT", "Stephen King", "horror"},
		{"Game of Thrones", "George RR Martin", "fantasy"},
		{"The lord of the rings", "JRR Tolkien", "fantasy"},
	}
	bookstore_title = os.Getenv("bookstore_title")
	bookstore_title, ok := os.LookupEnv("bookstore_title")
	if !ok {
		fmt.Println("bookstore_title env var is not present, setting to Default")
		bookstore_title = "Generic Bookstore"
	} else {
		fmt.Printf("Bookstore title is: %s\n", bookstore_title)
	}
}

func main() {
	app := iris.New()
	app.Favicon("./assets/favicon.ico")
	app.HandleDir("/assets", iris.Dir("./assets"))

	tmpl := iris.Django("./views", ".html")
	tmpl.Reload(true) //only for dev mode
	app.RegisterView(tmpl)

	// app.Get("/", func(ctx iris.Context) {
	// 	ctx.HTML("<h2> Hello World</h2>")
	// })
	app.Get("/", index)

	booksAPI := app.Party("/books")
	{
		booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
		// POST: http://localhost:8080/books
		booksAPI.Post("/", create)
	}

	app.Listen(":8080")
}
func index(ctx iris.Context) {
	// ctx.ViewData("title", "Hi Page")
	// ctx.ViewData("name", "iris")
	// ctx.ViewData("serverStartTime", startTime)
	// or if you set all view data in the same handler you can use the
	// iris.Map/pongo2.Context/map[string]interface{}, look below:

	if err := ctx.View("index.html", iris.Map{
		"title":           "Welcome Page",
		"name":            "iris Tut",
		"books":           books,
		"bookstore_title": bookstore_title,
	}); err != nil {
		ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func list(ctx iris.Context) {

	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func create(ctx iris.Context) {
	var b Book
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Book creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Book: " + b.Title + " by " + b.Author + "type " + b.Genre)
	books = append(books, b)
	ctx.StatusCode(iris.StatusCreated)
}

// curl -i -X POST --header 'Content-Type:application/json' --data "{\"title\":\"New Book\",\"author\":\"New Author\",\"genre\":\"Fiction\"}" http://localhost:8080/books
