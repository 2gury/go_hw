package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

func ExecutePipeline(jobs ...job) {
	waitGroup := &sync.WaitGroup{}
	lenJobs := len(jobs)
	chans := make([]chan interface{}, lenJobs+1)
	for i := 0; i < lenJobs+1; i++ {
		chans[i] = make(chan interface{})
	}

	waitGroup.Add(lenJobs)
	for i := 0; i < lenJobs; i++ {
		go func(jb job, wg *sync.WaitGroup, inCh, outCh chan interface{}) {
			jb(inCh, outCh)
			wg.Done()
			close(outCh)
		}(jobs[i], waitGroup, chans[i], chans[i+1])
	}
	waitGroup.Wait()
	close(chans[0])
}

func SingleHash(in, out chan interface{}) {
	globalWaitGroup := &sync.WaitGroup{}
	for data := range in {
		globalWaitGroup.Add(1)
		localWaitGroup := &sync.WaitGroup{}
		resParts := make([]string, 2)
		stringData := fmt.Sprintf("%v", data)
		md5Data := DataSignerMd5(stringData)
		inputData := []string{stringData, md5Data}

		go func(input []string, results []string, wg *sync.WaitGroup) {
			lenResParts := len(results)
			wg.Add(lenResParts)
			for i := 0; i < lenResParts; i++ {
				go func(val string, res *string) {
					*res = DataSignerCrc32(val)
					wg.Done()
				}(input[i], &results[i])
			}
			wg.Wait()
			out <- fmt.Sprintf("%s~%s", results[0], results[1])
			globalWaitGroup.Done()
		}(inputData, resParts, localWaitGroup)
	}
	globalWaitGroup.Wait()
}

func MultiHash(in, out chan interface{}) {
	globalWaitGroup := &sync.WaitGroup{}
	for data := range in {
		localWaitGroup := &sync.WaitGroup{}
		resParts := make([]string, 6)
		lenResParts := len(resParts)
		stringData := fmt.Sprintf("%v", data)
		localWaitGroup.Add(lenResParts)
		globalWaitGroup.Add(1)

		for i := 0; i < lenResParts; i++ {
			go func(val string, res *string, wg *sync.WaitGroup) {
				*res = DataSignerCrc32(val)
				wg.Done()
			}(fmt.Sprintf(strconv.Itoa(i)+stringData), &resParts[i], localWaitGroup)
		}

		go func(results []string, wg *sync.WaitGroup) {
			wg.Wait()
			result := ""
			for _, value := range results {
				result += value
			}
			out <- result
			globalWaitGroup.Done()
		}(resParts, localWaitGroup)
	}
	globalWaitGroup.Wait()
}

func CombineResults(in, out chan interface{}) {
	result := []string{}
	for data := range in {
		stringData := fmt.Sprintf("%v", data)
		result = append(result, stringData)
	}
	sort.Strings(result)

	var combineResult string
	for id, value := range result {
		if id != len(result)-1 {
			combineResult += fmt.Sprintf("%s_", value)
		} else {
			combineResult += value
		}
	}
	out <- combineResult
}

func main() {
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			out <- 5
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			for val := range in {
				log.Println(val)
			}
		}),
	}

	ExecutePipeline(freeFlowJobs...)
	time.Sleep(time.Second)

	log.Println("OK")
}
