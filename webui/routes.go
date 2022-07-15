package webui

import (
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

type WebUI struct {
	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func New(pgxPool *pgxpool.Pool) *WebUI {
	return &WebUI{pgxPool: pgxPool, dbQuerier: db.New()}
}

func (w *WebUI) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(LoggerMiddleware(&log.Logger))

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
