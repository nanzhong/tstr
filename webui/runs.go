package webui

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (ui *WebUI) getRun(w http.ResponseWriter, r *http.Request) {
	runIDstr := chi.URLParam(r, "runID")

	runID, err := uuid.Parse(runIDstr)

	if err != nil {
		log.Error().Err(err).Str("runID", runIDstr).Msg("unable to parse runID")
		w.WriteHeader(400)
		return
	}

	run, err := ui.dbQuerier.GetRun(context.Background(), ui.pgxPool, runID)
	if err != nil {
		log.Error().Err(err).Str("runID", runIDstr).Msg("unable to get runID")
		w.WriteHeader(404)
		return
	}

	render.JSON(w, r, run)
}

func (ui *WebUI) listRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := ui.dbQuerier.UIListRecentRuns(context.Background(), ui.pgxPool, 100)

	if err != nil {
		log.Err(err).Msg("unable to ListRuns")
		w.WriteHeader(500)
		return
	}

	render.JSON(w, r, runs)
}
