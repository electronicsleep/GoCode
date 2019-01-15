package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleError(info string, err error) {
	if err != nil {
		fmt.Println("info: ", info)
		panic(err)
	}
}

func templatePageHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Page endpoint: ", r.URL.Path)
	rfc_time := time.Now().Format(time.RFC1123)

	t := time.Now()
	tm := t.Format("20060102")

	tmpl := template.Must(template.ParseFiles("public/template.html"))

	type TmplPageData struct {
		Links template.HTML
		Body  template.HTML
		Save  string
		Code  string
		Time  string
	}

	err := r.ParseForm()
	handleError("parse form error:", err)
	code := r.Form.Get("code")
	log.Println("-----> Code")

	filename := randString(25)
	save := "public/" + filename + "-" + tm + ".txt"
	filename = "public/" + filename + "-" + tm + ".txt"
	log.Println("save: " + save)

	saveBytes := []byte(code)
	saveErr := ioutil.WriteFile(filename, saveBytes, 0644)
	handleError("file save error: ", saveErr)

	var lines []string = strings.Split(code, "\n")
	for index, line := range lines {
		i := strconv.Itoa(index + 1)
		log.Println("Line "+i+": ", line)
	}

	log.Println("-----> End")

	ParseErr := r.ParseForm()
	handleError("Read template file", ParseErr)

	path := r.URL.Path

	if path == "/" {
		path = "/index"
	}
	file := "public" + path + ".html"

	if _, err := os.Stat(file); !os.IsNotExist(err) {

		bodyString, rferr := ioutil.ReadFile(file)
		handleError("read page file", rferr)

		headerLinksString, rferr2 := ioutil.ReadFile("public/header_links.html")
		handleError("read header file", rferr2)

		headerLinks := template.HTML(string(headerLinksString))
		body := template.HTML(string(bodyString))

		data := TmplPageData{
			Links: headerLinks,
			Body:  body,
			Save:  save,
			Code:  code,
			Time:  rfc_time,
		}

		tmpl_err := tmpl.Execute(w, data)
		handleError("template execute error:", tmpl_err)
	} else {
		fmt.Fprintf(w, "404")
	}
}

func randString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(randInt(65, 90)) != temp {
			temp = string(randInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	root := http.HandlerFunc(templatePageHandler)
	http.HandleFunc("/", root)

	fmt.Println("Server: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
