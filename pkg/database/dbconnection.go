package database

//go:generate mockgen -source dbconnection.go -destination dbconnection_mock.go -package database

import (
	"database/sql"

	"bulk/logger"

	"bulk/config"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	l     logger.ILogger
	MYSQL config.DatabaseConfig
}

type MYSQL struct {
	Host               string
	User               string
	Password           string
	DBName             string
	MaxPoolSize        int
	MaxIdleConnections int
}

type DBConnection interface {
	DBConnect()
}

func NewDatabaseConnection(logger logger.ILogger, M config.DatabaseConfig) *Connection {
	return &Connection{l: logger, MYSQL: M}
}

func (db *Connection) DBConnect() *sql.DB {
	dbConn, errConn := sql.Open("mysql", db.MYSQL.Username+":"+db.MYSQL.Password+"@tcp("+db.MYSQL.Address+")/"+db.MYSQL.Database+"?parseTime=true")

	if errConn != nil {
		db.l.Fatalf("Error while connecting database. err= %v", errConn.Error())
		return nil
	}

	errPing := dbConn.Ping()
	if errPing != nil {
		db.l.Fatalf("Error while ping database. err= %v", errPing.Error())
		return nil
	}

	dbConn.SetMaxOpenConns(db.MYSQL.MaxOpenConnections)
	dbConn.SetMaxIdleConns(db.MYSQL.MaxIdleConnections)
	return dbConn
}
