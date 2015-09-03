`go.google.musicmanager` is a Go port of the [Musicmanager interface][1]
of [Simon Weber's unofficial Google Music API][2].  It can be used to
list, upload, and download songs in a Google Play Music library.  The
package also comes with the toy command-line client `gmusic` for testing
this functionality.

Note that use of this package likely constitutes a violation of
[Google's Terms of Service, section 2][3], paragraph 2: "For example,
don't interfere with our Services or try to access them using a method
other than the interface and the instructions that we provide."  Use at
your own risk.

Use of this package requires a Google OAuth 2.0 token with a scope of
`https://www.googleapis.com/auth/musicmanager`.  See [here][4] for
instructions on how to obtain one.  Note that the above scope need not
be manually activated for a Google Cloud project (in fact, it's not even
listed on the APIs page).

 - [API reference](https://godoc.org/github.com/lxr/go.google.musicmanager)
 - [Bugs and gotchas](https://godoc.org/github.com/lxr/go.google.musicmanager#pkg-note-bug)

[1]: https://unofficial-google-music-api.readthedocs.org/en/latest/reference/musicmanager.html
[2]: https://github.com/simon-weber/gmusicapi
[3]: https://www.google.com/intl/en/policies/terms/#toc-services
[4]: https://godoc.org/golang.org/x/oauth2/google
