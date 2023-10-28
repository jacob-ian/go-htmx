package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"jacobmatthe.ws/htmx"
	"jacobmatthe.ws/htmx/internal/template"
)

type index struct {
	Name string
}

var tmpl *template.TemplateEngine

type item struct {
	Text string
	Id   int
}

func main() {
	log.Println("Starting up...")

	var err error
	tmpl, err = template.NewTemplateEngine()
	if err != nil {
		log.Fatal(err)
	}

	var items []item
	items = make([]item, 0)

	router := httprouter.New()
	router.HandlerFunc("GET", "/", logRequest(NewHomeHandler(&items)))
	router.HandlerFunc("GET", "/about", logRequest(NewAboutHandler()))
	router.HandlerFunc("POST", "/items", logRequest(NewAddItemHandler(&items)))
	router.HandlerFunc("DELETE", "/items/:id", logRequest(NewDeleteItemHandler(&items)))
	router.HandlerFunc("GET", "/assets/:path", logRequest(htmx.NewStaticFileServer()))

	log.Println("Listening on localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}

func logRequest(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		next.ServeHTTP(w, r)
	})
}

func NewHomeHandler(items *[]item) http.Handler {
	type home struct {
		Title string
		Name  string
		Items []item
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(items)
		err := tmpl.Execute(w, "home.html", home{Title: "Home", Name: "Jacob", Items: *items})
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}
	})
}

func NewAboutHandler() http.Handler {
	type post struct {
		Title       string
		Description string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, "about.html", nil)
		if err != nil {
			log.Println(err)
		}
	})
}

func NewAddItemHandler(items *[]item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := item{
			Text: r.FormValue("text"),
			Id:   len(*items),
		}
		*items = append(*items, data)
		err := tmpl.Execute(w, "list-item.html", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func NewDeleteItemHandler(items *[]item) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		out := make([]item, len(*items)-1)
		for i, item := range *items {
			if i == id {
				continue
			}
			out = append(out, item)
		}
		*items = out
		w.WriteHeader(http.StatusOK)
	}
}
