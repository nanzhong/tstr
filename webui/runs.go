package webui

import (
	"context"
	"net/http"

	"github.com/nanzhong/tstr/webui/templates"
	"github.com/rs/zerolog/log"
)

func (ui *WebUI) listRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := ui.dbQuerier.UIListRecentRuns(context.Background(), ui.pgxPool, 100)

	if err != nil {
		log.Err(err).Msg("unable to ListRuns")
		w.WriteHeader(500)
		return
	}

	hasPendingRuns := false
	hasFinishedRuns := false

	// the UIListRecentRuns query sorts by pending, so:
	if len(runs) > 0 {
		if runs[0].IsPending {
			hasPendingRuns = true
		}
		if runs[len(runs)-1].IsPending {
			hasFinishedRuns = false
		}
	} else {
		hasPendingRuns = false
		hasFinishedRuns = false
	}

	for _, run := range runs {
		if !run.FinishedAt.Valid {
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
