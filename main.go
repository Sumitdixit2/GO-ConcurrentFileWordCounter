package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"strings"
)

func countWords(file string) int {
	content , err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(content))
	return len(words)
}

func listfiles(dir string) ([]string,[]int){
	files,err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	var list []string
	var words []int

	for _,entry := range files{
		if !entry.IsDir() {
		list = append(list,filepath.Join(dir,entry.Name()))
	}
	}

	for _,file := range list {
		words = append(words,countWords(file))
	}

	return list,words

}

func main(){
	dir := "./notes/"
	list,words := listfiles(dir)

	for i,filename := range list {
			fmt.Println("file: ",filename," has ",words[i]," words")
			}
}
