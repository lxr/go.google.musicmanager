# Introduction

Go.google.musicmanager is a Golang port of the [Musicmanager interface]
[1] of Simon Weber's [gmusicapi] [2].  It can be used to list, upload,
and download songs in a Google Play Music library.  The package also
comes with the gmusic toy command-line client for testing this
functionality.

Note that use of this package likely constitutes a violation of
[Google's Terms of Service] [3], section 2, paragraph 2: "For example,
don't interfere with our Services or try to access them using a method
other than the interface and the instructions that we provide."  Use at
your own risk.

# Usage

A complete API reference can be found on [godoc] [4].  For an example
project utilizing this library, see [google-musicmanager-web] [5].

Use of this package requires a Google OAuth 2.0 token with a scope of
`https://www.googleapis.com/auth/musicmanager`.  See [here] [6] for
instructions on how to obtain one.  Note that the above scope need not
be manually activated for a Google Cloud project (in fact, it's not even
listed on the APIs page).

# Bugs

See [godoc] [7].

[1]: https://unofficial-google-music-api.readthedocs.org/en/latest/reference/musicmanager.html
[2]: https://github.com/simon-weber/gmusicapi
[3]: https://www.google.com/intl/en/policies/terms/
[4]: https://godoc.org/github.com/lxr/go.google.musicmanager
[5]: https://github.com/lxr/google-musicmanager-web
[6]: https://godoc.org/golang.org/x/oauth2/google
[7]: https://godoc.org/github.com/lxr/go.google.musicmanager#pkg-note-bug
