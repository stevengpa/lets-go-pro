package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/stevengpa/lets-go-pro/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger                  `json:"logger,omitempty"`
	snippets       *models.SnippetModel          `json:"snippets,omitempty"`
	users          *models.UserModel             `json:"users,omitempty"`
	templateCache  map[string]*template.Template `json:"template_cache,omitempty"`
	formDecoder    *form.Decoder                 `json:"form_decoder,omitempty"`
	sessionManager *scs.SessionManager           `json:"session_manager,omitempty"`
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Access")
	dsn := flag.String("dsn", "web:pass@tcp(localhost:3386)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		//AddSource: true,
	}))

	db, err := openDSN(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:                *addr,
		Handler:             app.routes(),
		ErrorLog:            slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:           tlsConfig,
		IdleTimeout:         time.Minute,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        10 * time.Second,
		MaxHeaderValueCount: 100,
	}

	logger.Info("starting server", slog.String("addr", srv.Addr))
	err = srv.ListenAndServeTLS("./assets/tls/cert.pem", "./assets/tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDSN(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		err2 := db.Close()
		if err2 != nil {
			return nil, err2
		}
		return nil, err
	}

	return db, nil
}
