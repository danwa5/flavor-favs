// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.

package main

import (
    "context"
    "fmt"
    "github.com/joho/godotenv"
    "github.com/zmb3/spotify/v2"
    "github.com/zmb3/spotify/v2/auth"
    "log"
    "net/http"
    "os"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
    err = godotenv.Load()

    auth  = spotifyauth.New(
                spotifyauth.WithClientID(os.Getenv("SPOTIFY_ID")),
                spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_SECRET")),
                spotifyauth.WithRedirectURL(redirectURI),
                spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead))
    ch    = make(chan *spotify.Client)
    state = "abc123"
)

func init() {
    // err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: ", err)
    }
}

func main() {
    // first start an HTTP server
    http.HandleFunc("/callback", completeAuth)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Got request for:", r.URL.String())
    })
    go func() {
        err := http.ListenAndServe(":8080", nil)
        handleError(err, "Error starting server on port 8080")
    }()

    url := auth.AuthURL(state)
    fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

    // wait for auth to complete
    client := <-ch

    ctx := context.Background()

    // use the client to make calls that require authorization
    user, err := client.CurrentUser(ctx)
    handleError(err, "Error fetching current user")
    fmt.Println("You are logged in as:", user.ID)

    topArtists, err := client.CurrentUsersTopArtists(ctx)
    handleError(err, "Error fetching current user's top artists")
    fmt.Println("\nTOP ARTISTS (medium term)")

    for index, artist := range topArtists.Artists {
        fmt.Printf("%d. %s (%d)\n", index+1, artist.Name, artist.Popularity)
    }
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
    tok, err := auth.Token(r.Context(), state, r)
    if err != nil {
        http.Error(w, "Couldn't get token", http.StatusForbidden)
        log.Fatal(err)
    }
    if st := r.FormValue("state"); st != state {
        http.NotFound(w, r)
        log.Fatalf("State mismatch: %s != %s\n", st, state)
    }

    // use the token to get an authenticated client
    client := spotify.New(auth.Client(r.Context(), tok))
    fmt.Fprintf(w, "Login Completed!")
    ch <- client
}

func handleError(err error, message string) {
    if err != nil {
        log.Fatalf(message + "\n%v", err)
    }
}
