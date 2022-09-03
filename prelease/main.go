package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	gio "github.com/digisan/gotk/io"
	nt "github.com/digisan/gotk/net-tool"
	"github.com/digisan/gotk/strs"
)

func LiteralLocIP2PubIP(oldport, newport int, filepaths ...string) error {
	for _, fpath := range filepaths {
		if err := nt.ChangeLocalUrlPort(fpath, oldport, newport, false, true); err != nil {
			return err
		}
		if err := nt.LocIP2PubIP(fpath, false, true); err != nil {
			return err
		}
	}
	return nil
}

func CommentaryReplace(fpath string) error {
	m := map[string]string{
		"Updated@": "Updated@ " + time.Now().Format(time.RFC3339),
	}
	_, err := gio.FileLineScan(fpath, func(line string) (bool, string) {
		if strings.HasPrefix(line, "//") {
			for k, v := range m {
				if strings.Contains(line, k) {
					switch k {
					case "Updated@":
						line = strs.TrimTailFromLast(line, k) + k
						return true, strings.ReplaceAll(line, k, v)
					}
				}
			}
		}
		return true, line
	}, fpath)
	return err
}

func main() {
	fmt.Println("Usage: \n\t-s: go source file [../server/main.go]\n\t-c: commentary symbol replacement [true]\n\t-p: local ip to public ip [false]")
	var (
		s = flag.String("s", "../server/main.go", "go source file")
		c = flag.Bool("c", true, "commentary symbol replacement")
		p = flag.Bool("p", false, "local ip to public ip")
	)
	flag.Parse()

	if *c {
		CommentaryReplace(*s)
	}
	if *p {
		LiteralLocIP2PubIP(1323, 1323, *s)
	}
}
