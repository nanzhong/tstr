package webui

import (
	"net/http"

	"github.com/nanzhong/tstr/webui/templates"
	"github.com/rs/zerolog/hlog"
)

//go:generate qtc -dir=templates


// XXX TODO XXX TODO: double check the output tags I've been using

type WebUI struct{}

func NewWebUI() *WebUI {
	return &WebUI{}
}

func (ui *WebUI) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().Msg("serving dashboard")
	p := &templates.Dashboard{}
	templates.WritePageTemplate(w, p)
}

func (ui *WebUI) ServeRuns(w http.ResponseWriter, r *http.Request) {

	hlog.FromRequest(r).Info().Msg("serving runs page")
	runs := genRuns()
	hasPendingRuns := false
	hasFinishedRuns := false

	for _, run := range runs {
		if run.FinishedAt == nil {
			hasPendingRuns = true
			break
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

	mux.HandleFunc("/labels", func(w http.ResponseWriter, r *http.Request) {
		hlog.FromRequest(r).Info().Msg("serving labels page")

		lbl1 := templates.LabelSummary{
			Label:       "xlabel: zvalue",
			MonthReport: genDayReports(60),
		}

		monthsummary := []templates.LabelSummary{lbl1}

		p := &templates.LabelsPage{monthsummary}
		templates.WritePageTemplate(w, p)
	})

	// TODO: fix me. please.
	return hlog.MethodHandler("method")(hlog.URLHandler("url")(hlog.RemoteAddrHandler("peer")(mux)))
}
