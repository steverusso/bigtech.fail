// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/steverusso/bigtech.fail/util"
)

var (
	rootDir = flag.String("dir", "../website/content/platforms", "The project's website directory.")
	commit  = flag.Bool("w", false, "Write the changes to the platform data files.")
)

func main() {
	flag.Parse()

	platforms, err := readPlatforms(*rootDir, flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range platforms {
		fmt.Printf("scraping %s\n", p.Name)

		if err := util.ScrapePlatform(nil, p); err != nil {
			log.Fatal(err)
		}

		out := os.Stdout
		if *commit {
			out = nil
		}

		if err := p.Commit(out); err != nil {
			log.Fatal(err)
		}
	}
}

func readPlatforms(dir string, names []string) ([]*util.Platform, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var plats []*util.Platform
	for _, f := range files {
		if f.IsDir() && (contains(names, f.Name()) || len(names) == 0) {
			fpath := filepath.Join(dir, f.Name(), "index.md")

			p, err := util.NewPlatformFromFile(fpath)
			if err != nil {
				return nil, err
			}

			plats = append(plats, p)
		}
	}

	return plats, nil
}

func contains(in []string, str string) bool {
	for _, s := range in {
		if str == s {
			return true
		}
	}
	return false
}
