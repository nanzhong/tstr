package webui

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (ui *WebUI) getTest(w http.ResponseWriter, r *http.Request) {

	testIDstr := chi.URLParam(r, "testID")
	testID, err := uuid.Parse(testIDstr)
	if err != nil {
		log.Error().Err(err).Str("runID", testIDstr).Msg("unable to parse testID")
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

	// testIDs :=

	// testID = uuid.NullUUID{
	// 	UUID:  testID,
	// 	Valid: true,
	// }

	test, err := ui.dbQuerier.GetTest(context.Background(), ui.pgxPool, testID)
	if err != nil {
		log.Err(err).Msg("unable to GetTest")
		w.WriteHeader(500)
		return
	}

	runs, err := ui.dbQuerier.UIRunsSummary(context.Background(), ui.pgxPool,
		db.UIRunsSummaryParams{
			TestID: uuid.NullUUID{
				UUID:  testID,
				Valid: true,
			},
			Limit: int32(limit),
		})

	if err != nil {
		log.Err(err).Msg("unable to ListRuns")
		w.WriteHeader(500)
		return
	}

	resp := struct {
		db.GetTestRow
		RunsSummary []db.UIRunsSummaryRow
	}{
		test,
		runs,
	}

	// response := map[string]interface{}{
	// 	"RunsSummary": runs,
	// 	"Test":        test,
	// }

	render.JSON(w, r, resp)
}

func (ui *WebUI) listTests(w http.ResponseWriter, r *http.Request) {
	tests, err := ui.dbQuerier.UIListTests(context.Background(), ui.pgxPool)
	if err != nil {
		log.Err(err).Msg("unable to ListTests")
		w.WriteHeader(500)
		return
	}

	render.JSON(w, r, tests)

}
