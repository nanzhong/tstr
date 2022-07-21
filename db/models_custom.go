package db

import "github.com/jackc/pgtype"

type RunLog struct {
	Type string `json:"type"`
	Time string `json:"time"`
	Data []byte `json:"data"`
}

type RunResults []RunResult

func (r *RunResults) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	var dec pgtype.TextArray
	if err := dec.DecodeText(ci, src); err != nil {
		return err
	}
	for _, el := range dec.Elements {
		*r = append(*r, RunResult(el.String))
	}
	return nil
}

func (r RunResults) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	var enc pgtype.TextArray
	if err := enc.Set(r); err != nil {
		return nil, err
	}
	return enc.EncodeText(ci, buf)
}
