package db

import (
	"fmt"
	"path/filepath"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

type Config struct {
	Db          string
	ColExisting string
	ColText     string
	ColHtml     string
	DirExisting string
	DirText     string
	DirHtml     string
}

func (cfg Config) String() string {
	db := fmt.Sprintf("Entity: @database: [%s]; @collection(text): [%s]; @collection(html): [%s]", cfg.Db, cfg.ColText, cfg.ColHtml)
	dir := fmt.Sprintf("Entity: @directory(text): [%s]; @directory(html): [%s]", cfg.DirText, cfg.DirHtml)
	return db + lk.LF + dir
}

const (
	Database = "dictionaryTest"
)

var (
	DataDirEnt = "./data/inbound/entities"
	CfgEnt     = Config{
		Db:          Database,
		ColExisting: "entities",
		ColText:     "entities_text",
		ColHtml:     "entities_html",
		DirText:     filepath.Join(DataDirEnt, "text"),
		DirHtml:     filepath.Join(DataDirEnt, "html"),
		DirExisting: "./data/renamed",
	}

	DataDirCol = "./data/inbound/collections"
	CfgCol     = Config{
		Db:          Database,
		ColExisting: "collections",
		ColText:     "collections_text",
		ColHtml:     "collections_html",
		DirText:     filepath.Join(DataDirCol, "text"),
		DirHtml:     filepath.Join(DataDirCol, "html"),
		DirExisting: "./data/renamed/collections",
	}

	CfgGrp = map[string]Config{
		"entity":     CfgEnt,
		"collection": CfgCol,
	}
)

func init() {

	lk.Log("ingested entities data store under '%v'", DataDirEnt)

	gio.MustCreateDir(CfgEnt.DirText)
	gio.MustCreateDir(CfgEnt.DirHtml)

	lk.Log("%v", CfgEnt)

	/////////////////////////////////////////////////////////////

	lk.Log("ingested collections data store under '%v'", DataDirCol)

	gio.MustCreateDir(CfgCol.DirText)
	gio.MustCreateDir(CfgCol.DirHtml)

	lk.Log("%v", CfgCol)
}
