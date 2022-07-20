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
    "errors"
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

    itemType string
    limit int

    commandMap map[string]interface{}
    rangeOptions = [3]spotify.Range{spotify.ShortTermRange, spotify.MediumTermRange, spotify.LongTermRange}
)

func init() {
    handleError(err, "Error loading .env file")

    // parse command line options
    flag.StringVar(&itemType, "type", "artists", "\"artists\" or \"tracks\"")
    flag.IntVar(&limit, "limit", 10, "the number of results per data set")
    flag.Parse()
    validateOptions(itemType, limit)
}

func main() {
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

    // create map of item types and spotify functions
    commandMap = map[string]interface{} {
        "artists": client.CurrentUsersTopArtists,
        "tracks": client.CurrentUsersTopTracks,
    }

    for _, ro := range rangeOptions {
        timeRangeOpt := spotify.Timerange(ro)
        limitOpt := spotify.Limit(limit)
        res, _ := Call(itemType, ctx, timeRangeOpt, limitOpt)

        fmt.Printf("\nTOP %d %v (%s)\n", limit, itemType, ro)

        if itemType == "artists" {
            results := res.(*spotify.FullArtistPage)
            for index, artist := range results.Artists {
                fmt.Printf("%d. %s (%d)\n", index+1, artist.Name, artist.Popularity)
            }
        } else {
            results := res.(*spotify.FullTrackPage)
            for index, track := range results.Tracks {
                fmt.Printf("%d. %s %s (%d)\n", index+1, track.Name, buildArtistSentence(track.Artists), track.Popularity)
            }
        }
    }
}

func buildArtistSentence(artists []spotify.SimpleArtist) (result string) {
    artistSentence := "by "
    for i, artist := range artists {
        if i > 0 {
            artistSentence += ", "
        }
        artistSentence += artist.Name
    }
    return artistSentence
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

func validateOptions(itemType string, limit int) {
    if itemType != "artists" && itemType != "tracks" {
        err = errors.New("-type must be one of the following: artists, tracks")
        handleError(err, "Invalid option")
    }

    if limit < 1 || limit > 50 {
        err = errors.New("-limit must be between 1 and 50")
        handleError(err, "Invalid option")
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
