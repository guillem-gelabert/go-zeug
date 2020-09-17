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
	"strings"

	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

// command line tool
// zeug -u guillem@gmail.com -p blablabla
// > help
// get	see all your scheduled cards
// play	answer cards
// > get
// show scheduled cards
// id	article	substantive	dateScheduled
// > play
//  Haus (write "<der|die|das)> Haus")
// > der Haus
// Wrong! DAS Haus. Press ENTER for next.
// Land (write "<der|die|das)> Land")
// > der Land
// Right! DER Land. Press ENTER for next.
// You're done for today

func main() {
	host := flag.String("h", "http://localhost:8000", "Host URL")
	email := flag.String("u", "", "Email to login")
	password := flag.String("p", "", "Password to login")

	flag.Parse()
	reader := bufio.NewReader(os.Stdin)

	jwt, err := authenticate(*host, *email, *password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully logged in as %s\n", *email)

	session, err := getSession(*host, jwt)
	if err != nil {
		log.Fatal(err)
	}

	for _, card := range *session {
		fmt.Printf("%s (write \"<der|die|das)>\" %s\n", card.Substantive, card.Substantive)
		text, _ := reader.ReadString('\n')
		text = strings.Split(text, "\n")[0]
		if strings.ToLower(text) != strings.ToLower(card.Article) {
			err = answer(*host+"/card", jwt, card.ID, false)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print("WRONG ")
		} else {
			err = answer(*host+"/card", jwt, card.ID, true)
			fmt.Print("RIGHT ")
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("You said %q it is %q\n", text, card.Article)
	}

	fmt.Println("You're done for today")
}

func answer(host string, jwt string, id int, correct bool) error {
	answer := &dto.AnswerDTO{
		ID:      id,
		Correct: correct,
	}
	body, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	http.Post(host+"/cards", "application/json", bytes.NewBuffer(body))
	return nil
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
