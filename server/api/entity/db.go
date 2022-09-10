package entity

import (
	mh "github.com/digisan/db-helper/mongo"
)

func allEntities() ([]*EntityType, error) {

	mh.UseDbCol(cfg.db, cfg.colText)
	entitiesIn, err := mh.Find[EntityType](nil)
	if err != nil {
		return nil, err
	}

	/////////

	mh.UseDbCol(cfg.db, "entities")
	entitiesEx, err := mh.Find[EntityType](nil)
	if err != nil {
		return nil, err
	}

	/////////

	return append(entitiesIn, entitiesEx...), nil
}
