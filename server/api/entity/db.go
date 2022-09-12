package entity

import (
	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
)

func allEntities() ([]*EntityType, error) {

	// inbound db

	mh.UseDbCol(cfg.db, cfg.colText)
	entitiesIn, err := mh.Find[EntityType](nil)
	if err != nil {
		return nil, err
	}

	// existing db

	mh.UseDbCol(cfg.db, "entities")
	entitiesEx, err := mh.Find[EntityType](nil)
	if err != nil {
		return nil, err
	}

	/////////

	return append(entitiesIn, entitiesEx...), nil
}

func allEntityNames() ([]string, error) {
	entities, err := allEntities()
	if err != nil {
		return nil, err
	}
	return FilterMap(entities, nil, func(i int, e *EntityType) string { return e.Entity }), nil
}
