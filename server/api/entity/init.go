package entity

import (
	"fmt"
	"path/filepath"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

type Config struct {
	db      string
	colText string
	colHtml string
	dirText string
	dirHtml string
}

func (cfg Config) String() string {
	dbstr := fmt.Sprintf("Entity: @database: [%s]; @text collection: [%s]; @html collection: [%s]", cfg.db, cfg.colText, cfg.colHtml)
	dirstr := fmt.Sprintf("Entity: @text directory: [%s]; @html directory: [%s]", cfg.dirText, cfg.dirHtml)
	return dbstr + "\n" + dirstr + "\n"
}

const (
	DataDir = "./data/entities"
)

var (
	cfg = Config{
		db:      "dictionaryTest",
		colText: "entities_text",
		colHtml: "entities_html",
		dirText: filepath.Join(DataDir, "text"),
		dirHtml: filepath.Join(DataDir, "html"),
	}
)

func init() {

	lk.Log("ingested entities data store under '%v'", DataDir)

	gio.MustCreateDir(cfg.dirText)
	gio.MustCreateDir(cfg.dirHtml)

	lk.Log("%v", cfg)
}
