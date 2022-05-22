package webui

import (
	"context"
	"net/http"

	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/webui/templates"
	"github.com/rs/zerolog/log"
)

func (ui *WebUI) listTests(w http.ResponseWriter, r *http.Request) {
	tests, err := ui.dbQuerier.UIListTests(context.Background(), ui.pgxPool)
	if err != nil {
		log.Err(err).Msg("unable to ListTests")
		w.WriteHeader(500)
		return
	}

	if 1 == 2 {
		println(tests)
	}

	runs, err := ui.dbQuerier.UIListRecentRuns(context.Background(), ui.pgxPool, 500)
	if err != nil {
		log.Err(err).Msg("unable to ListRuns")
		w.WriteHeader(500)
		return
	}

	runsByTestID := map[string][]db.UIListRecentRunsRow{}

	for _, run := range runs {
		if _, exists := runsByTestID[run.TestID.String()]; exists {
			runsByTestID[run.TestID.String()] = append(runsByTestID[run.TestID.String()], run)
		} else {
			runsByTestID[run.TestID.String()] = []db.UIListRecentRunsRow{run}
		}
	}

	p := &templates.TestsPage{
		Tests:            tests,
		TestRunsByTestID: runsByTestID,
	}
	templates.WritePageTemplate(w, p)

}
