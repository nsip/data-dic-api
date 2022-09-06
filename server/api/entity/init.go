package entity

import (
	"fmt"

	mh "github.com/digisan/db-helper/mongo"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

const (
	dataFolder = "./data/"
)

type DbConfig struct {
	Database   string `form:"database" json:"database"`
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
	gio.MustCreateDir(dataFolder)
	lk.Log("ingested data store at %v", dataFolder)

	mh.UseDbCol(cfg.Database, cfg.Collection)
	lk.Log("%v", cfg)
}
