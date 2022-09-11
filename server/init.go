package main

import (
	lk "github.com/digisan/logkit"
	in "github.com/nsip/data-dic-api/server/ingest"
)

func init() {
	lk.FailOnErr("%v", in.Ingest())
}
