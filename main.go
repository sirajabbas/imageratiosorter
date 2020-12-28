package main

import (
	"fmt"
	"imagesorter/file"
	"log"
	"os"

	"github.com/sirajabbas/utilsgo"
)

func main() {
	utilsgo.Configure(utilsgo.Config{ErrorFileName: "err.log", LogFileName: "info.log"})
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}
	folder := args[0]

	switch folder {
	case "":
		log.Fatal("invalid folder name. Try help command")
	case "help":
		fmt.Println(`
		imagesorder foldername
		`)
	}
	utilsgo.PrintSuccess("staring application")
	//creating unsorted directory
	err := utilsgo.CreateDirectory(folder + "/unsorted")
	if err != nil {
		utilsgo.SugarLogger.Error(err)
		log.Fatal("Error creating unsorted directory")
	}
	utilsgo.PrintSuccess("created unsorted directory")
	file.SortFiles(folder)
}

func printHelp() {
	fmt.Println(`
		imagesorder foldername
		imagesorter help
	`)
}
