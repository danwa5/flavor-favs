# Flavor Favs

This is a Go CLI tool that uses a [wrapper](https://github.com/zmb3/spotify) to make calls to
the Spotify Web API to fetch your favorite artists and tracks. Your favorites from 3 time ranges
(4 weeks, 6 months, and several years) will be displayed.

## Getting Started

To allow Flavor Favs access to your personal Spotify data, you'll need to obtain a **client ID**
and **secret key** for the application to use.

1. Log in to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Click on *Create an App*
3. Register your app
- Enter an App Name
- Enter an App Description
- Click *Create*
4. After redirect, click on your app and take note of your **client ID** and **secret key**.
DO NOT share these with anyone!
5. Click on *Edit Settings* and enter a **redirect URI**. If you don't have one, just default
to http://localhost:8080/callback. This URI enables Spotify to invoke the Flavor Favs application
after successful authentication.
6. Set your environment variables. In your terminal, you can run:
```
$ export SPOTIFY_CLIENT_ID=<YOUR CLIENT ID>
$ export SPOTIFY_CLIENT_SECRET=<YOUR SECRET KEY>
$ export SPOTIFY_REDIRECT_URI=<YOUR REDIRECT URI>
```

## Testing and Building the Go Binary

Install Golang: https://go.dev/doc/install

Clone the repo
```
$ git clone <REPO>
```

Test the application
```
$ go run main.go
```

Build the executable binary
```
$ go build
```

## Usage

Before you can run the executable binary, make sure it has execute privileges
```
$ chmod 744 flavor-favs
```

Run the app in the directory of the `flavor-favs` binary. If you're using a Mac,
you may get a "Developer cannot be verified" warning, in which case you'll need
to go to your System Preferences to allow permissions for the program to run.

See the available command line options:
```
$ ./flavor-favs -h
Usage of ./flavor-favs:
  -limit int
        the number of results per data set (default 10)
  -type string
        "artists" or "tracks" (default "artists")
```

See your favorite artists:
```
$ ./flavor-favs -type artists
```

See your favorite tracks:
```
$ ./flavor-favs -type tracks
```

Change the number of results in each data set:
```
$ ./flavor-favs -type artists -limit 20
```

## Helpful Resources
- [Spotify Web API](https://developer.spotify.com/documentation/web-api/)
- [Spotify Web API Wrapper Library](https://github.com/zmb3/spotify)
- [Oauth2 Authentication](https://developer.spotify.com/documentation/general/guides/authorization/)
- [Golang Reflection Call](https://medium.com/@vicky.kurniawan/go-call-a-function-from-string-name-30b41dcb9e12)
