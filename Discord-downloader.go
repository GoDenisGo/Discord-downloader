package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	body := downloadFile("https://discord.com/api/download?platform=linux&format=tar.gz")
	gzipReader := unzip(body)
	untar(gzipReader, "/opt/")

	if err := body.Close(); err != nil {
		log.Fatalf("Error trying to close the response body:\n%s.\n", err)
	}
}

/*
Step 1:
func downloadFile returns the *http.Response from the target link.
The response body contains the software package that we want to download.
*/
func downloadFile(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error trying to reach the download source:\n%s.\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received status code %d from the link.\nMaybe try again?\n", resp.StatusCode)
	}

	return resp.Body
}

/*
Step 2:
func unzip returns the decompressed contents of the package using the gzip format.
This allows us to correctly extract the raw file contents and save them into the file system.
*/
func unzip(body io.ReadCloser) *gzip.Reader {
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		log.Fatalf("Error trying to unzip the package:\n%s.\n", err)
	}

	if err := gzipReader.Close(); err != nil {
		log.Fatalf("%s.\nTry running the script again, or maybe test the health of your device?\n", err)
	}

	return gzipReader
}

/*
Step 3:
func untar extracts the decompressed software package and saves it to /opt/Discord/ in the file system.
*/
func untar(gzipReader *gzip.Reader, targetDir string) {
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("Error trying to read the tar archive:\n%s.\n", err)
		}

		info := header.FileInfo()

		// Create directories or files depending on the type
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath.Join(targetDir, header.Name), 0755)
			if err != nil {
				log.Fatalf("Error trying to make a new folder at the location %s/%s:\n%s.\n",
					targetDir, header.Name, err,
				)
			}
		case tar.TypeReg:
			outFile, err := os.Create(filepath.Join(targetDir, header.Name))
			if err != nil {
				log.Fatalf("Error trying to make a new file at the location %s/%s:\n%s.\n",
					targetDir, header.Name, err,
				)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				log.Fatalf("Error while copying file contents into the new directory:\n%s.\n", err)
			}

			if err := os.Chmod(outFile.Name(), info.Mode()); err != nil {
				log.Fatalf("Error while changing file access modes:\n%s.\n", err)
			}

			if err := outFile.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
