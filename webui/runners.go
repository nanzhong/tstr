package webui

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"
)

func (ui *WebUI) listRunners(w http.ResponseWriter, r *http.Request) {
	runners, err := ui.dbQuerier.ListRunners(context.Background(), ui.pgxPool, sql.NullTime{Valid: true, Time: time.Now().Add(-24 * 365 * time.Hour)}) // XXX TODO FIXME

	if err != nil {
		log.Err(err).Msg("unable to ListRunners")
		w.WriteHeader(500)
		return
	}

	render.JSON(w, r, runners)
}

func (ui *WebUI) getRunner(w http.ResponseWriter, r *http.Request) {

	runnerIDstr := chi.URLParam(r, "runnerID")
	runnerID, err := uuid.Parse(runnerIDstr)
	if err != nil {
		log.Error().Err(err).Str("runnerID", runnerIDstr).Msg("unable to parse runnerID")
		w.WriteHeader(400)
		return
	}

	limit := 100

	numRuns := r.URL.Query().Get("runs")
	if numRuns != "" {
		limit, err = strconv.Atoi(numRuns)
		if err != nil {
			log.Error().Err(err).Str("runs", numRuns).Msg("unable to parse runs")
			w.WriteHeader(400)
			return
		}
	}

	runner, err := ui.dbQuerier.GetRunner(context.Background(), ui.pgxPool, runnerID)
	if err != nil {
		log.Err(err).Msg("unable to GetRunner")
		w.WriteHeader(500)
		return
	}

	runsSummary, err := ui.dbQuerier.UIRunnerSummary(context.Background(), ui.pgxPool, db.UIRunnerSummaryParams{
		RunnerID: uuid.NullUUID{UUID: runnerID, Valid: true},
		Limit:    int32(limit),
	})
	if err != nil {
		log.Err(err).Msg("unable to ListRunners")
		w.WriteHeader(500)
		return
	}

	response := map[string]interface{}{
		"Runner":      runner,
		"RunsSummary": runsSummary,
	}

	if 2 == 1 {
		render.JSON(w, r, map[string]interface{}{
			"Runner":      runner,
			"RunsSummary": []string{},
		})

	} else {

		render.JSON(w, r, response)
	}
}
