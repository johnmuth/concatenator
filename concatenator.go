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
	go channelGet(urls, bodyChannel, errorChannel, doneChannel)
	var mutex = &sync.Mutex{}
	for {
		select {
		case body := <-bodyChannel:
			log.Printf("Locking to add %s", body)
			mutex.Lock()
			megabody = megabody + strings.Trim(body, "\n")
			mutex.Unlock()
			log.Printf("Unlocked after adding %s ... megabody=%s", body, megabody)
		case err = <-errorChannel:
			fmt.Errorf("Error: %v", err)
		case <-doneChannel:
			log.Println("GOT A DONE")

			log.Printf("megabody=%s", megabody)
			return megabody, err
		}
	}
}

func channelGet(urls []string, bodyChannel chan string, errorChannel chan error, doneChannel chan bool) {
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
