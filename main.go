package main

import (
	"bufio"
	"fmt"
	"imagesorter/file"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/sirajabbas/utilsgo"
)

func main() {

	/**
	read the input directory
	read output dir and create it if not exist
	read the ratios and process it*/
	var sourceDir, outputDir, ratioParam string

	sourceDir = readFromConsole("source directory")

	switch sourceDir {
	case "":
		log.Fatal("invalid source directory name. Try help command")
	case "help":
		printHelp()
		return
	}

	//validating output directory
	outputDir = readFromConsole("destination directory")
	if outputDir == "" {
		log.Fatal("invalid destination directory. Try help command")
	} else if strings.ToLower(outputDir) == "default" {
		outputDir = fmt.Sprintf("%s/imagesorter", sourceDir)
	}
	//validating ratio param
	ratioParam = readFromConsole("ratio strings")
	if ratioParam == "" {
		log.Fatal("invalid ratio. try help command")
	} else if strings.ToLower(ratioParam) == "default" {
		ratioParam = "1:1 4:3 3:4 16:9"
	}

	//creating destination directory
	err := utilsgo.CreateDirectory(outputDir + "/unsorted")
	//creating unsorted directory
	if err != nil {
		log.Fatal("Error creating destination directory")
	}
	utilsgo.PrintSuccess("created destination directory")

	//stripping last / char
	if string(sourceDir[len(sourceDir)-1]) == "/" {
		sourceDir = string(sourceDir[0 : len(sourceDir)-1])

	}
	//parsing ratios
	ratioStrings := strings.Split(ratioParam, " ")

	if len(ratioStrings) == 0 {
		log.Fatal("error processing the ratio")
	}
	//processing ratios
	for _, rs := range ratioStrings {
		numStrings := strings.Split(rs, ":")
		if len(numStrings) > 0 && len(numStrings) == 2 {
			w, err := strconv.Atoi(numStrings[0])
			h, err := strconv.Atoi(numStrings[1])
			if err != nil {
				log.Fatal("error processing ratio")
			}
			file.AllowedRatios = append(file.AllowedRatios, new(big.Rat).SetFrac64(int64(w), int64(h)))
		}
		//creating destination directory for this ratio
		fmt.Println("creating directory: ", fmt.Sprintf("%s/%s", outputDir, rs))
		err = utilsgo.CreateDirectory(fmt.Sprintf("%s/%s", outputDir, rs))
		if err != nil {
			log.Fatal("Error creating destination directory.", err)
		}
	}

	file.SortFiles(sourceDir, outputDir)
}

func printHelp() {
	fmt.Println(`imagesorter`)
}

func readFromConsole(query string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(query)
	fmt.Print(":")

	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if text != "" {
			return text
		}

	}
}
