package models

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Requester func(string) (*http.Response, error)

type DownloadableFile interface {
	DownloadPhoto(Requester) (string, error)
	GetFileName() string
}

type House struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	PhotoUrl  string `json:"photoUrl"`
}

type HouseList struct {
	Houses []House `json:"houses"`
}

func (h House) GetFileName() string {
	filepathSlice := strings.Split(h.PhotoUrl, "/")
	return strconv.Itoa(h.Id) + "-" + filepathSlice[len(filepathSlice)-1]
}

func (h House) DownloadPhoto(requester Requester) (string, error) {
	// Get the data
	log.Println("Downloading photo from " + h.PhotoUrl)
	resp, err := requester(h.PhotoUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filepath := h.GetFileName()

	// Create the file
	out, err := os.Create("./" + filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return filepath, err
}
