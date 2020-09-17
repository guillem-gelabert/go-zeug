package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

// Game holds the application's shared state
type Game struct {
	host         string
	jwt          string
	session      *[]dto.CardDTO
	currentCards *[]dto.CardDTO
	reader       *bufio.Reader
}

func (game *Game) play() error {
	for _, card := range *game.session {
		fmt.Printf("%s (write \"<der|die|das)>\" %s\n", card.Substantive, card.Substantive)
		text, _ := game.reader.ReadString('\n')
		text = strings.Split(text, "\n")[0]
		if strings.ToLower(text) != strings.ToLower(card.Article) {
			err := game.answer(card.ID, false)
			if err != nil {
				return err
			}
			fmt.Print("WRONG ")
		} else {
			err := game.answer(card.ID, true)
			fmt.Print("RIGHT ")
			if err != nil {
				return err
			}
		}
		fmt.Printf("You said %q it is %q\n", text, card.Article)
	}
	return nil
}

func (game *Game) answer(id int, correct bool) error {
	answer := &dto.AnswerDTO{
		ID:      id,
		Correct: correct,
	}
	body, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	http.Post(game.host+"/cards", "application/json", bytes.NewBuffer(body))
	return nil
}
