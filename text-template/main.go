package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var _ = oldmain

func oldmain() {
	// type Inventory struct {
	// 	Material string
	// 	Count    uint
	// }

	//sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.}} items are made of {{.}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, "foo")
	if err != nil {
		panic(err)
	}

	lp := filepath.Join("templates", "layout.html")

	t, err := template.ParseFiles(lp)
	if err != nil {
		panic(err)
	}

	println(11, t.Name())

	err = t.ExecuteTemplate(os.Stdout, "layout.html", "wow")
	if err != nil {
		panic(err)
	}


}

func main() {

	log.SetFlags(log.Lshortfile)

	f := os.DirFS("")
	a := FileServer(f)

	err := http.ListenAndServe(":8080", a)
	panic(err)

}

type fileHandler struct {
	root fs.FS
}

func FileServer(root fs.FS) http.Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	// MIT license Alex Edwards
	lp := filepath.Join("templates", "master.html")
	fp := filepath.Join(filepath.Clean(r.URL.Path))
	//Return a 404 if the template doesn't exist

	info, err := fs.Stat(f.root, fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(rw, r)
			return
		}
		log.Println(err.Error())
		http.Error(rw, http.StatusText(500), 500)
		return
	}
	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(rw, r)
		return
	}

	body, err := fs.ReadFile(f.root, fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, http.StatusText(500), 500)
		return
	}

	// we could move this to creation-time, for efficiency
	//but we leave it here for development easy, we can change template while running
	tmpl, err := template.ParseFS(f.root, lp)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(rw, "master", body)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, http.StatusText(500), 500)
	}
}
