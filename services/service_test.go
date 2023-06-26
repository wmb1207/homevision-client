package services

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wmb1207/homevision_client/models"
)

type TestStruct struct{}

func (t TestStruct) DownloadPhoto(requester models.Requester) (string, error) {
	return "test_file.md", nil
}

func (t TestStruct) GetFileName() string {
	return "test_file.md"
}

func testRequester(url string) (*http.Response, error) {
	fileReader, err := os.Open("./test_file.md")
	if err != nil {
		panic(err)
	}

	resp := &http.Response{Body: fileReader}
	return resp, nil
}

func TestDownloadRealPhoto(t *testing.T) {

	want := "1-googlelogo_color_272x92dp.png"

	if os.Getenv("CALL_REAL_API") != "TRUE" {
		return
	}
	url := "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"

	result := DownloadPhoto(models.House{Id: 1, PhotoUrl: url}, NewSafeHttpRequest(3))
	if !result.Status {
		t.Errorf("Error downloading photo")
	}
	got := result.Filename
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	os.Remove(want)
}


func TestQueryRealHouses(t *testing.T) {
	if os.Getenv("CALL_REAL_API") != "TRUE" {
		return
	}

	url := "http://app-homevision-staging.herokuapp.com/api_project/houses"
	result, err := QueryHouses(url, 1, 1, NewSafeHttpRequest(3))
	if err != nil {
		t.Errorf("Error querying houses")
	}
	if len(result) != 1 {
		t.Errorf("got %d, want %d", len(result), 1)
	}
	got := result[0].Id
	want := 0
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
	gotName := result[0].GetFileName()
	wantName := "0-big-custom-made-luxury-house-260nw-374099713.jpg"
	if gotName != wantName {
		t.Errorf("got %s, want %s", gotName, wantName)
	}
}

func TestDownloadRealHousePhoto(t *testing.T) {

	want := "1-big-custom-made-luxury-house-260nw-374099713.jpg"

	if os.Getenv("CALL_REAL_API") != "TRUE" {
		return
	}

	url := "https://image.shutterstock.com/image-photo/big-custom-made-luxury-house-260nw-374099713.jpg"
	result := DownloadPhoto(models.House{Id: 1, PhotoUrl: url}, NewSafeHttpRequest(3))
	if !result.Status {
		t.Errorf("Error downloading photo")
	}
	got := result.Filename
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	os.Remove(want)
}

func TestDownloadPhoto(t *testing.T) {
	testStruct := TestStruct{}
	got, _ := testStruct.DownloadPhoto(testRequester)
	want := "test_file.md"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestSafeHttpRequestTwoErrors(t *testing.T) {
	expected := "{'data': 'dummy'}"

	try := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if try < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			try++
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, expected)
	}))

	defer server.Close()

	url := server.URL
	response, err := NewSafeHttpRequest(3)(url)
	if err != nil {
		t.Errorf("Error getting data from %s", url)
	}

	body := response.Body
	defer body.Close()
	got, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Error reading response body")
	}

	if string(got) != expected {
		t.Errorf("got %s, want %s", got, expected)
	}
}

func TestSafeHttpRequestOneError(t *testing.T) {
	expected := "{'data': 'dummy'}"

	try := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if try == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			try++
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, expected)
	}))

	defer server.Close()

	url := server.URL
	response, err := NewSafeHttpRequest(3)(url)
	if err != nil {
		t.Errorf("Error getting data from %s", url)
	}

	body := response.Body
	defer body.Close()
	got, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Error reading response body")
	}

	if string(got) != expected {
		t.Errorf("got %s, want %s", got, expected)
	}
}
