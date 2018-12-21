package main

import (
	"html/template"
	"net/http"
)

type tools struct {
	Name     string
	Info     string
	Size     string
	Download string
	Update   string
	Type     string
}

type Tools struct {
	Tool []tools
}

func tool(w http.ResponseWriter, r *http.Request) {
	t := Tools{}
	db := getDb()
	defer db.Close()
	res, err := db.Query("select Name, Info, Size, Download, Updata, tltype from tools")
	checkerr(err)
	for res.Next() {
		var name string
		var info string
		var size string
		var download string
		var update string
		var tltype string
		res.Scan(&name, &info, &size, &download, &update, &tltype)
		to := tools{name, info, size, download, update, tltype}
		t.Tool = append(t.Tool, to)
	}
	q, _ := template.ParseFiles("./static/html/tools.html")
	q.Execute(w, t)
}
