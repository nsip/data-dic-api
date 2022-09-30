package db

import (
	"fmt"
	"path/filepath"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

const (
	DATABASE = "dictionaryTest"
)

//////////////////////////////////////////////////////

type DbColType string

const (
	// Item
	Existing DbColType = "existing"
	Text     DbColType = "text"
	Html     DbColType = "html"

	// Action
	Submit  DbColType = "submit"
	Approve DbColType = "approve"
)

//////////////////////////////////////////////////////

type ItemConfig struct {
	DbColExisting string
	DbColText     string
	DbColHtml     string
	DirExisting   string
	DirText       string
	DirHtml       string
}

func (cfg ItemConfig) String() string {
	db := fmt.Sprintf("Item: @database: [%s]; @db-collection(text): [%s]; @db-collection(html): [%s]", DATABASE, cfg.DbColText, cfg.DbColHtml)
	dir := fmt.Sprintf("@directory(text): [%s]; @directory(html): [%s]", cfg.DirText, cfg.DirHtml)
	return db + lk.LF + dir
}

func (cfg ItemConfig) DbColVal(col DbColType) string {
	switch col {
	case Existing:
		return cfg.DbColExisting
	case Text:
		return cfg.DbColText
	case Html:
		return cfg.DbColHtml
	default:
		return ""
	}
}

/////////////////////////////////////////////////////////////

type ActionConfig struct {
	DbColSubmit  string
	DbColApprove string
}

func (cfg ActionConfig) String() string {
	db := fmt.Sprintf("Action: @database: [%s]; @db-collection(submit): [%s]; @db-collection(approval): [%s]", DATABASE, cfg.DbColSubmit, cfg.DbColApprove)
	return db
}

func (cfg ActionConfig) DbColVal(col DbColType) string {
	switch col {
	case Submit:
		return cfg.DbColSubmit
	case Approve:
		return cfg.DbColApprove
	default:
		return ""
	}
}

//////////////////////////////////////////////////////

var (
	DataDirEntity = "./data/inbound/entities"
	CfgEntity     = ItemConfig{
		DbColExisting: "entities",
		DbColText:     "entities_text",
		DbColHtml:     "entities_html",
		DirText:       filepath.Join(DataDirEntity, "text"),
		DirHtml:       filepath.Join(DataDirEntity, "html"),
		DirExisting:   "./data/renamed",
	}

	DataDirCollection = "./data/inbound/collections"
	CfgCollection     = ItemConfig{
		DbColExisting: "collections",
		DbColText:     "collections_text",
		DbColHtml:     "collections_html",
		DirText:       filepath.Join(DataDirCollection, "text"),
		DirHtml:       filepath.Join(DataDirCollection, "html"),
		DirExisting:   "./data/renamed/collections",
	}

	CfgGrp = map[string]ItemConfig{
		"entity":     CfgEntity,
		"collection": CfgCollection,
	}

	//////////////////////////////////////////////////

	CfgAction = ActionConfig{
		DbColSubmit:  "act_submit",
		DbColApprove: "act_approve",
	}
)

func init() {

	lk.Log("ingested entities data store under '%v'", DataDirEntity)

	gio.MustCreateDir(CfgEntity.DirText)
	gio.MustCreateDir(CfgEntity.DirHtml)

	lk.Log("%v", CfgEntity)

	/////////////////////////////////////////////////////////////

	lk.Log("ingested collections data store under '%v'", DataDirCollection)

	gio.MustCreateDir(CfgCollection.DirText)
	gio.MustCreateDir(CfgCollection.DirHtml)

	lk.Log("%v", CfgCollection)
}
