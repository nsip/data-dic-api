package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	proc "github.com/nsip/data-dic-api/server/process"
	"github.com/tidwall/gjson"
)

func init() {
	lk.Log2F(true, "./process.log")
	lk.WarnDetail(false)
}

func main() {

	var (
		wholePtr = flag.Bool("whole", false, "true: for whole process, including 'rename'; otherwise only 'process'")

		root = "./data/"

		// rename
		dirOriEntPtr = flag.String("oed", filepath.Join(root, "original"), "original entities json data directory")
		dirRnEntPtr  = flag.String("red", filepath.Join(root, "renamed"), "renamed entities json data directory")

		dirOriColPtr = flag.String("ocd", filepath.Join(root, "original/collections"), "original collections json data directory")
		dirRnColPtr  = flag.String("rcd", filepath.Join(root, "renamed/collections"), "renamed collections json data directory")

		// process
		dirInEntPtr  = flag.String("ie", filepath.Join(root, "renamed"), "input entities data directory")
		dirOutEntPtr = flag.String("oe", filepath.Join(root, "out"), "output entities data directory")
		dirErrEntPtr = flag.String("ee", filepath.Join(root, "err"), "error entities data directory")

		dirInColPtr  = flag.String("ic", filepath.Join(root, "renamed/collections"), "input collections data directory")
		dirOutColPtr = flag.String("oc", filepath.Join(root, "out/collections"), "output collections data directory")
		dirErrColPtr = flag.String("ec", filepath.Join(root, "err/collections"), "error collections data directory")
	)

	flag.Parse()

	if *wholePtr {

		//////////////////////////////////////////////////////////////

		dirOriEnt, dirRnEnt := *dirOriEntPtr, *dirRnEntPtr

		gio.MustCreateDir(dirRnEnt)

		// clear destination dir for putting renamed file
		lk.FailOnErr("%v", fd.RmFilesIn(dirRnEnt, false, true, "json"))

		// make sure each file's name is its entity value
		proc.FixFileName(dirOriEnt, dirRnEnt)

		//////////////////////////////////////////////////////////////

		dirOriCol, dirRnCol := *dirOriColPtr, *dirRnColPtr

		gio.MustCreateDir(dirRnCol)

		// clear destination dir for putting renamed file
		lk.FailOnErr("%v", fd.RmFilesIn(dirRnCol, false, true, "json"))

		// make sure each file's name is its entity value
		proc.FixFileName(dirOriCol, dirRnCol)

		/////////////////////////////

		mChk := map[string][]string{
			dirRnEnt: {"Element", "Object", "Abstract Element"},
			dirRnCol: {"Collection"},
		}

		for _, dir := range []string{dirRnEnt, dirRnCol} {
			fs, err := os.ReadDir(dir)
			lk.FailOnErr("%v", err)
			for _, f := range fs {
				if fname := f.Name(); strings.HasSuffix(fname, ".json") {
					fpath := filepath.Join(dir, fname)
					data, err := os.ReadFile(fpath)
					lk.FailOnErr("%v", err)
					lk.WarnOnErrWhen(NotIn(gjson.Get(string(data), "Metadata.Type").String(), mChk[dir]...), "%v@%s", errors.New("ERROR TYPE"), fpath)
				}
			}
		}

	}

	// ------------------------------------------------------------------------------------- //

	if !fd.DirExists(*dirRnEntPtr) || !fd.DirExists(*dirRnColPtr) {
		lk.FailOnErr("%v", errors.New("input 'Renamed' Dirs are NOT existing for Processing"))
	}

	mInOut := map[string]string{
		*dirInEntPtr: *dirOutEntPtr,
		*dirInColPtr: *dirOutColPtr,
	}
	mInErr := map[string]string{
		*dirInEntPtr: *dirErrEntPtr,
		*dirInColPtr: *dirErrColPtr,
	}

	gio.MustCreateDirs(*dirOutEntPtr, *dirOutColPtr, *dirErrEntPtr, *dirErrColPtr)

	for I, dir := range []string{*dirInEntPtr, *dirInColPtr} {

		out := mInOut[dir]       // "out" is final output directory for ingestion
		errfolder := mInErr[dir] // "err" is for incorrect format json dump into

		lk.FailOnErr("%v", fd.RmFilesIn(out, false, false))
		lk.FailOnErr("%v", fd.RmFilesIn(errfolder, false, false))

		proc.Preproc(dir, out, errfolder)

		proc.DumpClassLinkage(out, "class-link.json", "ver", "0.0.1")

		proc.DumpPathValue(out, "path_val")

		if I == 0 {
			proc.DumpCollection(out, "collection-entities.json", "ver", "0.0.1")
		}
	}

	// ------------------------------------------------------------------------------------- //

	// remove error folder for empty error files
	errdir := filepath.Join(root, "err")
	fpaths, _, err := fd.WalkFileDir(errdir, true)
	lk.FailOnErr("%v", err)
	if len(fpaths) == 0 {
		lk.FailOnErr("%v", os.RemoveAll(errdir))
	}
}
