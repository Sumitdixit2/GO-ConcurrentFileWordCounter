// This is a concurrent file word counter implemented in go that returns the numbers of words in files and returning their word count concurrently.
package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"strings"
	"sync"
) // standard packages

var wg sync.WaitGroup // variable defined to use weight groups for go routine synchronization

type WordCounter struct { // struct for sending and receiving both filename and their word count in channel communication
	filename string
	words int
}

func countWords(file string) int { // function defined to use standard go packages which returns a single files word count
	content , err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(string(content))
	return len(words)
}

func Jobworker (jobchn chan<- string,list []string) { // A seperate Jobworker defined whose only goal is to send the other workers their jobs and to remove the burden from main.
	defer wg.Done() // basically a signal called by our worker as an indicator that says the job has been done.

	for _,filename := range list { // loop through our list of files and sending each filename into the job channel.
		jobchn <- filename
	}
	close(jobchn) // closing the job channel after all the file names have been sent.
}

func worker (jobchn <-chan string , resultchn chan<- WordCounter) { // main worker that finds out the number of words in the file and send the WordCounter struct through the result channel. 
	defer wg.Done()
	for filename := range jobchn{ // looping through our job channel to receive filenames and calling the countWords function to find out the total words and sending them to the result channel.

		fmt.Println("Started reading file: ",filename)

		Countedwords := countWords(filename)

		fmt.Println("file: ",filename," has been successfully read!")

		resultchn <- WordCounter{filename: filename , words: Countedwords} // sending the struct to the result channel.
	}
}

func listfiles(dir string) ([]WordCounter){ //function that returns the filename and word count as a slice.
	jobchn := make(chan string) //creating the job channel.
	chn := make(chan WordCounter) // creating the result channel.
	files,err := os.ReadDir(dir)

	if err != nil { // error handling.
		log.Fatal(err)
	}

	var list []string
	var results []WordCounter

	for _,entry := range files{ // looping through the files slice and only appending files not dir's in the list slice.
		if !entry.IsDir() {
		list = append(list,filepath.Join(dir,entry.Name()))
	}
	}

	wg.Add(2) // defines that we are waiting for 2 instances of our main workers to finish their work.
	go worker(jobchn,chn)
	go worker(jobchn,chn)

	wg.Add(1) // waiting for jobworker to send all the jobs to the job channel.
	go Jobworker(jobchn,list)

	for i:= 0; i< len(list); i++ { // main thread looping through our list's length and appending the result channel data into our results slice
		data:= <- chn
		results = append(results,data)
	}

	wg.Wait() // Waiting for the workers to do their only then we can continue.

	fmt.Println("All the files have been read by the 2 workers")

	return results // returning the results slice.

}

func main(){
	dir := "./notes/"
	results := listfiles(dir) // recieving the results slice in our main function by calling the listfiles function.

	for _,result := range results { // looping the results slice and printing out the final result.
		fmt.Println("file: ",result.filename," has: ",result.words," words")
			}
}
