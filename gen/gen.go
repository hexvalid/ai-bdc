package main

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

//require php & php-cgi

//	php -c php/php.ini -t php/ -S 127.0.0.1:9000
//	curl http://127.0.0.1:9999/index.php
//	cat /tmp/bdc_pipe

//	rm /tmp/bdc_pipe && touch /tmp/bdc_pipe && rm /tmp/bdc_void/*

var folder = "/tmp/bdc_void/"
var okFolder = "ok/"

var count = 250000

var mode = 1

func main() {

	if mode == 1 {
		var c0 = &http.Client{}
		var c1 = &http.Client{}
		var c2 = &http.Client{}
		var c3 = &http.Client{}
		var c4 = &http.Client{}
		var c5 = &http.Client{}
		var c6 = &http.Client{}
		var c7 = &http.Client{}
		var c8 = &http.Client{}
		var c9 = &http.Client{}

		bar := pb.StartNew(count)

		for i := 0; i < (count / 10); i++ {
			var wg sync.WaitGroup
			wg.Add(10)

			go func() {
				getImage(c0, "9000", "a0u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()
			go func() {
				getImage(c1, "9001", "a1u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c2, "9002", "a2u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c3, "9003", "a3u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c4, "9004", "a4u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c5, "9005", "a5u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c6, "9006", "a6u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c7, "9007", "a7u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c8, "9008", "a8u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			go func() {
				getImage(c9, "9009", "a9u5q96vqe9leonrs28fhcpjvc")
				bar.Increment()
				wg.Done()
			}()

			wg.Wait()
		}

		bar.Finish()
	} else if mode == 2 {
		file, err := os.Open("/tmp/bdc_pipe")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		i := 0

		bar := pb.StartNew(count)

		for scanner.Scan() {
			ids := strings.Split(scanner.Text(), ":")

			input, err := ioutil.ReadFile(folder + ids[0] + ".jpg")
			if err != nil {
				panic(err)
				return
			}

			err = ioutil.WriteFile(okFolder+ids[1]+".jpg", input, 0644)
			if err != nil {
				panic(err)
				return
			}

			bar.Increment()
			i++

			if i == count {
				fmt.Println("fitted")
				os.Exit(1)
			} else if i > count {
				fmt.Println("fatted")
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		bar.Finish()

	}
}

func findInString(str, start, end string) string {
	var match []byte
	index := strings.Index(str, start)

	if index == -1 {
		return ""
	}

	index += len(start)

	for {
		char := str[index]

		if strings.HasPrefix(str[index:index+len(match)], end) {
			break
		}

		match = append(match, char)
		index++
	}

	return string(match)
}

func getImage(c *http.Client, port, session string) {

	req, _ := http.NewRequest("GET", "http://127.0.0.1:"+port+"/index.php", nil)

	req.Header.Set("Cookie", "PHPSESSID="+session)

	res, _ := c.Do(req)

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		defer res.Body.Close()

		id := findInString(bodyString, "S:", ":E")
		//imgURL := findInString(bodyString, "script src=\"", "\"")

		req2, _ := http.NewRequest("GET", "http://localhost:"+port+"/botdetect.php?get=image&c=DefaultCaptcha&t="+id, nil)
		req2.Header.Set("Cookie", "PHPSESSID="+session)
		res2, _ := c.Do(req2)

		defer res2.Body.Close()

		//fmt.Println(pipe[1])

		file, err := os.Create(folder + id + ".jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, res2.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}
