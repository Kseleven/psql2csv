package csv

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const ConnStr string = "user=%s password=%s host=%s port=%d database=%s sslmode=disable "

func NewDB(conf *Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf(ConnStr, conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName))
	if err != nil {
		return nil, err
	}

	return conn, conn.Ping(context.Background())
}
