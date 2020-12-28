package file

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirajabbas/utilsgo"
)

var AllowedRatios = []*big.Rat{
	new(big.Rat).SetFrac64(1, 1),  //1:1
	new(big.Rat).SetFrac64(4, 3),  //4:3
	new(big.Rat).SetFrac64(3, 4),  //3:4
	new(big.Rat).SetFrac64(16, 9), //16:9
}
var AllowedFileTypes = []string{"jpg", "jpeg", "png"}

type File struct {
	Path string
	Name string
}

var files []File

func SortFiles(path string) {
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
		utilsgo.SugarLogger.Error(err)
		log.Fatal("Error reading files from directory.")
	}
	for _, f := range files {
		w, h, err := GetImageDimension(f.Path)
		if err != nil {
			utilsgo.SugarLogger.Error("Error reading image dimention: ", f.Path, err)
			continue
		}
		ratio := big.NewRat(int64(w), int64(h))
		isAllowed := false
		for _, r := range AllowedRatios {
			if r.String() == ratio.String() {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			//move file to unsorted directory
			err := os.Rename(f.Path, fmt.Sprintf("%s/unsorted/%s", path, f.Name))
			if err != nil {
				utilsgo.SugarLogger.Error("error moving file to unsorted directory")
				continue
			}
			counter++
		}
	}
	utilsgo.PrintSuccess(fmt.Sprintf("%d files are identified and moved to unsorted directory", counter))

}

func GetImageDimension(filePath string) (w, h int, error error) {
	if reader, err := os.Open(filePath); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			utilsgo.SugarLogger.Error("Error reading image dimensions: ", err)
			error = err
			return
		}
		w = im.Width
		h = im.Height
		return
	} else {
		utilsgo.SugarLogger.Error("unable to open the file:", err)
		error = err
		return
	}
}
