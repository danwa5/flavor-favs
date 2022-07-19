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
    "flag"
    "fmt"
    "github.com/joho/godotenv"
    "github.com/zmb3/spotify/v2"
    "github.com/zmb3/spotify/v2/auth"
    "log"
    "net/http"
    "os"
    "reflect"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
    // load env vars before instantiating global vars
    err = godotenv.Load()

    auth  = spotifyauth.New(
                spotifyauth.WithClientID(os.Getenv("SPOTIFY_ID")),
                spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_SECRET")),
                spotifyauth.WithRedirectURL(redirectURI),
                spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserTopRead))
    ch    = make(chan *spotify.Client)
    state = "abc123"

    commandMap map[string]interface{}
    rangeOptions = [3]spotify.Range{spotify.ShortTermRange, spotify.MediumTermRange, spotify.LongTermRange}
)

func init() {
    handleError(err, "Error loading .env file")
}

func main() {
    // load command line options
    topItemType := flag.String("type", "artists", "\"artists\" or \"tracks\"")
    flag.Parse()

    // start a HTTP server
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
    fmt.Printf("You are logged in as: %s (%s)\n", user.DisplayName, user.ID)

    // create map of top item types and spotify functions
    commandMap = map[string]interface{} {
        "artists": client.CurrentUsersTopArtists,
        "tracks": client.CurrentUsersTopTracks,
    }

    for _, ro := range rangeOptions {
        timeRange := spotify.Timerange(ro)
        res, _ := Call(*topItemType, ctx, timeRange)

        fmt.Printf("\nTOP %v (%s)\n", *topItemType, ro)

        if *topItemType == "artists" {
            results := res.(*spotify.FullArtistPage)
            for index, artist := range results.Artists {
                fmt.Printf("%d. %s (%d)\n", index+1, artist.Name, artist.Popularity)
            }
        } else {
            results := res.(*spotify.FullTrackPage)
            for index, track := range results.Tracks {
                fmt.Printf("%d. %s (%d)\n", index+1, track.Name, track.Popularity)
            }
        }
    }
}

func Call(funcName string, params ...interface{}) (result interface{}, err error) {
    f := reflect.ValueOf(commandMap[funcName])

    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    var res []reflect.Value
    res = f.Call(in)

    // get and return value represented by reflect.Value
    result = res[0].Interface()
    return
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
