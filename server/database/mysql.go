package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	"emergency-his/server/config"
	"github.com/go-sql-driver/mysql"
)

func InitMySQL(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := mysql.NewConfig()
	dsn.User = cfg.Username
	dsn.Passwd = cfg.Password
	dsn.Net = "tcp"
	dsn.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	dsn.DBName = cfg.DBName
	dsn.ParseTime = true
	dsn.Loc = time.Local
	dsn.Params = map[string]string{"charset": "utf8mb4"}
	dsn.Timeout = 5 * time.Second

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return db, nil
}
