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
	role := flag.String("set-role", "", "set existing token role: admin or user (does not rotate the token)")
	flag.Parse()
	if *user == "" {
		log.Fatal("-user is required")
	}
	s, e := dashboard.NewSyncStore(*dir)
	if e != nil {
		log.Fatal(e)
	}
	if *role != "" {
		if *revoke {
			log.Fatal("-revoke and -set-role cannot be combined")
		}
		if e = s.SetRole(*user, *role); e != nil {
			log.Fatal(e)
		}
		fmt.Println("Updated token role for", *user, "to", *role)
		return
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
