// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package util

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Platform struct {
	Name     string             `yaml:"title"`
	Logo     string             `yaml:"logo"`
	Caption  string             `yaml:"caption"`
	Policies map[string]*policy `yaml:"legal"`
	URL      string             `yaml:"url"`
	Draft    bool               `yaml:"draft"`

	Content  string `yaml:"-"`
	FilePath string `yaml:"-"`
}

func NewPlatformFromFile(fpath string) (*Platform, error) {
	fm, content, err := readMdFile(fpath)
	if err != nil {
		return nil, err
	}

	p := Platform{FilePath: fpath, Content: content}
	return &p, yaml.Unmarshal([]byte(fm), &p)
}

func (p *Platform) Key() string {
	return strings.ToLower(p.Name)
}

// toMdFileContent returns the Platform as markdown file content.
func (p *Platform) toMdFileContent() (string, error) {
	var yamlBlob bytes.Buffer
	enc := yaml.NewEncoder(&yamlBlob)
	enc.SetIndent(2)
	defer enc.Close()

	err := enc.Encode(&p)
	return "---\n" + yamlBlob.String() + "---\n" + p.Content, err
}

func (p *Platform) Commit(to *os.File) error {
	if p.FilePath == "" {
		return fmt.Errorf("%s: empty file path", p.Name)
	}

	if to == nil {
		var err error
		if to, err = os.Create(p.FilePath); err != nil {
			return err
		}
		defer to.Close()
	}

	content, err := p.toMdFileContent()
	if err != nil {
		return err
	}

	_, err = to.Write([]byte(content))
	return err
}

type policy struct {
	Name    string `yaml:"name,omitempty"`
	AsOf    string `yaml:"asof"`
	URL     string `yaml:"url"`
	Sel     string `yaml:"sel,omitempty"`
	Arc     string `yaml:"arc,omitempty"`
	WordCnt int    `yaml:"wc"`
}

func (p *policy) LastCheckedNow() {
	p.AsOf = time.Now().Format("2006-01-02")
}
