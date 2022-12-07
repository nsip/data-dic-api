package db

import (
	"fmt"
	"path/filepath"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

const (
	DATABASE = "MyDictionary"
)

//////////////////////////////////////////////////////

type (
	DbColType string
	DbColName string
)

const (
	// Item
	Existing DbColType = "existing"
	Text     DbColType = "text"
	Html     DbColType = "html"

	// Action
	Submit    DbColType = "submit"
	Approve   DbColType = "approve"
	Subscribe DbColType = "subscribe"
)

//////////////////////////////////////////////////////

type ItemConfig struct {
	DbColExisting DbColName
	DbColText     DbColName
	DbColHtml     DbColName
	DirExisting   string
	DirText       string
	DirHtml       string
}

func (cfg ItemConfig) String() string {
	db := fmt.Sprintf("Item: @database: [%s]; @db-collection(text): [%s]; @db-collection(html): [%s]", DATABASE, cfg.DbColText, cfg.DbColHtml)
	dir := fmt.Sprintf("@directory(text): [%s]; @directory(html): [%s]", cfg.DirText, cfg.DirHtml)
	return db + lk.LF + dir
}

func (cfg ItemConfig) DbColName(col DbColType) DbColName {
	switch col {
	case Existing:
		return cfg.DbColExisting
	case Text:
		return cfg.DbColText
	case Html:
		return cfg.DbColHtml
	default:
		lk.FailOnErr("%v", fmt.Errorf("DbColType [%v] is unknown", col))
	}
	return ""
}

/////////////////////////////////////////////////////////////

type ActionConfig struct {
	DbColSubmit    DbColName
	DbColApprove   DbColName
	DbColSubscribe DbColName
}

func (cfg ActionConfig) String() string {
	db := fmt.Sprintf("Action: @database: [%s]; @db-collection(submit): [%s]; @db-collection(approval): [%s]; @db-collection(subscribe): [%s];",
		DATABASE,
		cfg.DbColSubmit,
		cfg.DbColApprove,
		cfg.DbColSubscribe,
	)
	return db
}

func (cfg ActionConfig) DbColName(col DbColType) DbColName {
	switch col {
	case Submit:
		return cfg.DbColSubmit
	case Approve:
		return cfg.DbColApprove
	case Subscribe:
		return cfg.DbColSubscribe
	default:
		lk.FailOnErr("%v", fmt.Errorf("DbColType [%v] is unknown", col))
	}
	return ""
}

//////////////////////////////////////////////////////

var (
	DirEntity = "./data/inbound/entities"
	CfgEntity = ItemConfig{
		DbColExisting: "entities",
		DbColText:     "entities_text",
		DbColHtml:     "entities_html",
		DirText:       filepath.Join(DirEntity, "text"),
		DirHtml:       filepath.Join(DirEntity, "html"),
		DirExisting:   "./data/renamed",
	}

	DirCollection = "./data/inbound/collections"
	CfgCollection = ItemConfig{
		DbColExisting: "collections",
		DbColText:     "collections_text",
		DbColHtml:     "collections_html",
		DirText:       filepath.Join(DirCollection, "text"),
		DirHtml:       filepath.Join(DirCollection, "html"),
		DirExisting:   "./data/renamed/collections",
	}

	CfgGrp = map[string]ItemConfig{
		"entity":     CfgEntity,
		"collection": CfgCollection,
	}

	// Computed
	ColEntities DbColName = "colentities"
	Class       DbColName = "class"
	PathVal     DbColName = "pathval"

	//////////////////////////////////////////////////

	CfgAction = ActionConfig{
		DbColSubmit:    "act_submit",
		DbColApprove:   "act_approve",
		DbColSubscribe: "act_subscribe",
	}
)

func init() {

	lk.Log("starting...db")

	lk.Log("ingested entities data store under '%v'", DirEntity)

	gio.MustCreateDir(CfgEntity.DirText)
	gio.MustCreateDir(CfgEntity.DirHtml)

	lk.Log("%v", CfgEntity)

	/////////////////////////////////////////////////////////////

	lk.Log("ingested collections data store under '%v'", DirCollection)

	gio.MustCreateDir(CfgCollection.DirText)
	gio.MustCreateDir(CfgCollection.DirHtml)

	lk.Log("%v", CfgCollection)
}
