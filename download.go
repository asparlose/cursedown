package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/asparlose/golib/zipfile"
)

func modFileURL(projectID, fileID int) (string, error) {
	modURL := fmt.Sprintf("https://minecraft.curseforge.com/mc-mods/%d", projectID)
	resp, err := http.Get(modURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	return fmt.Sprintf("%s/files/%d/download", resp.Request.URL, fileID), nil
}

func download(projectID, fileID int, modsDir string) error {
	fileURL, err := modFileURL(projectID, fileID)
	if err != nil {
		return err
	}
	fmt.Printf("[%d] %s\n", projectID, fileURL)

	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	filename := resp.Request.URL.String()
	filename2 := filename[strings.LastIndex(filename, "/")+1:]
	filename, err = url.QueryUnescape(filename2)
	if err != nil {
		filename = filename2
	}

	tempfile, err := ioutil.TempFile("", "cursedown")
	if err != nil {
		return err
	}
	defer tempfile.Close()
	defer os.Remove(tempfile.Name())

	tempname := tempfile.Name()

	_, err = io.Copy(tempfile, resp.Body)
	if err != nil {
		return err
	}
	tempfile.Close()

	dir := fmt.Sprintf("%s/mods", modsDir)

	tempfile2, err := os.Open(tempname)
	if err != nil {
		return err
	}

	z, err := zipfile.OpenReader(tempname)
	if err != nil {
		return err
	}
	defer z.Close()

	if len(z.Find(zipfile.DescendantsOf("META-INF")).Slice()) == 0 {
		dir = fmt.Sprintf("%s/resourcepacks", modsDir)
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, tempfile2)
	if err != nil {
		return err
	}

	fmt.Printf("[%d] completed: %s\n", projectID, filename)

	return nil
}
