package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/zmb3/spotify"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/clientcredentials"
)

func main() {
    accessToken := getAccessToken()
    client := spotify.Authenticator{}.NewClient(accessToken)

    // https://open.spotify.com/playlist/37i9dQZEVXcO0WnzMtBl3g?si=9894c660aeb846cb
    playlistID := "37i9dQZEVXcO0WnzMtBl3g"
    spotifyPlaylistID := spotify.ID(playlistID)
    playlist, err := client.GetPlaylist(spotifyPlaylistID)
    handleError(err, "Error fetching playlist " + playlistID)

    fmt.Println("playlist id:", playlist.ID)
    fmt.Println("playlist name:", playlist.Name)
    fmt.Println("playlist description:", playlist.Description)
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
