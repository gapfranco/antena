package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"antena/config"
	"antena/internal/models"
	"antena/ui"
)

type templateData struct {
	CurrentYear     int
	IsAuthenticated bool
	Events          []*models.Event
	Centrals        []interface{} // Simplified for now
	SearchType      string
	SearchCentral   int
	CurrentPage     int
	TotalPages      int
	HasNextPage     bool
	HasPrevPage     bool
}

type application struct {
	events        *models.EventModel
	templateCache map[string]*template.Template
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	addr := cfg.Port
	if addr == "" {
		addr = ":4000"
	}

	if cfg.TursoURL == "" || cfg.TursoToken == "" {
		log.Fatal("TURSO_URL and TURSO_TOKEN must be set in antena.conf or as environment variables")
	}

	db, err := openDB(cfg.TursoURL, cfg.TursoToken)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		events:        &models.EventModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	log.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(primaryURL, authToken string) (*sql.DB, error) {
	u, err := url.Parse(primaryURL)
	if err != nil {
		return nil, fmt.Errorf("parse turso url: %w", err)
	}

	q := u.Query()
	q.Set("auth_token", authToken)
	u.RawQuery = q.Encode()

	db, err := sql.Open("libsql", u.String())
	if err != nil {
		return nil, fmt.Errorf("open remote: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping remote: %w", err)
	}

	return db, nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /{$}", app.eventList)
	mux.HandleFunc("GET /events", app.eventList)

	return mux
}

func (app *application) eventList(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := 12
	eventType := r.URL.Query().Get("event_type")
	centralID, _ := strconv.Atoi(r.URL.Query().Get("central_id"))

	events, err := app.events.All(page, pageSize, eventType, centralID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalEvents, err := app.events.Count(eventType, centralID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := templateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: true,
		Events:          events,
		SearchType:      eventType,
		SearchCentral:   centralID,
		CurrentPage:     page,
		TotalPages:      (totalEvents + pageSize - 1) / pageSize,
	}
	data.HasNextPage = data.CurrentPage < data.TotalPages
	data.HasPrevPage = data.CurrentPage > 1

	app.render(w, http.StatusOK, "events.html", data)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist", page), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"add": func(a, b int) int {
		return a + b
	},
	"formatUnixMs": func(ms int64) string {
		if ms == 0 {
			return ""
		}
		t := time.UnixMilli(ms)
		return t.Format("02/01/2006 15:04:05")
	},
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/pages/*.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
