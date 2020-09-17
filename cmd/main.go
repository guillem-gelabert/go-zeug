package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

func main() {
	host := flag.String("h", "http://localhost:8000", "Host URL")
	email := flag.String("u", "", "Email to login")
	password := flag.String("p", "", "Password to login")

	flag.Parse()

	jwt, err := authenticate(*host, *email, *password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully logged in as %s\n", *email)
	session, err := getSession(*host, jwt)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		host:    *host,
		jwt:     jwt,
		reader:  bufio.NewReader(os.Stdin),
		session: session,
	}

	err = game.play()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You're done for today")
}

func getSession(host, jwt string) (*[]dto.CardDTO, error) {
	req, err := http.NewRequest("GET", host+"/cards", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+jwt)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cards []dto.CardDTO

	err = json.Unmarshal(body, &cards)
	if err != nil {
		return nil, err
	}

	return &cards, nil
}

func authenticate(host, email, password string) (string, error) {
	loginDTO := dto.LoginDTO{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(loginDTO)
	if err != nil {
		return "", err
	}

	res, err := http.Post(host+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(res.Status)
	}

	jwt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(jwt), nil
}
