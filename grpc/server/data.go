package server

type DataServer struct {
	datav1.UnimplementedAdminServiceServer

	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func NewDataServer(pgxPool *pgxpool.Pool) datav1.DataServiceServer {
	return &DataServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
	}
}


