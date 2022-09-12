package process

import (
	"errors"
	"os"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func Do(dirRnEnt, dirInEnt, dirOutEnt, dirErrEnt, dirRnCol, dirInCol, dirOutCol, dirErrCol string) error {

	if !fd.DirExists(dirRnEnt) || !fd.DirExists(dirRnCol) {
		err := errors.New("input 'Renamed' Dirs are NOT existing for Processing")
		lk.WarnOnErr("%v", err)
		return err
	}

	mInOut := map[string]string{
		dirInEnt: dirOutEnt,
		dirInCol: dirOutCol,
	}
	mInErr := map[string]string{
		dirInEnt: dirErrEnt,
		dirInCol: dirErrCol,
	}

	gio.MustCreateDirs(dirOutEnt, dirErrEnt, dirOutCol, dirErrCol)

	for I, dir := range []string{dirInEnt, dirInCol} {

		out := mInOut[dir]       // "out" is final output directory for ingestion
		errfolder := mInErr[dir] // "err" is for incorrect format json dump into

		if err := fd.RmFilesIn(out, false, false); err != nil {
			lk.WarnOnErr("%v", err)
			return err
		}

		if err := fd.RmFilesIn(errfolder, false, false); err != nil {
			lk.WarnOnErr("%v")
			return err
		}

		if err := Preproc(dir, out, errfolder); err != nil {
			lk.WarnOnErr("Preproc: %v", err)
			return err
		}

		if err := DumpClassLinkage(out, "class-link.json", "RefName", "ClassLinkage"); err != nil {
			lk.WarnOnErr("DumpClassLinkage: %v", err)
			return err
		}

		if err := DumpPathValue(out, "path_val"); err != nil {
			lk.WarnOnErr("DumpPathValue: %v", err)
			return err
		}

		if I == 0 {
			if err := DumpCollection(out, "collection-entities.json", "RefName", "CollectionEntities"); err != nil {
				lk.WarnOnErr("DumpCollection: %v", err)
				return err
			}
		}
	}

	// ------------------------------------------------------------------------------------- //

	// remove error folder for empty error files, collection error folder can be also deleted here!

	fpaths, _, err := fd.WalkFileDir(dirErrEnt, true)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}
	if len(fpaths) == 0 {
		if err := os.RemoveAll(dirErrEnt); err != nil {
			lk.WarnOnErr("%v", err)
			return err
		}
	}

	// fpaths, _, err = fd.WalkFileDir(dirErrCol, true)
	// if err != nil {
	// 	lk.WarnOnErr("%v", err)
	// 	return err
	// }
	// if len(fpaths) == 0 {
	// 	if err := os.RemoveAll(dirErrCol); err != nil {
	// 		lk.WarnOnErr("%v", err)
	// 		return err
	// 	}
	// }

	return nil
}
