package server

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoTimestamp(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if ts.Status == pgtype.Null {
		return nil
	}

	var t time.Time
	if err := ts.AssignTo(&t); err != nil {
		// This should never fail...
		log.Error().Err(err).Msg("converting pgtype.timestamptz to time.Time")
		return nil
	}

	return timestamppb.New(t)
}
