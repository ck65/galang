package main

import (
	"html/template"
	"net/http"
)

type Book struct {
	Name   string
	Author string
	Photo  string
	Info   string
	Buy    string
}

type Books struct {
	Books []Book
}

func book(w http.ResponseWriter, r *http.Request) {
	bo := Books{}
	db := getDb()
	defer db.Close()
	res, _ := db.Query("select Name, Author, Photo, Info ,buy from books")
	for res.Next() {
		var name string
		var author string
		var photo string
		var info string
		var buy string
		res.Scan(&name, &author, &photo, &info, &buy)
		b := Book{name, author, photo, info, buy}
		bo.Books = append(bo.Books, b)
	}
	t, _ := template.ParseFiles("./static/html/book.html")
	t.Execute(w, bo)
}
