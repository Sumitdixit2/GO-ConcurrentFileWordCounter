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
	filename []string
	words []int
}

func countWords(file string) int {
	content , err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(content))
	return len(words)
}

func listfiles(dir string) ([]string,[]int){
	chn := make(chan int)
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

	var wg sync.WaitGroup

	for _,file := range list {
		wg.Add(1)
		 go func (file string) {
			 defer wg.Done()
			 fmt.Println("Started: ",file)

			 words:= countWords(file)
			 chn <- words
			 data := <- chn
			 words= append(words,data)


			 fmt.Println("Completed: ",file)
		}(file)

		
	}

	wg.Wait()

	return list,words

}

func main(){
	dir := "./notes/"
	list,words := listfiles(dir)

	for i,filename := range list {
		fmt.Println("file: ",filename," has: ",words[i]," words")
			}
}
