package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"strings"
	"sync"
)

type WordCounter struct {
	filename string
	words int
}

func countWords(file string) int {
	content , err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(content))
	return len(words)
}

func listfiles(dir string) ([]WordCounter){
	chn := make(chan WordCounter)
	files,err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	var list []string
	var results []WordCounter

	for _,entry := range files{
		if !entry.IsDir() {
		list = append(list,filepath.Join(dir,entry.Name()))
	}
	}

	var wg sync.WaitGroup

	for _,file := range list {
		wg.Add(1)
		 go func (file string) {
			 defer wg.Done()
			 fmt.Println("Started: ",file)

			 Countedwords:= countWords(file)
			 chn <- WordCounter{filename: file,words: Countedwords}
			 
			 
			 fmt.Println("Completed: ",file)
			 
		}(file)

	}

	for i:=0; i< len(list); i++ {
		data:= <- chn
		results = append(results,data)
	}

	wg.Wait()

	return results

}

func main(){
	dir := "./notes/"
	results := listfiles(dir)

	for _,result := range results {
		fmt.Println("file: ",result.filename," has: ",result.words," words")
			}
}
