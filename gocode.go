package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func handleError(info string, err error) {
	if err != nil {
		fmt.Println("info: ", info)
		panic(err)
	}
}

func templatePageHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Page endpoint")
	time := time.Now().Format(time.RFC1123)

	tmpl := template.Must(template.ParseFiles("public/template.html"))

	type TmplPageData struct {
		Links template.HTML
		Body  template.HTML
		Time  string
		Code  string
	}

	err := r.ParseForm()
	handleError("parse form error:", err)
	code := r.Form.Get("code")

	log.Println("Page template ", r.URL.Path)
	err := r.ParseForm()
	handleError("Read template file", err)

	file := r.URL.Path

	if file == "/" {
		file = "/index"
	}

	if _, err := os.Stat("public" + file + ".html"); !os.IsNotExist(err) {

		bodyString, rferr := ioutil.ReadFile("public" + file + ".html")
		handleError("read page file", rferr)

		headerLinksString, rferr2 := ioutil.ReadFile("public/header_links.html")
		handleError("read header file", rferr2)

		headerLinks := template.HTML(string(headerLinksString))
		body := template.HTML(string(bodyString))

		data := TmplPageData{
			Links: headerLinks,
			Body:  body,
			Time:  time,
			Code:  code,
		}

		tmpl_err := tmpl.Execute(w, data)
		handleError("template execute error:", tmpl_err)
	} else {
		fmt.Fprintf(w, "404")
	}
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	root := http.HandlerFunc(templatePageHandler)
	http.HandleFunc("/", root)

	fmt.Println("Server: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
