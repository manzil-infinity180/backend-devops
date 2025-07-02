package db

import (
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func StartDb() (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost:5432" // or default to something
	}
	opts = &pg.Options{
		Addr:     host,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	// helps you to connect
	db := pg.Connect(opts)

	// migration
	collection := migrations.NewCollection()
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}
	//start the migrations
	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}
	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return nil, err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}
	return db, err
}
