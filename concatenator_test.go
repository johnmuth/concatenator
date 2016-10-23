package concatenator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"time"
)

func TestGetOne(t *testing.T) {
	expected := makeExpectedResponse(1)
	ts := makeTestServer(1)
	defer ts.Close()
	actual, err := get(makeTestUrl(ts.URL, 1))
	if err != nil {
		t.Error("Got unexpected error", err)
	}
	actual = strings.Trim(actual,"\n")
	if actual != expected {
		t.Errorf("'%s' was not equal to '%s'", actual, expected)
	}
}

func TestConcatenator(t *testing.T) {
	expected := makeExpectedResponse(2)
	ts := makeTestServer(2)
	testUrls := makeTestUrls(ts.URL, 2)
	defer ts.Close()
	actual, err := Concatenator(testUrls...)
	if err != nil {
		t.Error("Got unexpected error", err)
	}
	actual = strings.Trim(actual,"\n")
	if actual != expected {
		t.Errorf("'%s' was not equal to '%s'", actual, expected )
	}
}

func BenchmarkConcatenator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		expected := makeExpectedResponse(100)
		ts := makeTestServer(100)
		defer ts.Close()
		testUrls := makeTestUrls(ts.URL, 100)
		actual, err := Concatenator(testUrls...)
		if err != nil {
			b.Error("Got unexpected error", err)
		}
		actual = strings.Trim(actual,"\n")
		if actual != expected {
			b.Errorf("'%s' was not equal to '%s'", actual, expected )
		}
	}
}

func TestGetOneErr(t *testing.T) {
	ts := makeTestServer(0)
	defer ts.Close()
	pathWithInvalidHex := "/%zz"
	_, err := get(ts.URL + pathWithInvalidHex)
	if err == nil {
		t.Error("Didn't get expected error")
	}
}

func TestConcatenatorErr(t *testing.T) {
	ts := makeTestServer(0)
	defer ts.Close()
	_, err := Concatenator(ts.URL +"/x")
	if err == nil {
		t.Error("Didn't get expected error")
	}
}

func makeTestServer(numResponses int) *httptest.Server {
	expected := make(map[string]string)
	for i := 1; i <= numResponses; i++ {
		expected[fmt.Sprintf("/%d", i)]=fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		expectedResp, ok := expected[r.URL.Path]
		if ok {
			fmt.Fprintln(w, expectedResp)
		} else {
			http.NotFound(w,r)
		}

	}))
}

func makeExpectedResponse(numResponses int) (expectedResponse string) {
	for i := 1; i <= numResponses; i++ {
		expectedResponse+=fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)
	}
	return
}

func makeTestUrls(baseUrl string, numResponses int) (testUrls []string) {
	for i := 1; i <= numResponses; i++ {
		testUrls = append(testUrls, makeTestUrl(baseUrl, i))
	}
	return
}
func makeTestUrl(baseUrl string, num int) (testUrl string) {
	return fmt.Sprintf("%s/%d", baseUrl, num)
}
