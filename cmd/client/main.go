package main

import (
	"github.com/wmb1207/homevision_client/services"
	"sync"
)

const URL = "http://app-homevision-staging.herokuapp.com/api_project/houses"

func main() {

	// TODO: Implement a way to pass the page number to start at and the amount of items per page that you want to download
	// TODO: Implement a way to pass the amount of tries that you want to make for each request
	wg := new(sync.WaitGroup)
	for page := 1; page < 11; page++ {
		houses, err := services.QueryHouses(URL, page, 10, services.NewSafeHttpRequest(3))
		if err != nil {
			panic(err)
		}
		for _, house := range houses {
			wg.Add(1)
			go services.DownloadPhotoAsync(wg, house, services.NewSafeHttpRequest(3))
		}
	}
	wg.Wait()

}
