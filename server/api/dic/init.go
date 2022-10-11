package dic

import (
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	in "github.com/nsip/data-dic-api/server/ingest"
)

var (
	chDone      = make(chan bool)               // finish
	tkIngestAll = time.NewTicker(1 * time.Hour) // re ingest all to restructure
)

var (
	mtx = sync.Mutex{}
	//                    dbcol      kind     names
	mListCache = make(map[string]map[string][]string)
)

func init() {

	//
	mListCache["existing"] = make(map[string][]string)
	mListCache["text"] = make(map[string][]string)
	mListCache["html"] = make(map[string][]string)

	//
	go func() {
		for {
			select {
			case <-chDone:
				// SHOULD NOT BE HERE
				return

			case <-tkIngestAll.C:
				// Re ingest all, then update db(entities/collections)
				lk.Log("Channel Enter: tkIngestAll")

				err := in.IngestViaCmd(false)
				lk.WarnOnErr("%v", err)
				if err == nil {
					lk.Log("re-ingested all, and restructured all in existing folder")
				}
			}
		}
	}()
}
