# Weather API Service

A simple Weather API service built with Golang that fetches weather information for a given city using the OpenWeather API.

## Features
- Fetch weather data for any city.
- Simple RESTful API endpoints.
- Uses `.env` file for API key configuration.
- Graceful error handling.

## Prerequisites
- Golang installed (Go 1.18+ recommended).
- OpenWeather API key.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/weather-api.git
   cd weather-api
   ```

2. Create a `.env` file in the root directory and add your OpenWeather API key:
   ```sh
   OPENWEATHER_API_KEY=your_api_key_here
   ```

3. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

### Start the Server
Run the following command to start the server:
```sh
 go run main.go
```

The server will run on `http://localhost:8080`

### API Endpoints

#### 1. Welcome Message
**GET** `/hello`
- Response:
  ```json
  {
    "message": "Hi, welcome to Weather API BOT!"
  }
  ```

#### 2. Get Weather for a City
**GET** `/weather/{city_name}`
- Example:
  ```sh
  curl http://localhost:8080/weather/London
  ```
- Response:
  ```json
  {
    "name": "London",
    "main": {
      "temp": 289.15
    }
  }
  ```

## Deployment
To deploy the application, build the executable and run it:
```sh
 go build -o weather-api
 ./weather-api
```

## Contributing
Feel free to fork the repository and submit pull requests.

## License
This project is licensed under the MIT License.

