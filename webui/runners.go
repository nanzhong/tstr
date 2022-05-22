package webui

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/nanzhong/tstr/webui/templates"
	"github.com/rs/zerolog/log"
)

func (ui *WebUI) listRunners(w http.ResponseWriter, r *http.Request) {
	runners, err := ui.dbQuerier.ListRunners(context.Background(), ui.pgxPool, sql.NullTime{Valid: true, Time: time.Now().Add(-24 * 30 * time.Hour)}) // XXX TODO FIXME

	if err != nil {
		log.Err(err).Msg("unable to ListRunners")
		w.WriteHeader(500)
		return
	}

	p := &templates.RunnersPage{
		Runners: runners,
	}
	templates.WritePageTemplate(w, p)
}
