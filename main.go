package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/zmb3/spotify"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/clientcredentials"
)

type Playlist struct {
    ID          string  `json:"id"`
    IsPublic    bool    `json:"public"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
}

func main() {
    accessToken := getAccessToken()
    fmt.Println("Access token: " + accessToken.AccessToken)

    // https://open.spotify.com/playlist/37i9dQZEVXcO0WnzMtBl3g?si=9894c660aeb846cb
    playlistID := "37i9dQZEVXcO0WnzMtBl3g"

    url := "https://api.spotify.com/v1/playlists/" + playlistID
    var bearer = "Bearer " + accessToken.AccessToken

    req, err := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", bearer)

    client := &http.Client{}
    resp, err := client.Do(req)
    handleError(err, "Error making HTTP request")
    defer resp.Body.Close()

    byteValue, err := ioutil.ReadAll(resp.Body)
    handleError(err, "Error while reading the response bytes")
    // log.Println(string([]byte(byteValue)))

    var playlist Playlist
    json.Unmarshal([]byte(byteValue), &playlist)

    fmt.Println("playlist id:", playlist.ID)
    fmt.Println("playlist name:", playlist.Name)
    fmt.Println("playlist description:", playlist.Description)
    fmt.Println("playlist public:", playlist.IsPublic)
}

func getAccessToken() (*oauth2.Token) {
    err := godotenv.Load()
    handleError(err, "Error loading .env file")

    authConfig := &clientcredentials.Config{
        ClientID: os.Getenv("SPOTIFY_CLIENT_ID"),
        ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
        TokenURL: spotify.TokenURL,
    }

    accessToken, err := authConfig.Token(context.Background())
    handleError(err, "Error retrieving access token")

    return accessToken
}

func handleError(err error, message string) {
    if err != nil {
        log.Fatalf(message + "\n%v", err)
    }
}
