package templates

import (
	"fmt"

	"github.com/jackc/pgtype"
)

func jsonbToLabels(jobj pgtype.JSONB) map[string]string {
	m := map[string]string{}

	err := jobj.AssignTo(&m)
	if err != nil {
		panic(err)
	}
	return m
}

func labelsAsSlice(jobj pgtype.JSONB, separator string) []string {
	s := []string{}
	for k, v := range jsonbToLabels(jobj) {
		s = append(s, fmt.Sprintf("%s%s%s", k, separator, v))
	}

	return s
}
