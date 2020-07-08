package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func checkExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	fmt.Println(startMessage)
	flag.Usage = usage
	filePtr := flag.String("f", "img.jpg", "")
	outputPtr := flag.String("o", "default", "")
	vidPtr := flag.Bool("vid", false, "")
	collagePtr := flag.Bool("collage", false, "")
	flag.Parse()
	file := *filePtr
	if !checkExists(file) {
		fmt.Printf("File %s doesn't exists", file)
		return
	}
	ext := filepath.Ext(file)
	if ext != ".jpg" {
		fmt.Println("Unsupported file extension")
		return
	}
	bName := filepath.Base(file)
	kind := "jpg"
	if *vidPtr {
		kind = "mp4"
	}
	rName := *outputPtr
	if rName == "default" {
		rName = fmt.Sprintf("%s_portacli.%s", bName[:len(bName)-len(ext)], kind)
	}
	pFile, err := createPortrait(file, *collagePtr, *vidPtr)
	if err != nil {
		fmt.Println("Can't create portrait:", err)
		return
	}
	fmt.Println("OK! Got portrait:", pFile)
	fmt.Println("Downloading...")
	err = downloadPortrait(rName, pFile)
	if err != nil {
		fmt.Println("Can't download portrait:", err)
		return
	}
	fmt.Printf("Done! Portrait is saved as %s.", rName)
}
