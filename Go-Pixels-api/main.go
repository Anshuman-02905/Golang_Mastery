package main

import (
	"fmt"
	"log"

	"github.com/Anshuman-02905/Golang-Mastery/Go-Pixels-api/handler"
)

func main() {
	//LOAD THE ENV FILE

	ac := handler.NewAuthenticatedClient()
	//results, err := ac.SearchPhotos("waves", "15", "1")
	//results, err := ac.CuratedPhotos("1", "1")
	//results, err := ac.GetPhoto("2014422")
	//results, err := ac.SearchVideos("waves", "1", "2")

	//results, err := ac.PopularVideos("1")

	results, err := ac.GetVideos("2499611")

	fmt.Println(ac.RemainingTimes)
	if err != nil {
		log.Fatalf("Unable to SeatchPhotos %v", err)
	}
	fmt.Println(results)

}
