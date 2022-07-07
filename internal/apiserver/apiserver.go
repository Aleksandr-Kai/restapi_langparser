package apiserver

import (
	"errors"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store/sqlstore"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Start(cfg *config.Config) error {
	var srv *server

	switch cfg.Type {
	case config.MemStore:
		//srv = newServer(memstore.New(), cfg)
	case config.SQLStore:
		db, err := newDBConnection(cfg.DatabaseURL)
		if err != nil {
			return err
		}
		srv = newServer(sqlstore.New(db), cfg)
	default:
		return errors.New("store type incorrect")
	}

	return srv.Start(cfg.BindAddr)
}

func newDBConnection(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		databaseURL = "user=postgres dbname=postgres password=password sslmode=disable"
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.Domain{}, &model.Proxy{}, &model.Request{}, &model.Queue{})
	return db, err
}
