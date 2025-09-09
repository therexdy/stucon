package internal

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	RedisDB *redis.Client
	PSQLDB *sql.DB
}

func InitConn() (s *Server, err error){
	s = &Server{}
	s.PSQLDB, err = sql.Open("pgx", "postgres://appuser:GTAC@localhost:5432/stucon?sslmode=disable")
	if err != nil {
		return nil, err
	}
	s.RedisDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return s, nil
}

func (s *Server) CloseConn() (){
	s.PSQLDB.Close()
	s.RedisDB.Close()
}
