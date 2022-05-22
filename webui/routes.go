package webui

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/webui/templates"
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

func (ui *WebUI) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	p := &templates.Dashboard{}
	templates.WritePageTemplate(w, p)
}

func (ui *WebUI) ServeLabels(w http.ResponseWriter, r *http.Request) {

	testByLabel, err := ui.dbQuerier.UITestsByLabels(context.Background(), ui.pgxPool)

	if err != nil {
		log.Err(err).Msg("unable to UITestsByLabels")
		w.WriteHeader(500)
		return
	}

	resultsByTestRows, err := ui.dbQuerier.UITestResults(context.Background(), ui.pgxPool)
	if err != nil {
		log.Err(err).Msg("unable to UITestResults")
		w.WriteHeader(500)
		return
	}

	resultsByTest := map[string][]db.RunResult{}
	for _, _ = range resultsByTestRows {
		// resultsByTest[r.TestID.String()] = r.Results.([]string)
	}

	p := &templates.LabelsPage{
		TestsByLabel:  testByLabel,
		ResultsByTest: resultsByTest,
	}
	templates.WritePageTemplate(w, p)
}

func (w *WebUI) Handler() http.Handler {
	// mux := http.NewServeMux()

	r := chi.NewRouter()

	// r.Use(LoggerMiddleware(&log.Logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tests", http.StatusTemporaryRedirect)
	})

	r.Route("/runs", func(r chi.Router) {
		r.Get("/", w.listRuns)
	})

	r.Route("/tests", func(r chi.Router) {
		r.Get("/", w.listTests)
	})

	r.Route("/runners", func(r chi.Router) {
		r.Get("/", w.listRunners)
	})

	return r
}
