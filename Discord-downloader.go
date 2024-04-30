package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/*
Do not run this script with Discord active. Quit the application first, then execute the script again.
*/
func main() {
	body, err := downloadFile("https://discord.com/api/download?platform=linux&format=tar.gz")
	if err != nil {
		log.Fatal(err.Error())
	}
	gzipReader, err := unzip(body)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = untar(gzipReader, "/opt/"); err != nil {
		log.Fatal(err.Error())
	}

	if err := body.Close(); err != nil {
		log.Fatalf("Error trying to close the response body:\n%s.\n", err)
	}
}

/*
Step 1:
func downloadFile returns the *http.Response from the target link.
The response body contains the software package that we want to download.
*/
func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error trying to reach the download source:\n%s.\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received status code %d from the link.\nMaybe try again?\n", resp.StatusCode)
	}

	return resp.Body, nil
}

/*
Step 2:
func unzip returns the decompressed contents of the package using the gzip format.
This allows us to correctly extract the raw file contents and save them into the file system.
*/
func unzip(body io.ReadCloser) (*gzip.Reader, error) {
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("Error trying to unzip the package:\n%s.\n", err)
	}

	if err := gzipReader.Close(); err != nil {
		return nil,
			fmt.Errorf("%s.\nTry running the script again, or maybe test the health of your device?\n", err)
	}

	return gzipReader, nil
}

/*
Step 3:
func untar extracts the decompressed software package and saves it to /opt/Discord/ in the file system.
*/
func untar(gzipReader *gzip.Reader, targetDir string) error {
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("Error trying to read the tar archive:\n%s.\n", err)
		}

		info := header.FileInfo()

		// Create directories or files depending on the type:
		switch header.Typeflag {
		case tar.TypeDir:
			if err := copyFolder(targetDir, header); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := copyFile(tarReader, header, targetDir, info); err != nil {
				return err
			}
		}
	}
}

/*
func copyFolder copies the directory structure of the archive into the file system, to be saved permanently.
*/
func copyFolder(targetDir string, header *tar.Header) error {
	if err := os.MkdirAll(filepath.Join(targetDir, header.Name), 0755); err != nil {
		return fmt.Errorf("Error trying to make a new folder at the location %s%s:\n%s.\n",
			targetDir, header.Name, err,
		)
	}

	return nil
}

/*
func copyFile copies the contents of the archive files into the file system, to be saved permanently.
*/
func copyFile(tarReader io.Reader, header *tar.Header, targetDir string, info fs.FileInfo) error {
	outFile, err := os.OpenFile(filepath.Join(targetDir, header.Name), os.O_CREATE|os.O_RDWR, info.Mode())
	if err != nil {
		return fmt.Errorf(
			"Error trying to create a new file at the location %s%s:\n%s.\nIs Discord still running?\n",
			targetDir, header.Name, err,
		)
	}

	if _, err := io.Copy(outFile, tarReader); err != nil {
		_ = outFile.Close()
		return fmt.Errorf("Error while copying file contents into the new directory:\n%s.\n", err)
	}

	if err = outFile.Close(); err != nil {
		return err
	}

	return nil
}
