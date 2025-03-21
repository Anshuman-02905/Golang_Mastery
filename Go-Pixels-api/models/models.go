package models

import (
	"net/http"
)

// Client represents an API client with authentication token, HTTP client, and remaining request quota.
type Client struct {
	Token          string      // API authentication token
	HC             http.Client // HTTP client for making requests
	RemainingTimes int32       // Number of remaining API calls allowed
}

// SearchResults represents the response structure for a search query on a photo API.
type SearchResults struct {
	Page         int     `json:"page"`          // Current page number
	Perpage      int     `json:"per_page"`      // Number of results per page
	TotalResults int     `json:"total_results"` // Total number of search results
	PrevPage     string  `json:"prev_page"`     // URL for the previous page
	NextPage     string  `json:"next_page"`     // URL for the next page
	Photos       []Photo `json:"photos"`        // List of photos returned in search results
}

// Photo represents an individual photo entity with its metadata.
type Photo struct {
	Id               int         `json:"id"`               // Unique identifier of the photo
	Width            int         `json:"width"`            // Width of the photo in pixels
	Height           int         `json:"height"`           // Height of the photo in pixels
	Url              string      `json:"url"`              // URL to the photo
	Photographer     string      `json:"photographer"`     // Name of the photographer
	Photographer_url string      `json:"photographer_url"` // URL to the photographer's profile
	Photographer_id  int         `json:"photographer_id"`  // Unique ID of the photographer
	AvgColor         string      `json:"avg_color"`        // Average color of the photo
	Src              PhotoSource `json:"src"`              // Different size versions of the photo
	Liked            bool        `json:"liked"`            // Whether the user has liked the photo
}

// PhotoSource represents different size versions of a photo for various use cases.
type PhotoSource struct {
	Original  string `json:"original"`  // Original high-resolution image URL
	Large2x   string `json:"large2x"`   // Double-sized large image URL
	Large     string `json:"large"`     // Large image URL
	Medium    string `json:"medium"`    // Medium-sized image URL
	Small     string `json:"small"`     // Small-sized image URL
	Portrait  string `json:"portrait"`  // Portrait-oriented image URL
	Landscape string `json:"landscape"` // Landscape-oriented image URL
	Tiny      string `json:"tiny"`      // Very small image URL
}

// CuratedPhotos represents a collection of curated photos.
type CuratedPhotos struct {
	Page         int     `json:"page"`          // Current page number
	PerPage      int     `json:"per_page"`      // Number of results per page
	Photos       []Photo `json:"photos"`        // List of curated photos
	NextPage     string  `json:"next_page"`     // URL for the next page
	TotalResults string  `json:"total_results"` // Total number of curated results
	PrevPage     string  `json:"prev_page"`     // URL for the previous page
}

// VideoSeach represents the response structure for a video search query.
type VideoSeach struct {
	Page         int     `json:"page"`          // Current page number
	PerPage      int     `json:"per_page"`      // Number of results per page
	TotalResults int     `json:"total_results"` // Total number of video results
	Url          string  `json:"url"`           // API URL for fetching more results
	Videos       []Video `json:"videos"`        // List of videos returned in search results
	PrevPage     string  `json:"prev_page"`     // URL for the previous page
	NextPage     string  `json:"next_page"`     // URL for the next page
}

// Video represents an individual video entity with its metadata.
type Video struct {
	Id       int    `json:"id"`       // Unique identifier of the video
	Width    int    `json:"width"`    // Width of the video in pixels
	Height   int    `json:"height"`   // Height of the video in pixels
	Url      string `json:"url"`      // URL to the video file
	Image    string `json:"image"`    // Thumbnail image URL
	Duration int    `json:"duration"` // Duration of the video in seconds
	User     User   `json:"user"`     // User who uploaded the video
}

// User represents an individual user profile.
type User struct {
	Id   int    `json:"id"`   // Unique identifier of the user
	Name string `json:"name"` // Name of the user
	Url  string `json:"url"`  // Profile URL of the user
}

// PopularVideos represents a collection of trending or popular videos.
type PopularVideos struct {
	Page         int     `json:"page"`          // Current page number
	PerPage      int     `json:"per_page"`      // Number of results per page
	TotalResults int     `json:"total_results"` // Total number of popular videos available
	Url          string  `json:"url"`           // API URL for fetching more results
	Videos       []Video `json:"videos"`        // List of popular videos
	PrevPage     string  `json:"prev_page"`     // URL for the previous page
	NextPage     string  `json:"next_page"`     // URL for the next page
}
