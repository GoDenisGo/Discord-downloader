package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// URL of the tar.gz file
	url := "https://discord.com/api/download?platform=linux&format=tar.gz"

	// Directory to extract files to
	targetDir := "/opt/"

	// Step 1: Download the file
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error trying to reach the link: %s.\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received status code %d from the link.\nMaybe try again?", resp.StatusCode)
	}

	// Step 2: Decompress gzip
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalf("Error trying to unzip the package: %s", err)
	}
	defer gzipReader.Close()

	// Step 3: Un-tar the contents
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("Error trying to read the tar archive: %s", err)
		}

		// Create directories or files depending on the type
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetDir+"/"+header.Name, 0755)
			if err != nil {
				log.Fatalf("Error trying to make a new directory at the location %s/%s: %s",
					targetDir, header.Name, err,
				)
			}
		case tar.TypeReg:
			outFile, err := os.Create(targetDir + "/" + header.Name)
			if err != nil {
				log.Fatalf("Error trying to make a new file at the location %s/%s: %s",
					targetDir, header.Name, err,
				)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				log.Fatalf("Error while copying file contents into the new directory: %s", err)
			}
			outFile.Close()
		}
	}
}
