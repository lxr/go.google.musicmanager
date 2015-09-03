// Gmusic is a toy command line client for managing Google Play Music
// libraries.
package main

// BUG(lor): There should probably be an option to disable diagnostic
// output in the download and upload commands.

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	os.Args = os.Args[1:]
	if len(os.Args) == 0 {
		os.Args = []string{""}
	}
	cmd := os.Args[0]
	f := map[string]func() error{
		"register": register,
		"list":     list,
		"download": download,
		"upload":   upload,
	}[cmd]
	if f == nil {
		fmt.Fprintln(os.Stderr, "usage: gmusic (register | list | download | upload)")
		os.Exit(1)
	}
	if err := f(); err != nil {
		fmt.Fprintf(os.Stderr, "gmusic: %s: %v\n", cmd, err)
		os.Exit(1)
	}
}

func getScanner() *bufio.Scanner {
	if len(os.Args) == 1 {
		return bufio.NewScanner(os.Stdin)
	}
	buf := new(bytes.Buffer)
	for _, arg := range os.Args[1:] {
		buf.WriteString(arg + "\n")
	}
	return bufio.NewScanner(buf)
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func println(s string) {
	fmt.Println(s)
}
