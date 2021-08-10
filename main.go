package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Result type definition
type Result string

//Search type definition
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

var (
	//Web1 fake search
	Web1 = fakeSearch("web1")
	//Web2 fake search
	Web2 = fakeSearch("web2")
	//Image1 fake search
	Image1 = fakeSearch("image1")
	//Image2 fake search
	Image2 = fakeSearch("image2")
	//Video1 fake search
	Video1 = fakeSearch("video1")
	//Video2 fake search
	Video2 = fakeSearch("video2")
)

func main() {
	//Seed value to generate rand values.
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	end := time.Since(start)
	fmt.Println(results)
	fmt.Printf("%v ms\n", end)
}

//Actual implementation to search
func Google(query string) (results []Result) {
	c := make(chan Result)
	//Go routine for each query to look on two Search engines.
	//Te first one to finish is the one that is appended
	go func() { c <- Replicate(query, Web1, Web2) }()
	go func() { c <- Replicate(query, Image1, Image2) }()
	go func() { c <- Replicate(query, Video1, Video2) }()
	//If timeout occured then return the results we have in the time set.
	timeout := time.After(time.Duration(time.Millisecond * 60))
	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		case <-timeout:
			return
		}

	}
	return results
}

//Replicate will create a go routine for each search type received and send the query.
func Replicate(query string, search ...Search) Result {
	c := make(chan Result)
	for i := range search {
		s := search[i]
		go func() {
			c <- s(query)
		}()
	}
	return <-c
}
