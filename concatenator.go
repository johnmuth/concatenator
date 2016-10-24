package concatenator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func Concatenator(urls ...string) (megabody string, err error) {
	bodyChannel := make(chan string)
	errorChannel := make(chan error)
	doneChannel := make(chan bool)
	go multiGet(urls, bodyChannel, errorChannel, doneChannel)
	for {
		select {
		case body := <-bodyChannel:
			body = strings.Trim(body, "\n")
			megabody = megabody + body
			log.Printf("Added %s megabody=%s", body, megabody)
		case err = <-errorChannel:
			fmt.Errorf("Error: %v", err)
		case <-doneChannel:
			log.Printf("GOT A DONE! megabody=%s", megabody)
			return megabody, err
		}
	}
}

func multiGet(urls []string, bodyChannel chan string, errorChannel chan error, doneChannel chan bool) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			body, err := get(url)
			if err != nil {
				errorChannel <- err
			} else {
				bodyChannel <- body
			}
		}(url)
	}
	wg.Wait()
	log.Println("END OF WAIT")
	doneChannel <- true
}

func get(url string) (body string, err error) {
	var resp *http.Response
	log.Printf("GET %s", url)
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		log.Printf("GOT %s", url)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			body = string(bodyBytes)
		} else {
			log.Printf("ERROR %s : %v", url, err)
			return "", err
		}
	} else {
		err = errors.New("Non-200 response attempting to fetch " + url + " : " + resp.Status)
	}
	return
}
