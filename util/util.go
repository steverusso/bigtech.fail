// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package util

import (
	"bufio"
	"os"
)

// readMdFile reads the given markdown file and returns the frontmatter and
// markdown content.
func readMdFile(fpath string) (fm, md string, _ error) {
	file, err := os.Open(fpath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	inFrontMatter := false

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		ln := sc.Text()
		if ln == "+++" || ln == "---" {
			inFrontMatter = !inFrontMatter
			continue
		}

		if inFrontMatter {
			fm += ln + "\n"
		} else {
			md += ln + "\n"
		}
	}

	return fm, md, sc.Err()
}
