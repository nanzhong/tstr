package server

import (
	"github.com/jackc/pgx/v4/pgxpool"
	datav1 "github.com/nanzhong/tstr/api/data/v1"
	"github.com/nanzhong/tstr/db"
)

type DataServer struct {
	datav1.UnimplementedDataServiceServer

	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func NewDataServer(pgxPool *pgxpool.Pool) datav1.DataServiceServer {
	return &DataServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
	}
}
