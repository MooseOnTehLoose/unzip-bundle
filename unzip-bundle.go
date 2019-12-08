//copied and modified from https://golangcode.com/unzip-files-in-go/
//This will unzip a go bundle into the current directory, or at a location supplied via arg3

package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var file, output string

	fmt.Println()
	//if a zip file and extraction path defined
	if len(os.Args) == 3 {
		file, output = os.Args[1], os.Args[2]
		//if only zip file specified
	} else if len(os.Args) == 2 {
		file, output = os.Args[1], strings.TrimSuffix(os.Args[1], ".zip")
		//if no arguments specified get mad
	} else if len(os.Args) == 1 {
		fmt.Println("No zip archive specified!")
		os.Exit(1)
	}

	err := Unzip(file, output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unzipped " + file + " to: " + output + "successfully")

}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).

func Unzip(src string, dest string) error {

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		cleanPathName := f.Name
		if strings.Contains(f.Name, ":") {
			cleanPathName = strings.Replace(f.Name, ":", "_", -1)
		}
		if strings.Contains(cleanPathName, "*") {
			cleanPathName = strings.Replace(f.Name, "*", "_", -1)
		}
		if strings.Contains(cleanPathName, "|") {
			cleanPathName = strings.Replace(f.Name, "|", "_", -1)
		}

		fpathDest := filepath.Join(dest, cleanPathName)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpathDest, os.ModePerm)
			continue
		}

		// Make folder for files under it
		if err = os.MkdirAll(filepath.Dir(fpathDest), os.ModePerm); err != nil {
			return err
		}

		//write files to folders created
		outFile, err := os.OpenFile(fpathDest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		//check if this file needs to be gunzip-ed
		if strings.Contains(fpathDest, ".gz") {

			if err := gunzip(fpathDest); err != nil {
				return err
			}

			if err := deleteFile(fpathDest); err != nil {
				return err
			}
		}

	}
	return nil
}

//extract .gz files in the bundle as they're written to dest folder
func gunzip(src string) error {

	filename := src

	gzipfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer gzipfile.Close()

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		return err
	}

	defer reader.Close()

	newfilename := strings.TrimSuffix(filename, ".gz")

	writer, err := os.Create(newfilename)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}

func deleteFile(file string) error {
	// delete file
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}
