package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	conn, err := connect("budgo", "devdev", "budgo", "5432")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := &controller{conn}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/in", withError(c.in))
	r.Post("/out", withError(c.out))
	r.Get("/*", withError(c.serve))
	r.Post("/*", withError(c.serve))
	fmt.Println("Starting server on :8081")
	log.Fatalln(http.ListenAndServe(":8081", r))
}

type controller struct {
	conn *conn
}

func withError(next func(w http.ResponseWriter, r *http.Request) (int, error)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		code, err := next(w, r)
		if err != nil {
			fmt.Println(err)
			http.Error(w, Err(err).JSON(), code)
		}
	}

	return http.HandlerFunc(fn)
}

// Data is the data for the homepage
type Data struct {
	Total   int
	Records []*Record
}

func (c *controller) serve(w http.ResponseWriter, r *http.Request) (int, error) {
	fm := template.FuncMap{
		"sum": func(x, y int) int {
			return x + y
		},
		"divide": func(a, b int) int {
			return a / b
		}}

	t, err := template.New("index.html").Funcs(fm).ParseFiles("./static/index.html")
	if err != nil {
		panic(err)
	}

	records, err := c.conn.records()
	if err != nil {
		return 400, err
	}
	total := 0
	for _, record := range records {
		total += record.Cents
	}

	data := &Data{
		Total:   total / 100,
		Records: records,
	}

	err = t.Execute(w, data)
	if err != nil {
		return 400, err
	}
	return 200, nil
}

func (c *controller) in(w http.ResponseWriter, r *http.Request) (int, error) {

	r.ParseForm()
	description := r.FormValue("description")
	if description == "" {
		return 400, errors.New("description not provided")
	}
	category := r.FormValue("category")
	if category == "" {
		return 400, errors.New("category not provided")
	}
	amount := r.FormValue("amount")
	if amount == "" {
		return 400, errors.New("amount not provided")
	}

	dollars, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 400, err
	}
	err = c.conn.saveRecord(description, category, int(dollars*100))
	if err != nil {
		return 400, err
	}
	http.Redirect(w, r, "/?success=true&type=in", http.StatusTemporaryRedirect)
	return 200, nil
}

func (c *controller) out(w http.ResponseWriter, r *http.Request) (int, error) {
	r.ParseForm()
	description := r.FormValue("description")
	if description == "" {
		return 400, errors.New("description not provided")
	}
	category := r.FormValue("category")
	if category == "" {
		return 400, errors.New("category not provided")
	}
	amount := r.FormValue("amount")
	if amount == "" {
		return 400, errors.New("amount not provided")
	}

	dollars, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 400, err
	}
	err = c.conn.saveRecord(description, category, int(dollars*-100))
	if err != nil {
		return 400, err
	}
	http.Redirect(w, r, "/?success=true&type=out", http.StatusTemporaryRedirect)
	return 200, nil
}
