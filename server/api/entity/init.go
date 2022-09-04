package entity

import (
	"fmt"

	mh "github.com/digisan/db-helper/mongo"
	lk "github.com/digisan/logkit"
)

type DbConfig struct {
	Database   string `form:"database" json:"database"` // ``
	Collection string `form:"collection" json:"collection"`
}

func (cfg DbConfig) String() string {
	return fmt.Sprintf("Entity: @database: [%s]; @collection: [%s]", cfg.Database, cfg.Collection)
}

var (
	cfg = DbConfig{
		Database:   "dictionary",
		Collection: "entity",
	}
)

func init() {
	lk.Log("%v", cfg)
	mh.UseDbCol(cfg.Database, cfg.Collection)
}
