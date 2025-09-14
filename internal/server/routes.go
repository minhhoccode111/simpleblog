package server

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Page struct {
	Title   string
	Slug    string
	Body    []byte
	PubDate time.Time
}

type Metadata struct {
	Title, Slug string
	PubDate     string
}

const timeFormat = "2006-01-02"

var templates = template.Must(template.ParseGlob("template/*.html"))
var linkRegex = regexp.MustCompile(`\[(.*?)\]`)
var validPath = regexp.MustCompile(`^/(admin/)?(articles(/([a-zA-Z0-9-]+))?/?(\?.*)?)?$`)

func (p *Page) save() error {
	filename := "data/" + p.Slug + ".md"
	var body []byte
	metadata := []string{
		"---",
		p.Title,
		p.PubDate.Format(timeFormat),
		"---",
	}

	body = append(body, []byte(strings.Join(metadata, "\n")+"\n")...) // 1 extra new line
	body = append(body, p.Body...)

	return os.WriteFile(filename, body, 0600)
}

// loadMetadata works like loadPage but returns Metadata, ignore the Body for performance
func loadMetadata(slug string) (*Metadata, error) {
	filename := "data/" + slug + ".md"

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var title string
	var pubDate time.Time

	// read the first 4 lines
	for range 4 {
		line, err := reader.ReadString('\n')
		line = line[:len(line)-1]
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if line == "---" {
			continue
		}
		if title == "" {
			title = line
			continue
		}

		lineTime, err := time.Parse(timeFormat, line)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		pubDate = lineTime
	}

	return &Metadata{title, slug, pubDate.Format(timeFormat)}, nil
}

func loadPage(slug string) (*Page, error) {
	filename := "data/" + slug + ".md"
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var title string
	var pubDate time.Time

	// read the first 4 lines
	for range 4 {
		line, err := reader.ReadString('\n')
		line = line[:len(line)-1]
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		if line == "---" {
			continue
		}
		if title == "" {
			title = line
			continue
		}

		lineTime, err := time.Parse(timeFormat, line)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		pubDate = lineTime
	}

	rest, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Page{title, slug, rest, pubDate}, nil
}

func renderErrorTemplate(w http.ResponseWriter, err error) {
	// try to render error
	err = templates.ExecuteTemplate(w, "error.html", err.Error())
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error executing: %v", err),
			http.StatusInternalServerError,
		)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	var err error

	switch tmpl {
	case "view":
		// parse markdown before viewing
		var buf bytes.Buffer
		if err = goldmark.Convert(p.Body, &buf); err != nil {
			break
		}
		err = templates.ExecuteTemplate(w, "view.html", struct {
			Title, Slug, PubDate string
			Body                 template.HTML
		}{p.Title, p.Slug, p.PubDate.Format(timeFormat), template.HTML(buf.String())})
	case "edit":
		err = templates.ExecuteTemplate(w, "edit.html", struct {
			Title, Slug, PubDate string
			Body                 []byte
		}{p.Title, p.Slug, p.PubDate.Format(timeFormat), p.Body})
	case "all-published":
		today := time.Now().Format(timeFormat)
		files, err := os.ReadDir("./data")
		if err != nil {
			break
		}
		metadataSlice := []*Metadata{}
		for _, f := range files {
			slug := strings.TrimSuffix(f.Name(), ".md")
			metadata, err := loadMetadata(slug)
			if err != nil {
				// because we can't use `break` to exit the switch so we render error view right here
				renderErrorTemplate(w, err)
				return
			}
			if metadata.PubDate > today {
				continue
			}
			metadataSlice = append(metadataSlice, metadata)
		}
		err = templates.ExecuteTemplate(w, "all-published.html", metadataSlice)
	case "all-admin":
		files, err := os.ReadDir("./data")
		if err != nil {
			break
		}
		metadataSlice := []*Metadata{}
		for _, f := range files {
			slug := strings.TrimSuffix(f.Name(), ".md")
			metadata, err := loadMetadata(slug)
			if err != nil {
				renderErrorTemplate(w, err)
				return
			}
			metadataSlice = append(metadataSlice, metadata)
		}
		err = templates.ExecuteTemplate(w, "all-admin.html", metadataSlice)
	}

	if err != nil {
		renderErrorTemplate(w, err)
	}
}

func unauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (s *Server) basicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorizedResponse(w)
			return
		}

		encodedBase64 := authHeader[len("Basic "):]
		payload, err := base64.StdEncoding.DecodeString(encodedBase64)
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if parts[0] != "admin" || parts[1] != "admin" {
			unauthorizedResponse(w)
			return
		}

		// log.Printf("logged in: %v", parts)
		next.ServeHTTP(w, r)
	})
}

// public
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/articles", http.StatusFound)
}
func (s *Server) GetAllPublishedArticlesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "all-published", nil)
}
func (s *Server) GetPublishedArticleHandler(w http.ResponseWriter, r *http.Request, slug string) {
	page, err := loadPage(slug)
	if err != nil && os.IsNotExist(err) {
		http.Redirect(
			w,
			r,
			fmt.Sprintf("/admin/articles?action=create&slug=%s", slug),
			http.StatusSeeOther,
		)
		return
	}

	if err != nil {
		renderErrorTemplate(w, err)
		return
	}

	today := time.Now().Format(timeFormat)
	if page.PubDate.Format(timeFormat) > today {
		renderErrorTemplate(w, fmt.Errorf("Unpublished Article"))
		return
	}

	renderTemplate(w, "view", page)
}

// auth
func (s *Server) AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/articles", http.StatusFound)
}
func (s *Server) AdminGetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: if action=create, which mean now article found
	// if slug=provided-slug, then redirect to `edit` since user already have a specific slug in mind
	// else
	// if slug=some-slug, pre-fill slug to form (or we should redirect to edit handler?)
	// else display all articles
	renderTemplate(w, "all-admin", nil)
}
func (s *Server) AdminCreateArticleHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AdminUpdateArticleGetHandler(w http.ResponseWriter, r *http.Request, slug string) {
	page, err := loadPage(slug)
	// error occurs and not not found file error
	if err != nil {
		if !os.IsNotExist(err) {
			renderErrorTemplate(w, err)
			return
		}

		// not found file, redirect to create
		http.Redirect(
			w,
			r,
			fmt.Sprintf("/admin/articles?action=create&slug=%s", slug),
			http.StatusSeeOther,
		)
		return
	}

	renderTemplate(w, "edit", page)
}
func (s *Server) AdminUpdateArticleHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) AdminDeleteArticleHandler(w http.ResponseWriter, r *http.Request) {}

// makeHandler act like a middleware that extract slug from Path and pass it to
// handler, we might not need to use it since mux is good enough
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		fn(w, r, slug)
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// serve static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", s.IndexHandler).Methods("GET")
	r.HandleFunc("/articles", s.GetAllPublishedArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{slug}", makeHandler(s.GetPublishedArticleHandler)).Methods("GET")

	// Admin subrouter with authentication
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(s.basicAuthentication)

	adminRouter.HandleFunc("", s.AdminIndexHandler).Methods("GET")
	adminRouter.HandleFunc("/articles", s.AdminGetAllArticlesHandler).Methods("GET")
	adminRouter.HandleFunc("/articles", s.AdminCreateArticleHandler).Methods("POST")
	// assume that admin get to edit, so return edit view
	adminRouter.HandleFunc("/articles/{slug}", makeHandler(s.AdminUpdateArticleGetHandler)).
		Methods("GET")
	// with ?action=edit because a html form can't send PUT request
	adminRouter.HandleFunc("/articles/{slug}", s.AdminUpdateArticleHandler).Methods("POST")
	// with ?action=delete because a html form can't send DELETE request
	adminRouter.HandleFunc("/articles/{slug}", s.AdminDeleteArticleHandler).Methods("POST")
	return r
}
