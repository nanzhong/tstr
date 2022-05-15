package db

type RunLog struct {
	Type string `json:"type"`
	Time string `json:"time"`
	Data []byte `json:"data"`
}
