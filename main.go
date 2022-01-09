package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	ImagePath string
	Title     string
	Body      []byte
}

func (p *Page) save() error {
	path := "./data/"
	filename := path + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
func (p *Page) delete() error{
	path := "./data/"
	filename := path + p.Title + ".txt"
	return os.Remove(filename)
}

func loadPage(title string) (*Page, error) {
	path := "./data/"
	filename := path + title + ".txt"
	imagePath := "/assets/" + title + ".jpg"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body, ImagePath: imagePath}, nil
}



func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// if title == home, then execute diff code for home page
	if title == "home" {
		pages, err := ioutil.ReadDir("./data/")
		if err != nil {
			log.Fatal(err)
		}
		pagesToRend := make([]string,len(pages))
		for i,page := range pages{
			pageName := page.Name()
			pagesToRend[i] = strings.Trim(pageName,".txt")
		}
		templates.ExecuteTemplate(w,"view.html",pagesToRend)
		
  }else{
	   p, err := loadPage(title)
	   if err != nil {
		    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		    return
	     }
     renderTemplate(w, "view", p)
  }
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	var p *Page
	// if method post
	if r.Method == http.MethodPost{
		title := r.FormValue("page")
		p = &Page{Title: title}
	}
	// if method get
	if r.Method == http.MethodGet{
		title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// check if the page is the home page
	if title == "home"{
		http.Redirect(w,r,"/view/home",http.StatusForbidden)
		return
	}

	p, err = loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadFile(w, r, title)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func deleteHandler(w http.ResponseWriter, r *http.Request){
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
		err = p.delete()
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		
		err = os.Remove("./assets/"+title+".jpg")
		if err != nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		http.Redirect(w,r,"/view/home",http.StatusFound)
}

func uploadFile(w http.ResponseWriter, r *http.Request, pageName string) {
	r.ParseMultipartForm(32<<20)
	file, _, err := r.FormFile("image")
	if file == nil {
		fmt.Println("empty file")
		return
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("./assets/"+pageName+".jpg", fileBytes, 0644)
}

//creating valid path array
var validPathArr [4] string = [4]string{"/view/","/save/","/edit/","/delete/"}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	path := r.URL.Path
	for i := range validPathArr{
		if strings.HasPrefix(path,validPathArr[i]) == true {
			return strings.TrimPrefix(path,validPathArr[i]),nil
		}
	}
	return "", errors.New("Invalid Page Request")
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/view/", viewHandler)
	mux.HandleFunc("/edit/", editHandler)
	mux.HandleFunc("/save/", saveHandler)
	mux.HandleFunc("/delete/",deleteHandler)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
