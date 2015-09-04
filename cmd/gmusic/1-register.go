/*

Registering the client

Usage:

	gmusic register id [name]

Gmusic must be registered with the user's Google Play account before it
can be used to manage Google Play Music libraries.  The registration
process asks you to navigate to a special URL where you can grant access
permissions for the Google Play Music Manager.  Doing this gives you an
authorization code, which is then input to gmusic to register it.  Once
gmusic has been registered, it creates a file called ".gmusic.json" in
the user's home directory; other gmusic commands refer to this file for
their access credentials.

The ID under which you register gmusic in your Google Play Music library
needs to be unique on Google's side, so pick it reasonably randomly.
Remember that there are limits to how many devices a single account can
have authorized, with how many accounts a single device can be
authorized, and how many devices one account can deauthorize in a year,
so be careful in using this command.

Note that downloading tracks has been known to fail unless the ID is
sufficiently MAC address-like.  The exact threshold is unknown; perhaps
the server only checks for a colon.

If a human-readable name under which to register gmusic is not given,
it defaults to "gmusic".

*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/lxr/go.google.musicmanager"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google's official Music Manager's credentials
var conf = &oauth2.Config{
	ClientID:     "652850857958",
	ClientSecret: "ji1rklciNp2bfsFJnEH_i6al",
	RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	Scopes:       []string{musicmanager.Scope},
	Endpoint:     google.Endpoint,
}

var creds struct {
	ID string
	oauth2.Token
}

func init() {
	cmds["register"] = register
}

func register() error {
	var id, name string
	switch len(os.Args) {
	case 2:
		id, name = os.Args[1], "gmusic"
	case 3:
		id, name = os.Args[1], os.Args[2]
	default:
		return errors.New("takes 1-2 arguments")
	}

	url := conf.AuthCodeURL("", oauth2.AccessTypeOffline)
	logf(`Please open the following URL in a browser to
authorize gmusic with your Google account.  Copy the code given to you
at the end of the authorization process below.

%s

> `, url)
	var code string
	if _, err := fmt.Scanln(&code); err != nil {
		return err
	}
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}
	httpclient := conf.Client(oauth2.NoContext, tok)
	client, err := musicmanager.NewClient(httpclient, id)
	if err != nil {
		return err
	}
	if err := client.Register(name); err != nil {
		return err
	}
	credsPath, err := getCredsPath()
	if err != nil {
		return err
	}
	f, err := os.Create(credsPath)
	if err != nil {
		return err
	}
	defer f.Close()
	creds.ID = id
	creds.Token = *tok
	if err := json.NewEncoder(f).Encode(creds); err != nil {
		return err
	}
	logf("registration successful\n")
	return nil
}

func getCredsPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, ".gmusic.json"), nil
}

func loadClient() (*musicmanager.Client, error) {
	credsPath, err := getCredsPath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(credsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&creds); err != nil {
		return nil, err
	}
	client := conf.Client(oauth2.NoContext, &creds.Token)
	return musicmanager.NewClient(client, creds.ID)
}
