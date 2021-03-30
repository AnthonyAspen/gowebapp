package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}


func (p *Page) save() error {
	path := "./data/"
	filename := path + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	path := "./data/"
	filename := path + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title,err := getTitle(w,r)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title,err := getTitle(w,r)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title,err := getTitle(w,r)
	if err != nil{
		http.Error(w,err.Error(),http.StatusNotFound)
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {

	  path := r.URL.Path
		if strings.HasPrefix(path,"/view/")||strings.HasPrefix(path,"/save/")||strings.HasPrefix(path,"/edit/"){
			validPath := path[6:]
    return validPath, nil 
		}
			return "",errors.New("Invalid Page Request")
}
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is a normal way to find the valid path, but I changed it to
// comprehend it better
//var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


func main() {
  fs := http.FileServer(http.Dir("assets"))
  http.Handle("/assets/",http.StripPrefix("/assets/",fs))
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
