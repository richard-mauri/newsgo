package main

import (
	"flag"
	"fmt"
	"github.com/richard-mauri/newsgo/internal/dashboard"
	"log"
)

func main() {
	dir := flag.String("dir", "data", "NewsGo server data directory")
	user := flag.String("user", "", "unique user ID")
	label := flag.String("label", "", "optional display label")
	revoke := flag.Bool("revoke", false, "revoke the named user token")
	flag.Parse()
	if *user == "" {
		log.Fatal("-user is required")
	}
	s, e := dashboard.NewSyncStore(*dir)
	if e != nil {
		log.Fatal(e)
	}
	if *revoke {
		if e = s.Revoke(*user); e != nil {
			log.Fatal(e)
		}
		fmt.Println("Revoked", *user)
		return
	}
	token, e := s.IssueToken(*user, *label)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("User:", *user, "\nToken (copy now; never stored in plaintext):", token)
}
