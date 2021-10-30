package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/seehait-c/bulk-rename/sorter"
)

var opts struct {
	StartID         int    `short:"i" long:"start-id" description:"Start ID which will be assigned to the first file" default:"1"`
	Zeros           int    `short:"z" long:"zeros" description:"Number of additional zeros" default:"0"`
	Prefix          string `short:"p" long:"prefix" description:"File name prefix" default:"file"`
	SpaceBeforeID   bool   `short:"s" long:"space-before-id" description:"Add space before ID"`
	IncludeDotFiles bool   `short:"d" long:"include-dot-files" description:"Include dot files"`
	DryRun          bool   `short:"t" long:"dry-run" description:"Dry-run"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}
	log.Printf("Executing script with arguments: start id: %d, prefix: %s\n", opts.StartID, opts.Prefix)

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error while getting the working directory, %s\n", err)
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Printf("Error while listing files from %s\n", err)
	}
	var filteredFiles []fs.FileInfo
	for _, file := range files {
		if opts.IncludeDotFiles || file.Name()[0] != '.' {
			if file.IsDir() {
				continue
			}

			filteredFiles = append(filteredFiles, file)
		}
	}
	fileCounts := len(filteredFiles)
	lastID := opts.StartID + fileCounts
	idLen := opts.Zeros
	for lastID > 0 {
		idLen++
		lastID /= 10
	}
	space := ""
	if opts.SpaceBeforeID {
		space = " "
	}

	sortedFiles := sorter.NatSort(filteredFiles)
	for i, file := range sortedFiles {
		newFile := fmt.Sprintf("%s%s%0*d%s", opts.Prefix, space, idLen, opts.StartID+i, file.Ext)
		log.Printf("%s -> %s", file.NameExt, newFile)
		if !opts.DryRun {
			if err := os.Rename(file.NameExt, newFile); err != nil {
				log.Printf("Error while renaming %s to %s, %s", file.NameExt, newFile, err)
			}
		}
	}
}
