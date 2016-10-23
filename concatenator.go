package concatenator

import (
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
)

func Concatenator(urls ...string) (megabody string, err error) {
	for _, url := range urls {
		var body string
		body, err = get(url)
		if err != nil {
			return
		} else {
			megabody = megabody + strings.Trim(body, "\n")
		}
	}
	return
}

func get(url string)(body string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode==200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			body = string(bodyBytes)
		}
	} else {
		err = errors.New("Non-200 response attempting to fetch " + url + " : " + resp.Status)
	}
	return
}
