package services

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/wmb1207/homevision_client/models"
)

type Result struct {
	Filename string
	Status   bool
}

func NewSafeHttpRequest(retry int) models.Requester {
	return func(url string) (*http.Response, error) {
		for try := 0; try < retry; try++ {
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Failed to get response from " + url)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Println("Failed to get response from " + url)
				continue
			}
			return resp, nil
		}
		return nil, errors.New("Failed to get response")
	}
}

func DownloadPhoto(downloadable models.DownloadableFile, requester models.Requester) Result {
	filepath, err := downloadable.DownloadPhoto(requester)
	if err != nil {
		return Result{Filename: "", Status: false}
	}
	return Result{Filename: filepath, Status: true}
}

func DownloadPhotoAsync(wg *sync.WaitGroup, downloadable models.DownloadableFile, requester models.Requester) {
	defer wg.Done()
	result := DownloadPhoto(downloadable, requester)
	if result.Status {
		log.Println("Successfully downloaded photo " + result.Filename)
		return
	}
	log.Println("Failed to download photo " + downloadable.GetFileName())
}

func QueryHouses(url string, page int, per_page int, requester models.Requester) ([]models.House, error) {
	finalUrl := url + "?page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(per_page)
	resp, err := requester(finalUrl)
	defer resp.Body.Close()

	// Create the house struct
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	houses := models.HouseList{Houses: []models.House{}}
	err = json.Unmarshal(body, &houses)
	if err != nil {
		return nil, err
	}
	return houses.Houses, nil
}
