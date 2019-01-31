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
	defer fmt.Println("\n\nFinish")

	header := res.Header
	length, err := strconv.Atoi(header["Content-Length"][0])

	if err != nil {
		panic(err)
	}

	count := getParseCount(length)

	fmt.Println("\nParse count : " + strconv.Itoa(count) + "\n")

	lenSub := length / count
	body := make([]string, count+1)
	for i := 0; i < count; i++ {
		wg.Add(1)

		start := lenSub * i
		end := lenSub * (i + 1)

		if i == count-1 {
			end += length % count
		}

		go func(start int, end int, i int) {
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)

			if err != nil {
				panic(err)
			}

			req.Header.Add("Range", "bytes="+strconv.Itoa(start)+"-"+strconv.Itoa(end-1))
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
			wg.Done()
		}(start, end, i)
	}

	wg.Wait()
	ioutil.WriteFile(path.Base(url), []byte(string(strings.Join(body, ""))), 0644)
}

func getParseCount(length int) int {
	if length > 104857600 {
		return 100
	} else if length > 1048576 {
		return 50
	} else {
		return 3
	}
}
