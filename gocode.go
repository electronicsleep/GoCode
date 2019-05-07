package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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

func templateHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Page endpoint: ", r.URL.Path)
	rfc_time := time.Now().Format(time.RFC1123)

	t := time.Now()
	tm := t.Format("20060102")
	savefile := ""

	tmpl := template.Must(template.ParseFiles("public/index.html"))

	type TmplPageData struct {
		Header template.HTML
		Links  template.HTML
		Body   template.HTML
		Footer template.HTML
		Save   string
		Get    string
		Code   string
		Time   string
	}

	err := r.ParseForm()
	handleError("parse form error:", err)
	code := r.Form.Get("code")
	public := r.Form.Get("public")
	getfile := r.Form.Get("getfile")
	log.Println("-----> Found getfile1")
	log.Println(getfile)
	if getfile != "" {
		log.Println("-----> Found getfile")
		log.Println(getfile)
		file, err := ioutil.ReadFile("public/code/" + getfile)
		handleError("error getting getfile", err)
		code = string(file)
		getfile = "?getfile=" + getfile

	} else if code == "" {
		log.Println("-----> Start")
		savefile = "NA"

	} else {
		log.Println("-----> Code")

		filename := randString(25)
		if public == "public" {
			public = "-" + public
		}
		//static savefile
		savefile = "public/code/" + filename + "-" + tm + public + ".txt"
		getfile = "?getfile=" + filename + "-" + tm + public + ".txt"
		log.Println("save: " + savefile)

		saveBytes := []byte(code)
		saveErr := ioutil.WriteFile(savefile, saveBytes, 0644)
		handleError("file save error: ", saveErr)

		var lines []string = strings.Split(code, "\n")
		for index, line := range lines {
			i := strconv.Itoa(index + 1)
			log.Println("Line "+i+": ", line)
		}

		log.Println("-----> End")
	}

	ParseErr := r.ParseForm()
	handleError("Read template file", ParseErr)

	path := r.URL.Path

	if path == "/" {
		path = "/index"
	}
	file := "public" + path + ".html"

	bodyString, rferr := ioutil.ReadFile(file)
	handleError("read page file", rferr)

	headerString, rferr2 := ioutil.ReadFile("public/header.html")
	handleError("read header file", rferr2)

	headerLinksString, rferr3 := ioutil.ReadFile("public/header_links.html")
	handleError("read header links file", rferr3)

	footerString, rferr3 := ioutil.ReadFile("public/footer.html")
	handleError("read footer file", rferr3)

	header := template.HTML(string(headerString))
	headerLinks := template.HTML(string(headerLinksString))
	body := template.HTML(string(bodyString))
	footer := template.HTML(string(footerString))

	data := TmplPageData{
		Header: header,
		Links:  headerLinks,
		Body:   body,
		Footer: footer,
		Save:   savefile,
		Get:    getfile,
		Code:   code,
		Time:   rfc_time,
	}

	tmpl_err := tmpl.Execute(w, data)
	handleError("template execute error:", tmpl_err)
}

func templateHandlerAbout(w http.ResponseWriter, r *http.Request) {
	log.Println("About endpoint: ", r.URL.Path)

	tmpl := template.Must(template.ParseFiles("public/about.html"))
	type TmplPageData struct {
		Header template.HTML
		Links  template.HTML
		Footer template.HTML
	}

	headerString, rferr2 := ioutil.ReadFile("public/header.html")
	handleError("read header file", rferr2)

	headerLinksString, rferr3 := ioutil.ReadFile("public/header_links.html")
	handleError("read header links file", rferr3)

	footerString, rferr3 := ioutil.ReadFile("public/footer.html")
	handleError("read footer file", rferr3)

	header := template.HTML(string(headerString))
	headerLinks := template.HTML(string(headerLinksString))
	footer := template.HTML(string(footerString))

	data := TmplPageData{
		Header: header,
		Links:  headerLinks,
		Footer: footer,
	}

	tmpl_err := tmpl.Execute(w, data)
	handleError("template execute error:", tmpl_err)
}

func templateHandlerHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("---> History endpoint: ", r.URL.Path)

	files, readErr := ioutil.ReadDir("public/code/")
	if readErr != nil {
		log.Fatal(readErr)
	}

	codeFiles := ""
	for _, f := range files {
		if f.Name() != "README.md" {
			if strings.Contains(f.Name(), "-public") {

				log.Println(f.Name())
				codeFiles += "<p><a href=public/code/"
				codeFiles += f.Name()
				codeFiles += " target=_blank>"
				codeFiles += f.Name()
				codeFiles += "</a>"
				codeFiles += "\n"
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("public/history.html"))
	type TmplPageData struct {
		Header template.HTML
		Links  template.HTML
		Footer template.HTML
		Files  template.HTML
	}

	headerString, rferr2 := ioutil.ReadFile("public/header.html")
	handleError("read header file", rferr2)

	headerLinksString, rferr3 := ioutil.ReadFile("public/header_links.html")
	handleError("read header links file", rferr3)

	footerString, rferr3 := ioutil.ReadFile("public/footer.html")
	handleError("read footer file", rferr3)

	header := template.HTML(string(headerString))
	headerLinks := template.HTML(string(headerLinksString))
	footer := template.HTML(string(footerString))

	Code := template.HTML(string(codeFiles))

	data := TmplPageData{
		Header: header,
		Links:  headerLinks,
		Footer: footer,
		Files:  Code,
	}

	tmpl_err := tmpl.Execute(w, data)
	handleError("template execute error:", tmpl_err)
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

	home_www := http.HandlerFunc(templateHandler)
	http.HandleFunc("/", home_www)

	about := http.HandlerFunc(templateHandlerAbout)
	http.HandleFunc("/about", about)

	history := http.HandlerFunc(templateHandlerHistory)
	http.HandleFunc("/history", history)

	fmt.Println("Server: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
