package handler

// Why Use Struct Methods for API Handling?

// Feature	Struct Method (c.FetchData())	Normal Function (FetchData())
// Encapsulation	✅ Client state is managed inside struct	❌ No client state, needs manual passing
// Reusability	✅ Can be reused with different API calls	✅ Also reusable, but with manual parameters
// Dependency Injection	✅ Inject a custom HTTP client for testing	❌ Harder to inject a mock client
// Readability	✅ Organized under Client struct	❌ Scattered across multiple function

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Anshuman-02905/Golang-Mastery/Go-Pixels-api/config"
	"github.com/Anshuman-02905/Golang-Mastery/Go-Pixels-api/models"
)

const (
	PhotoApi = "https://api.pexels.com/v1/"
	VideoApi = "https://api.pexels.com/videos/"
)

type AuthenticatedClient struct {
	*models.Client
}

func NewAuthenticatedClient() *AuthenticatedClient {
	c := config.NewClient()

	return &AuthenticatedClient{Client: c}
}
func CreateClient() {

}

func (ac *AuthenticatedClient) RequestwithAuth(method, url string) (*http.Response, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Request did not generate %v", err)
		return nil, err
	}
	req.Header.Add("Authorization", ac.Token)

	resp, err := ac.HC.Do(req)

	if err != nil {
		log.Fatalf("Response did not generate %v", err)
		return resp, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		log.Fatal("error in fetching X-Ratelimit-Remaining")
		return resp, err
	} else {
		ac.RemainingTimes = int32(times)
	}
	return resp, nil
}

func (ac *AuthenticatedClient) SearchPhotos(query, per_page, page string) (*models.SearchResults, error) {

	url := fmt.Sprintf("search?query=%v&per_page=%v&page%v", query, per_page, page)
	fmt.Println(url)
	resp, err := ac.RequestwithAuth("GET", url)
	if err != nil {
		log.Fatalf("Error at executing request %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error at Reading the  response %v", err)
		return nil, err
	}
	var result models.SearchResults
	json.Unmarshal(data, &result)
	return &result, err

}

func (ac *AuthenticatedClient) CuratedPhotos(page, per_page string) (*models.CuratedPhotos, error) {

	url := fmt.Sprintf(PhotoApi+"curated?per_page=%v&page=%v", per_page, page)
	fmt.Println(url)
	resp, err := ac.RequestwithAuth("GET", url)
	if err != nil {
		log.Fatalf("Error at executing request %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error at Reading the  response %v", err)
		return nil, err
	}
	var result models.CuratedPhotos
	json.Unmarshal(data, &result)
	return &result, err

}
func (ac *AuthenticatedClient) GetPhoto(id string) (*models.Photo, error) {
	url := fmt.Sprintf(PhotoApi + "photos/" + id)
	fmt.Println(url)
	resp, err := ac.RequestwithAuth("GET", url)

	if err != nil {
		log.Fatalf("Error at executing request %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro At getting the respnse %v", err)
		return nil, err
	}
	var result models.Photo
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Error at decoding the response data %v", err)
		return nil, err
	}
	return &result, nil

}

func (ac *AuthenticatedClient) SearchVideos(query, per_page, page string) (*models.VideoSeach, error) {
	url := fmt.Sprintf(VideoApi+"search?query=%v&per_page=%v&page=%v", query, per_page, page)
	fmt.Println(url)
	resp, err := ac.RequestwithAuth("GET", url)
	if err != nil {
		log.Fatalf("error at gettting the response %v", err)
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading the Response  %v", err)
		return nil, err
	}
	var results models.VideoSeach

	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Fatalf("error at Unmarshal the response data  %v", err)
		return nil, err
	}

	return &results, nil
}

func (ac *AuthenticatedClient) PopularVideos(per_page string) (*models.PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"popular?per_page=%v", per_page)
	resp, err := ac.RequestwithAuth("GET", url)
	if err != nil {
		log.Fatalf("error at Requesting Popular Videos %v", err)
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error at Reading the response of Popular Videos %v", err)
		return nil, err
	}

	var results models.PopularVideos
	err = json.Unmarshal(data, &results)

	if err != nil {
		log.Fatalf("error at Unmarshal the response data  %v", err)
		return nil, err
	}

	return &results, nil

}
func (ac *AuthenticatedClient) GetVideos(id string) (*models.Video, error) {

	url := fmt.Sprintf(VideoApi+"videos/%v", id)
	fmt.Println(url)
	resp, err := ac.RequestwithAuth("GET", url)

	if err != nil {
		log.Fatalf("error at Requesting Single Videos %v", err)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error at Reading the response of Popular Videos %v", err)
		return nil, err
	}
	var results models.Video
	err = json.Unmarshal(data, &results)

	if err != nil {
		log.Fatalf("error at Unmarshal the response data  %v", err)
		return nil, err
	}
	return &results, err

}
