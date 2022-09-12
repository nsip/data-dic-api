package collection

import (
	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
)

func allCollections() ([]*CollectionType, error) {

	// inbound db

	mh.UseDbCol(cfg.db, cfg.colText)
	collectionsIn, err := mh.Find[CollectionType](nil)
	if err != nil {
		return nil, err
	}

	// existing db

	mh.UseDbCol(cfg.db, "collections")
	collectionsEx, err := mh.Find[CollectionType](nil)
	if err != nil {
		return nil, err
	}

	/////////

	return append(collectionsIn, collectionsEx...), nil
}

func allCollectionNames() ([]string, error) {
	collections, err := allCollections()
	if err != nil {
		return nil, err
	}
	return FilterMap(collections, nil, func(i int, e *CollectionType) string { return e.Entity }), nil
}
