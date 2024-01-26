package main

import (
	"fmt"
	"log"
	"runtime"
)

const sum = 2024
const maxIndex = 6
var indexes = []int {438,240,418,276,256,348,258,382,292,310,296,328}

var workers = runtime.NumCPU()

type Result struct {
	indexes []int 
}

type Job struct {
	s []int
	results chan<- Result
}

func (job Job) Do() {
	job.results <- Result {job.s}
}

func main() {
	fmt.Println("Start")
	if maxIndex > len(indexes) {
		log.Fatalf("Max index %d is bigger than the array length", maxIndex)
	}
	lenOfResults := factorial(len(indexes)) / ( factorial(maxIndex) * factorial (len(indexes) - maxIndex ))
	jobs := make(chan Job, workers)
	results := make (chan Result, lenOfResults) 
	done := make(chan struct{}, workers)
	
	go addNumber(jobs, results)
	for i := 0; i < workers; i++ {
		go doJobs(done, jobs)
	}
	go awaitCompletion(done, results)
	processResults(results)
}

func processResults(results <-chan Result) {
	for r := range results {
		s := 0
		for _, v := range r.indexes {
			if s += v; s > sum {
				break;
			}
		}
		if s == sum {
			fmt.Printf("Result = %d\n", r.indexes)
		}
    }
}

func doJobs(done chan<- struct{}, jobs <-chan Job) {
    for job := range jobs {
        job.Do()
    }
    done <- struct{}{}
}

func awaitCompletion(done <-chan struct{}, results chan Result) {
    for i := 0; i < workers; i++ {
        <-done
    }
    close(results)
}

func addNumber(jobs chan<- Job, results chan<- Result) {
	s := make([]int, 0, maxIndex)
	for i := 0; i < len(indexes); i++ {
		addNumbers(jobs, s, i, results)
	}
	close(jobs)
}

func addNumbers(jobs chan<- Job, subset []int, t int, results chan<- Result) {
	s := make([]int, 0, maxIndex)
	s = append(subset, indexes[t])				
	if len(s) == maxIndex {
		ss := make([]int, 0, maxIndex)
		ss = append(ss, s...)
		jobs <- Job{ss, results}		
	} else {
		for i := t + 1; i < len(indexes); i++ {
			if len(indexes) - i >= maxIndex - len(s){
				addNumbers(jobs, s, i, results)
			} 		
		}
	}
}

func factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}