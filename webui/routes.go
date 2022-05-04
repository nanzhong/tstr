package webui

import (
	"context"
	"net/http"

	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/webui/templates"
	"github.com/rs/zerolog/hlog"
)

//go:generate qtc -dir=templates

// XXX TODO XXX TODO: double check the output tags I've been using

type WebUI struct {
	querier db.Querier
}

func NewWebUI(querier db.Querier) *WebUI {
	return &WebUI{querier}
}

func (ui *WebUI) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().Msg("serving dashboard")
	p := &templates.Dashboard{}
	templates.WritePageTemplate(w, p)
}

func (ui *WebUI) ServeLabels(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().Msg("serving labels page")

	testByLabel, err := ui.querier.UITestsByLabels(context.Background())

	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("unable to UITestsByLabels")
		w.WriteHeader(500)
		return
	}

	resultsByTestRows, err := ui.querier.UITestResults(context.Background())
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("unable to UITestResults")
		w.WriteHeader(500)
		return
	}

	resultsByTest := map[string][]db.RunResult{}
	for _, r := range resultsByTestRows {
		resultsByTest[r.TestID] = r.Results
	}

	p := &templates.LabelsPage{
		TestsByLabel:  testByLabel,
		ResultsByTest: resultsByTest,
	}
	templates.WritePageTemplate(w, p)
}

func (ui *WebUI) ServeRuns(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().Msg("serving runs page")

	runs, err := ui.querier.UIListRecentRuns(context.Background(), 100)

	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("unable to ListRuns")
		w.WriteHeader(500)
		return
	}

	hasPendingRuns := false
	hasFinishedRuns := false

	for _, run := range runs {
		if run.FinishedAt.Status == pgtype.Null {
			hasPendingRuns = true
		} else {
			hasFinishedRuns = true
		}
	}

	p := &templates.RunsPage{
		Runs:            runs,
		HasPendingRuns:  hasPendingRuns,
		HasFinishedRuns: hasFinishedRuns,
	}
	templates.WritePageTemplate(w, p)
}

func (w *WebUI) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", w.ServeDashboard)

	mux.HandleFunc("/runs", w.ServeRuns)

	mux.HandleFunc("/labels", w.ServeLabels)

	// TODO: fix me. please.
	return hlog.MethodHandler("method")(hlog.URLHandler("url")(hlog.RemoteAddrHandler("peer")(mux)))
}
