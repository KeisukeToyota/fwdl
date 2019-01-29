package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
)

func main() {
	wg := &sync.WaitGroup{}

	url := os.Args[1]
	fmt.Println(url)
	res, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start")

	maps := res.Header
	length, err := strconv.Atoi(maps["Content-Length"][0])
	limit := getParseCount(length)
	bar := pb.StartNew(limit)

	fmt.Println("\nParse count : " + strconv.Itoa(limit) + "\n")

	lenSub := length / limit
	body := make([]string, limit+1)
	for i := 0; i < limit; i++ {
		wg.Add(1)

		min := lenSub * i
		max := lenSub * (i + 1)

		if i == limit-1 {
			max += length % limit
		}

		go func(min int, max int, i int) {
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)

			if err != nil {
				panic(err)
			}

			req.Header.Add("Range", "bytes="+strconv.Itoa(min)+"-"+strconv.Itoa(max-1))
			res, err := client.Do(req)

			if err != nil {
				panic(err)
			}

			defer res.Body.Close()
			render, err := ioutil.ReadAll(res.Body)

			if err != nil {
				panic(err)
			}

			body[i] = string(render)
			bar.Increment()
			wg.Done()
		}(min, max, i)
	}

	wg.Wait()
	ioutil.WriteFile(path.Base(url), []byte(string(strings.Join(body, ""))), 0x777)

	fmt.Println("\n\nFinish")
}

func getParseCount(length int) int {
	var count int
	if length > 104857600 {
		count = 100
	} else if length > 1048576 {
		count = 10
	} else {
		count = 2
	}
	return count
}
