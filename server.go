package main

import (
	"bufio"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// routes //////////////////////////////////////////////////////////////////////

// / 				- index page
// t: index.html, head.html, header.html, footer.html

// /memory/ 		- memory index page
// t: memory.html, head.html, header.html, footer.html

// /memory/:path 	- memory post page
// t: post.html, head.html, header.html, footer.html

// /memory/:path.md - memory post source page
// t: n/a

// /memory.xml 		- memory index rss feed
// t: memory.xml

// /paste/ 			- paste index page
// t: paste.html, head.html

// /paste/:path 	- view paste
// t: n/a

// /paste/flush 	- flush paste table
// t: n/a

// /upload/ 		- upload index page
// t: upload.html, head.html

// /upload/:path 	- download upload
// t: n/a

// /upload/flush 	- flush paste table
// t: n/a

// /ping 			- pong
// t: n/a

// ...

// GET  - /share           - webrtc file share
// GET  - /chess           - webrtc chess

// GET  - /color           - color image api
// GET  - /random          - generate random number

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./web.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS
		paste (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			content TEXT,
			date TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS
		upload (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			content BLOB,
			date TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", indexHandler)

	// http.HandleFunc("/memory/", memoryHandler)
	// http.HandleFunc("/memory.xml", memoryFeedHandler)

	serveFile("/favicon.ico")
	serveFile("/robots.txt")
	serveDir("/css/")
	serveDir("/js/")
	serveDir("/files/")

	// http.HandleFunc("/paste/", pasteHandler)
	// flushTable("paste")

	// http.HandleFunc("/upload/", uploadHandler)
	// flushTable("upload")

	// http.HandleFunc("/ping", pingHandler)

	log.Fatal(http.ListenAndServe(":8080", logRequest(http.DefaultServeMux)))
}

// handlers ////////////////////////////////////////////////////////////////////

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	page := indexData{
		MetaTitle: "Bartol Deak",
	}

	err := indexTemplates.Execute(w, page)
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
	}
}

// templates ///////////////////////////////////////////////////////////////////

var indexTemplates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/meta.html",
	"templates/header.html",
	"templates/footer.html",
))

// var memoryTemplates = template.Must(template.ParseFiles(
// 	"templates/memory.html",
// 	"templates/meta.html",
// 	"templates/header.html",
// 	"templates/footer.html",
// ))

// var postTemplates = template.Must(template.ParseFiles(
// 	"templates/post.html",
// 	"templates/meta.html",
// 	"templates/header.html",
// 	"templates/footer.html",
// ))

// var memoryFeedTemplates = template.Must(template.ParseFiles(
// 	"templates/memory.xml",
// ))

// var pasteTemplates = template.Must(template.ParseFiles(
// 	"templates/paste.html",
// 	"templates/meta.html",
// ))

// var uploadTemplates = template.Must(template.ParseFiles(
// 	"templates/upload.html",
// 	"templates/meta.html",
// ))

// template data ///////////////////////////////////////////////////////////////

type indexData struct {
	MetaTitle string
}

// type memoryData struct{}

// type postData struct{}

// type memoryFeedData struct {}

// type pasteData struct {}

// type uploadData struct {}

// old.handlers - delete ///////////////////////////////////////////////////////
/*
type Page struct {
	Title       string
	Stylesheets []string
	Scripts     []string
}

var indexTemplates = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

	page := Page{
		Title: "Bartol Deak",
	}

	err := indexTemplates.Execute(w, page)
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
	}
}

type Post struct {
	Title string
	Path  string
}

type PostPage struct {
	Page
	Content template.HTML
}

type PostIndexPage struct {
	Page
	Posts []Post
}

var postTemplates = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/post.html",
))

var postIndexTemplates = template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/post_index.html",
))

func tilHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/til/"):]
	if path != "" {
		realPath := getRealPath("./til/" + path)
		if !isPost(realPath) {
			notFoundHandler(w, r)
			return
		}

		md, err := ioutil.ReadFile(realPath)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}

		if strings.HasSuffix(path, ".md") {
			w.Write(md)
			return
		}

		title, err := getTitle(realPath)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
		}

		renderer := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
		var buf bytes.Buffer
		err = renderer.Convert(md, &buf)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
		}

		content := template.HTML(buf.Bytes())

		page := PostPage{
			Page{
				Title: title + " :: Today I Learned :: Bartol Deak",
			},
			content,
		}

		err = postTemplates.Execute(w, page)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
		}
		return
	}

	var posts []Post
	err := filepath.Walk("./til/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isPost(path) {
			title, err := getTitle(path)
			if err != nil {
				return err
			}

			post := Post{
				Title: title,
				Path:  "/" + path[:len(path)-len(".md")] + "/",
			}

			posts = append(posts, post)
		}

		return nil
	})
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
	}

	page := PostIndexPage{
		Page{
			Title: "Today I Learned :: Bartol Deak",
		},
		posts,
	}

	err = postIndexTemplates.Execute(w, page)
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
	}
}

func tilFeedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("til feed"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query != "" {
		w.Write([]byte("search query: " + query))
		return
	}

	w.Write([]byte("enter search query"))
}

type item struct {
	ID   int
	Name string
	Date string
}

type pastePage struct {
	Items []item
}

func pasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		content := r.FormValue("content")

		if name == "" || content == "" {
			w.Write([]byte("bad request"))
			return
		}

		result, err := db.Exec(`
			INSERT INTO
				paste (
					name,
					content,
					date
				)
				VALUES (
					?,
					?,
					DATETIME('NOW')
				)
		`, name, content)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}
		location := "/paste/" + strconv.FormatInt(id, 10)
		w.Header().Add("Location", location)
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	id := r.URL.Path[len("/paste/"):]
	if id != "" {
		row := db.QueryRow("SELECT content FROM paste WHERE id=?", id)
		var content string
		err := row.Scan(&content)
		if err != nil {
			notFoundHandler(w, r)
			return
		}
		w.Write([]byte(content))
		return
	}

	page := struct {
		Page
		Items []item
	}{
		Page: Page{
			Title: "Paste",
		},
	}

	rows, err := db.Query("SELECT id,name,date FROM paste ORDER BY id DESC")
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id   int
			name string
			date string
		)
		err := rows.Scan(&id, &name, &date)

		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}
		page.Items = append(page.Items, item{id, name, date})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/header.html", "templates/footer.html", "templates/paste.html")
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
		return
	}
	tmpl.Execute(w, page)
}

type uploadPage struct {
	Items []item
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, metadata, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("bad request" + err.Error()))
			return
		}
		defer file.Close()
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, file)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}
		name := metadata.Filename
		content := buf.Bytes()

		if metadata.Size > 20000000 {
			ok := basicAuth(w, r)
			if !ok {
				http.Error(w, "Not authorized", 401)
				return
			}
		}

		_, err = db.Exec(`
			INSERT INTO
				upload (
					name,
					content,
					date
				)
				VALUES (
					?,
					?,
					DATETIME('NOW')
				)
		`, name, content)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}

		w.Header().Add("Location", "/upload/")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	name := r.URL.Path[len("/upload/"):]
	if name != "" {
		row := db.QueryRow("SELECT content FROM upload WHERE name=?", name)
		var content string
		err := row.Scan(&content)
		if err != nil {
			notFoundHandler(w, r)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
		w.Write([]byte(content))
		return
	}

	page := pastePage{}

	rows, err := db.Query("SELECT id,name,date FROM upload ORDER BY id DESC")
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id   int
			name string
			date string
		)
		err := rows.Scan(&id, &name, &date)

		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}
		page.Items = append(page.Items, item{id, name, date})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		w.Write([]byte("internal server error" + err.Error()))
		return
	}
	tmpl.Execute(w, page)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ping-Pong, Umjetnost Zdravog Đira"))
}
*/

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404"))
}

// internalServerErrorHandler
// unauthorizedHandler

// utils ///////////////////////////////////////////////////////////////////////

func serveFile(path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "."+path)
	})
}

func serveDir(path string) {
	fs := http.FileServer(http.Dir("." + path))
	http.Handle(path, http.StripPrefix(path, fs))
}

func flushTable(table string) {
	path := "/" + table + "/flush"
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ok := basicAuth(w, r)
		if !ok {
			http.Error(w, "Not authorized", 401)
			return
		}

		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			w.Write([]byte("internal server error" + err.Error()))
			return
		}

		w.Write([]byte("done"))
	})
}

func basicAuth(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	// env[pw]...
	return authOK && username == "usr" && password == "pw"
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func getRealPath(path string) string {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1] + ".md"
	}

	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}

	return path
}

func isPost(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && filepath.Ext(path) == ".md"
}

func getTitle(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return firstLine[len("# ") : len(firstLine)-1], nil
}
