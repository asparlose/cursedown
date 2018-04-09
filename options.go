package main

var options struct {
	Modpack          string `short:"m" long:"modpack" required:"true" description:"modpack file"`
	OutputDirectory  string `short:"o" long:"out" required:"true" description:"target directory"`
	AllFlag          bool   `short:"a" long:"all" description:"download non-required mods"`
	ParallelDownload int    `short:"p" long:"parallel-download" default:"4" description:"number of parallel download"`
	OverrideOnly     bool   `short:"r" long:"override-only" description:"override only"`
	DownloadOnly     bool   `short:"d" long:"download-only" description:"download only"`
}
