package concatenator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetOne(t *testing.T) {
	ts := makeTestServer(1)
	defer ts.Close()
	actual, _ := get(makeTestUrl(ts.URL, 1))
	isCorrect(actual, t, 1)
}

func TestConcatenator(t *testing.T) {
	numUrls := 2
	ts := makeTestServer(numUrls)
	testUrls := makeTestUrls(ts.URL, numUrls)
	defer ts.Close()
	actual, _ := Concatenator(testUrls...)
	isCorrect(actual, t, numUrls)
}

func BenchmarkConcatenator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ts := makeTestServer(100)
		defer ts.Close()
		testUrls := makeTestUrls(ts.URL, 100)
		actual, _ := Concatenator(testUrls...)
		isCorrect(actual, b, 100)
	}
}

type Test interface {
	Errorf(format string, args ...interface{})
}

func isCorrect(actual string, t Test, numUrls int) {
	actual = strings.Trim(actual, "\n")
	for _, expected := range makeExpectedResponseParts(numUrls) {
		if !strings.Contains(actual, expected) {
			t.Errorf("'%s' does not contain '%s'", actual, expected)
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
	_, err := Concatenator(ts.URL + "/x")
	if err == nil {
		t.Error("Didn't get expected error")
	}
}

func makeTestServer(numResponses int) *httptest.Server {
	expected := make(map[string]string)
	for i := 1; i <= numResponses; i++ {
		expected[fmt.Sprintf("/%d", i)] = fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		expectedResp, ok := expected[r.URL.Path]
		if ok {
			fmt.Fprintln(w, expectedResp)
		} else {
			http.NotFound(w, r)
		}

	}))
}

func makeExpectedResponseParts(numResponses int) (expectedResponseParts []string) {
	for i := 1; i <= numResponses; i++ {
		expectedResponseParts = append(expectedResponseParts, fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i))
	}
	return
}

func makeTestUrls(baseUrl string, numResponses int) (testUrls []string) {
	for i := 1; i <= numResponses; i++ {
		testUrls = append(testUrls, makeTestUrl(baseUrl, i))
	}
	return
}

func makeTestUrl(baseUrl string, num int) string {
	return fmt.Sprintf("%s/%d", baseUrl, num)
}
