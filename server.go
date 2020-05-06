package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/russross/blackfriday/v2"
)

// Verb - Route            - File on server           - Description
//
// GET  - /                - ./index.html             - home
// GET  - /:path.ext       - ./static/:path.ext       - static files
//
// GET  - /til             - ./til/index.html         - til index
// GET  - /til.xml         - ./til.xml                - til index rss
// GET  - /til/:path       - ./til/:path/index.html   - til subindex
// GET  - /til/:path       - ./til/:path.html         - til post content
// GET  - /til/:path.md    - ./til/:path.md           - til post source
//
// GET  - /blog            - ./blog/index.html        - blog index
// GET  - /blog.xml        - ./blog.xml               - blog index rss
// GET  - /blog/:path      - ./blog/:path/index.html  - blog subindex
// GET  - /blog/:path      - ./blog/:path.html        - blog post content
// GET  - /blog/:path.md   - ./blog/:path.md          - blog post source
//
// GET  - /files           - ./files/index.html       - files index
// GET  - /files/:path     - ./files/:path/index.html - files subindex
// GET  - /files/:path.ext - ./files/:path.ext        - file
// GET  - /upload          - ./upload.html            - upload files web gui
// POST - /upload          - n/a                      - upload files
//
// GET  - /share           - n/a                      - webrtc file share
// GET  - /chess           - n/a                      - webrtc chess
//
// /ping - pong

func main() {
	fmt.Println("starting go server")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)

	// r.Route("/til", func(r chi.Router) {
	// 	r.Get("/", handler)
	// 	r.Get("/*", handler)
	// })

	r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/files/", http.FileServer(http.Dir("./files")))
		fs.ServeHTTP(w, r)
	})

	r.NotFound(notFoundHandler)

	http.ListenAndServe(":3000", r)
}

// -----------------------------------------------------------------------------

func handler(w http.ResponseWriter, r *http.Request) {
	urlFormat := getUrlFormat(r)

	switch urlFormat {
	case "md":
		path := "." + r.RequestURI
		if !postExists(path) {
			notFoundHandler(w, r)
			return
		}
		md, err := ioutil.ReadFile(path)
		if err != nil {
			internalServerError(w, r)
			return
		}
		w.Write(md)
	default:
		path := "." + r.RequestURI + ".md"
		if !postExists(path) {
			notFoundHandler(w, r)
			return
		}

		metadata, err := parseMetadata(path)
		if err != nil {
			internalServerError(w, r)
			return
		}

		md, err := ioutil.ReadFile(path)
		if err != nil {
			internalServerError(w, r)
			return
		}
		body := parseMarkdown(md)

		post := Post{metadata, body}

		t, err := template.ParseFiles("templates/post.html")
		if err != nil {
			internalServerError(w, r)
			return
		}

		t.Execute(w, post)
	}
}

func getUrlFormat(r *http.Request) string {
	return r.Context().Value(middleware.URLFormatCtxKey).(string)
}

func postExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && filepath.Ext(path) == ".md"
}

type Post struct {
	Metadata PostMetadata
	Body     template.HTML
}

type PostMetadata struct {
	Title string
	Date  time.Time
	Size  int64
}

func parseMarkdown(md []byte) template.HTML {
	return template.HTML(blackfriday.Run(md))
}

func parseMetadata(path string) (PostMetadata, error) {
	info, err := os.Stat(path)
	if err != nil {
		return PostMetadata{}, err
	}

	fn := info.Name()
	fnWithoutExtension := strings.TrimSuffix(fn, filepath.Ext(fn))
	fnWithoutDashes := strings.ReplaceAll(fnWithoutExtension, "-", " ")
	fnWithoutUnderlines := strings.ReplaceAll(fnWithoutDashes, "_", " ")
	title := strings.Title(fnWithoutUnderlines)
	date := info.ModTime()
	size := info.Size()

	return PostMetadata{title, date, size}, nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("nothing here"))
	// http.NotFound(w, r)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte("internal server error\nplease contact me b@bartol.dev"))
}
