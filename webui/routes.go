package webui

import (
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

//go:generate qtc -dir=templates

// XXX TODO XXX TODO: double check the output tags I've been using

type WebUI struct {
	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func New(pgxPool *pgxpool.Pool) *WebUI {
	return &WebUI{pgxPool: pgxPool, dbQuerier: db.New()}
}

func (w *WebUI) Handler() http.Handler {
	// mux := http.NewServeMux()

	r := chi.NewRouter()

	r.Use(LoggerMiddleware(&log.Logger))

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/tests", http.StatusTemporaryRedirect)
	// })

	r.Route("/api", func(r chi.Router) {

		r.Route("/tests", func(r chi.Router) {
			r.Get("/", w.listTests)
			r.Get("/{testID}", w.getTest)
		})

		r.Route("/runs", func(r chi.Router) {
			r.Get("/", w.listRuns)
			r.Get("/{runID}", w.getRun)
		})

		r.Route("/runners", func(r chi.Router) {
			r.Get("/", w.listRunners)
			r.Get("/{runnerID}", w.getRunner)
		})

	})

	fs := http.FileServer(http.Dir("webui/app/dist"))
	r.Handle("/*", fs)

	return r
}
