package file

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirajabbas/utilsgo"
)

//variable to hold the filtering rations
var AllowedRatios = []*big.Rat{}

//variable to hold the allowed file extensions
var AllowedFileTypes = []string{"jpg", "jpeg", "png"}

type File struct {
	Path string
	Name string
}

var files []File

func SortFiles(path, destination string) {
	counter := 0
	err := filepath.Walk(path,
		func(fpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ext := strings.Split(info.Name(), ".")[len(strings.Split(info.Name(), "."))-1]
			allowedType := false
			for _, e := range AllowedFileTypes {
				if ext == e {
					allowedType = true
					break
				}
			}
			if allowedType {
				files = append(files, File{Path: fpath, Name: info.Name()})
			}

			return nil
		})
	if err != nil {
		log.Fatal("Error reading files from directory.", err)
	}
	for _, f := range files {
		w, h, err := getImageDimension(f.Path)
		if err != nil {
			log.Println("Error reading image dimention: ", f.Path, err)
			continue
		}
		ratio := big.NewRat(int64(w), int64(h))
		isAllowed := false
		for _, r := range AllowedRatios {
			if r.String() == ratio.String() {
				isAllowed = true
				//move the file to respective folder
				err = f.copyFile(fmt.Sprintf("%s/%s", destination, strings.Replace(r.String(), "/", ":", 1)))
				if err != nil {
					fmt.Println(fmt.Sprintf("Error copying file %s", f.Path))
				}
				break
			}
		}
		if !isAllowed {
			//move file to unsorted directory
			err := f.copyFile(fmt.Sprintf("%s/%s", destination, "unsorted"))
			if err != nil {
				log.Println("error moving file to unsorted directory", err)
				continue
			}
			counter++
		}
	}
	utilsgo.PrintSuccess(fmt.Sprintf("%d files are unidentified and moved to unsorted directory", counter))
	utilsgo.PrintSuccess("done")

}

func getImageDimension(filePath string) (w, h int, error error) {
	if reader, err := os.Open(filePath); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			log.Println("Error reading image dimensions: ", err)
			error = err
			return
		}
		w = im.Width
		h = im.Height
	} else {
		log.Println("unable to open the file:", err)
		error = err
	}
	return
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func (f *File) copyFile(dst string) error {
	in, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(fmt.Sprintf("%s/%s", dst, f.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
