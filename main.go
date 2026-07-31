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

var wg sync.WaitGroup

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

func worker (jobchn chan string , resultchn chan WordCounter) {
		defer wg.Done()
		filename := <- jobchn

		Countedwords := countWords(filename)

		resultchn <- WordCounter{filename: filename , words: Countedwords}
}

func listfiles(dir string) ([]WordCounter){
	jobchn := make(chan string)
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


	for _,filename := range list {
	 	jobchn <- filename	
	}

	wg.Add(2)
	go worker(jobchn,chn)
	go worker(jobchn,chn)

	for i:=0; i< len(list); i++ {
		data:= <- chn
		results = append(results,data)
	}

	wg.Wait()

	fmt.Println("All the files have been read by the 2 workers")

	return results

}

func main(){
	dir := "./notes/"
	results := listfiles(dir)

	for _,result := range results {
		fmt.Println("file: ",result.filename," has: ",result.words," words")
			}
}
