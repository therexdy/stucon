package internal

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Server struct {
	RedisDB *redis.Client
	PSQLDB *sql.DB
	Minio *minio.Client
	Ctx context.Context
}

func InitConn() (s *Server, err error){
	s = &Server{}

	s.PSQLDB, err = sql.Open("postgres", "postgres://appuser:GTAC@stucon_postgres:5432/stucon?sslmode=disable")
	if err != nil {
		return nil, err
	}

	s.RedisDB = redis.NewClient(&redis.Options{
		Addr: "stucon_redis:6379",
	})

	s.Minio , err = minio.New("stucon_minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("appuser", "GTAC@gtac", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	
	s.Ctx = context.Background()

	return s, nil
}

func (s *Server) CloseConn() (){
	s.PSQLDB.Close()
	s.RedisDB.Close()
}
