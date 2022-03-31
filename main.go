package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"crypto/tls"

	"github.com/gliderlabs/ssh"
	"github.com/patrickmn/go-cache"
)

type Attacker struct {
	addr     string
	username string
	password string
}

func authHandler(ctx ssh.Context, password string) bool {
	if len(ctx.User()) > 0 {
		log.Printf("User: %s connecting from %s with password: %s\n",
	ctx.User(), ctx.RemoteAddr(), password)
	}
	return true
}

func (a *Attacker) String() string {
	if len(a.username) > 0 {
		return fmt.Sprintf("%s - %s:%s", a.addr, a.username, a.password)
	}
	return fmt.Sprintf("%s - SSH Key Attempt", a.addr)
}

func sessionHandler(s ssh.Session) {
	io.WriteString(s, "Welcome!\n")
}

var c = cache.New(400*time.Minute , 300*time.Second)

func notify(attacker *Attacker) {
	log.Println("Attempt", attacker.String())
	_, found := c.Get(attacker.addr)
	if !found {
		go pushNotify(attacker)
		c.Set(attacker.addr, 1, 0)
	}
}

var conf = &Configuration{}

type Configuration struct {
	UserId string
	Token  string
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



func main() {
	file, err := os.Open("./conf.json")
	if err != nil {
		log.Fatal("Error opening config file")
	}
	decoder := json.NewDecoder(file)
	decoder.Decode(&conf)
	
	s := &ssh.Server {
		Addr:            ":2222",
		Handler:         sessionHandler,
		PasswordHandler: authHandler,
	}

	log.Println("Starting ssh server on port 2222..")
	log.Fatal(s.ListenAndServe())
}
