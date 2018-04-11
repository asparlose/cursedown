package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/asparlose/golib/zipfile"
)

func override(z *zip.Reader, target string) error {
	modpack := zipfile.NewReader(z)

	cnt := 0

	for f := range modpack.Find(zipfile.And(zipfile.File(), zipfile.DescendantsOf("overrides"))).Iter(context.Background()) {
		t := f.Name[10:]

		os.MkdirAll(fmt.Sprintf("%s/%s", target, filepath.Dir(t)), 0777)
		stream, err := f.Open()
		if err != nil {
			return err
		}

		err = func() error {
			defer stream.Close()
			file, err := os.Create(fmt.Sprintf("%s/%s", target, t))
			if err != nil {
				return err
			}

			io.Copy(file, stream)
			return nil
		}()

		if err != nil {
			return err
		}
		fmt.Println(t)
		cnt++
	}

	fmt.Printf("Overrided %d files\n", cnt)

	return nil
}
