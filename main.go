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

func countWords(file string) int {
	content , err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(content))
	return len(words)
}

func listfiles(dir string) []string{
	files,err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	var list []string

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

				 countWords(file)

			 fmt.Println("Completed: ",file)
		}(file)
	}

	wg.Wait()

	return list

}

func main(){
	dir := "./notes/"
	list := listfiles(dir)

	for _,filename := range list {
			fmt.Println("file: ",filename)
			}
}
