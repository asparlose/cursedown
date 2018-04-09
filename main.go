package main

import (
	"archive/zip"
	"fmt"
	"sync"

	"github.com/jessevdk/go-flags"
)

func main() {
	p := flags.NewParser(&options, flags.HelpFlag)
	_, err := p.Parse()

	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return
		}
		panic(err)
	}

	modpack, err := zip.OpenReader(options.Modpack)
	if err != nil {
		panic(err)
	}

	m, err := manifestFromZipFile(&modpack.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s (%s)\n", m.Name, m.Version)

	fmt.Printf("%d mods (%d required, %d not required)\n", len(m.Files), m.RequiredFiles, len(m.Files)-m.RequiredFiles)

	if !options.OverrideOnly {

		count := len(m.Files)
		if !options.AllFlag {
			count = m.RequiredFiles
		}
		fmt.Printf("Downloading %d mods...\n", count)

		wg := new(sync.WaitGroup)
		ch := make(chan manifestFile)

		for i := 0; i < options.ParallelDownload; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, ch <-chan manifestFile) {
				defer wg.Done()
				for file := range ch {
					err := download(file.ProjectID, file.FileID, options.OutputDirectory)
					if err != nil {
						fmt.Printf("[%d] failed: %s\n", file.ProjectID, err)
					}
				}
			}(wg, ch)
		}

		for _, file := range m.Files {
			if file.Required || options.AllFlag {
				ch <- file
			}
		}

		close(ch)
		wg.Wait()
	}

	if !options.DownloadOnly {
		fmt.Println("Override files...")

		err = override(&modpack.Reader, options.OutputDirectory)
		if err != nil {
			panic(err)
		}
	}
}
