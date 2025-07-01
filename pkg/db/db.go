package db

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"log"
	"os"
)

func StartDb() (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)
	if os.Getenv("ENV") == "PROD" {
		_, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}
	} else {
		opts = &pg.Options{
			Addr:     "db:5432",
			User:     "postgres",
			Password: "postgres",
		}
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
