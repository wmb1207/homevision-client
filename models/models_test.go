package models

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

func testRequester(url string) (*http.Response, error) {
	fileReader, err := os.Open("./test_file.md")
	if err != nil {
		panic(err)
	}

	resp := &http.Response{Body: fileReader}
	return resp, nil
}

func TestDownloadPhoto(t *testing.T) {
	house := House{Id: 1, Address: "123 Main St", Homeowner: "John Doe", PhotoUrl: "http://example.com/photo.jpg"}

	got, _ := house.DownloadPhoto(testRequester)
	want := "1-photo.jpg"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	gotFile, err := os.Open(want)
	if err != nil {
		t.Errorf("File %s not found", want)
	}

	testResponse, err := testRequester("http://example.com/photo.jpg")
	if err != nil {
		t.Errorf("Error getting test file")
	}

	gotFileContent, err := io.ReadAll(gotFile)
	if err != nil {
		t.Errorf("Error reading file %s", want)
	}

	wantedFileContent, err := io.ReadAll(testResponse.Body)
	if err != nil {
		t.Errorf("Error reading test file")
	}

	if bytes.Equal(wantedFileContent, append(gotFileContent, '\n')) { // Correct for editors adding a new line at the end of the test file
		t.Errorf("got %s, want %s", gotFileContent, wantedFileContent)
	}
	os.Remove(want)
}

func TestGetFileName(t *testing.T) {
	house := House{Id: 2, Address: "123 Main St", Homeowner: "John Doe", PhotoUrl: "http://example.com/photo.jpg"}

	got := house.GetFileName()
	want := "2-photo.jpg"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
