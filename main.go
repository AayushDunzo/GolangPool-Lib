package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Result struct {
	work_id int
	Timeout int
}

var links = []string{
	"http://google.com",
	"http://facebook.com",
	"http://stackoverflow.com",
	"http://golang.org",
	"http://amazon.com",
	"http://youtube.com",
	"http://apple.com",
	"http://cloudfare.com",
	"http://mozilla.org",
	"http://wordpress.org",
}

func worker(id int, jobs <-chan int, results chan<- Result) {

	for j := range jobs {

		fmt.Println("worker", id, "started  job", j, "with len of jobs", len(jobs), "remaining")
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(3000*time.Millisecond))
		statusCode := make(chan int)
		timeout_val := 0
		go func() {
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, links[j-1], nil)
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {

				log.Println(err)

			} else {
				statusCode <- res.StatusCode
			}
		}()

		select {
		case <-ctx.Done():
			{
				fmt.Println("timeout dropping the request", j)
				timeout_val = 1

			}
		case sc := <-statusCode:
			{
				fmt.Println("successfully got response: ", sc)
				fmt.Println("worker", id, "finished job", j)
			}

		}

		results <- Result{j, timeout_val}

	}

}

func main() {

	start := time.Now()

	numJobs := 10
	jobs := make(chan int, numJobs)
	results := make(chan Result, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	var count int = 0
	for a := 1; a <= numJobs; a++ {
		res := <-results
		fmt.Println(res.work_id, res.Timeout)
		count++
	}
	fmt.Println("The count is", count)
	fmt.Println(time.Since(start))

}
