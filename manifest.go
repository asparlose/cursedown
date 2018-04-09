package main

import (
	"archive/zip"
	"encoding/json"
	"io"

	"github.com/asparlose/golib/zipfile"
	"github.com/pkg/errors"
)

type manifestRoot struct {
	Name          string         `json:"name"`
	Version       string         `json:"version"`
	Files         []manifestFile `json:"files"`
	RequiredFiles int
}

type manifestFile struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}

func manifestFromZipFile(z *zip.Reader) (manifestRoot, error) {
	zf := zipfile.NewReader(z).Find(zipfile.FullName("manifest.json")).Slice()

	if len(zf) == 0 {
		return manifestRoot{}, errors.New("manifest.json not found")
	}

	mf := zf[0]

	rc, err := mf.Open()
	if err != nil {
		return manifestRoot{}, err
	}
	defer rc.Close()

	buf := make([]byte, mf.UncompressedSize)
	_, err = io.ReadFull(rc, buf)
	if err != nil {
		return manifestRoot{}, err
	}

	manifest := manifestRoot{}

	err = json.Unmarshal(buf, &manifest)

	for _, f := range manifest.Files {
		if f.Required {
			manifest.RequiredFiles++
		}
	}

	return manifest, err
}
