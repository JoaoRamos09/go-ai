package database

import (
	"context"
	"database/sql"
	"log"
	"time"
	_"github.com/lib/pq"
)

func NewPostgres(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime string, port string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections to the database
	db.SetMaxOpenConns(maxOpenConns)

	// Set the maximum number of idle connections to the database
	db.SetMaxIdleConns(maxIdleConns)

	// Set the maximum idle time for a connection to the database
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle time for a connection to the database
	db.SetConnMaxIdleTime(duration)

	// Set the context for the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database to check if the connection is successful, with a timeout of ctx
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected successfully to database in port %s", port)

	return db, nil

}
