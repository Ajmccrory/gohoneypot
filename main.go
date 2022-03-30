package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"log/syslog"
	"encoding/json"

	"github.com/gliderlabs/ssh"
)

type Attacker struct {
	addr	 string
	username string
	password string
}

func authHandler(ctx ssh.Context, a *Attacker) bool {
	if len(a.username) > 0 {
		return fmt.Sprintf("%s - %s:%s", a.addr, a.username, a.password)
	}
	return fmt.Sprintf("%s - SSH Key Attempt", a.addr)
}

func (a *Attacker) String() string {
	if len(a.username) > 0 {
		return fmt.Sprintf("%s - %s:%s", a.addr, a.username, a.password)
	}
	return fmt.Sprintf("%s - SSH Key Attempt", a.addr)
	return &Attacker{addr, username, password}
}

func sessionHandler(s ssh.Session) {
	io.WriteString(s, "Welcome!\n")
}

var c = cache.New(10**time.Minute, 300*time.Second)

func notify(attacker *Attacker) {
	log.Println("Attempt", attacker.String())
	_, found := c.Get(attacker.addr)
	if !found {
		go pushNotify(attacker)
		c.Set(attacker.addr, 1, 0)
	}
}

func pushNotify(attacker *Attacker) {
	if conf.Token == "" || conf.UserId == "" {
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://api.pushover.net/1/messages.json",
		url.Values{"token": {conf.Token}, "user": {conf.UserId}, "message": {attacker.String()}})
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(resp.Status)
}

var conf = &Configuration{}

type Configuration struct {
	UserID  string
	Token	string
}

func main() {
	file, err := os.Open("./conf.json")
	if err != nil {
		log.Fatal("Error opening config file")
	}
	decoder := json.NewDecoder(file)
	decoder.Decode(&conf)

	logwriter ,e := syslog.New(syslog.LOG_INFO, os.Args[0])
	if e == nil {
		log.SetOutput(logwriter)
	}

	s := &ssh.Sever {
		Addr:					":2222",
		Handler:		sessionHandler,
		PAsswordHandler:	authHandler,
	}
	log.Println("Starting ssh server on port 2222..")
	log.Fatal(s.ListenAndServe())
}