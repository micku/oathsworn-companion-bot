package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command = func(params []string, repliedMessage string) string
type Commands = map[string]Command

type Die = [6]string
type Dice = map[string]Die

var availableCommands = Commands{ "/roll": roll, "/reroll": reroll }
var availableDice = Dice{
	"w": [6]string{"x", "x", "1", "1", "2", "{2}"},
	"y": [6]string{"x", "x", "1", "2", "3", "{3}"},
	"r": [6]string{"x", "x", "2", "3", "3", "{4}"},
	"b": [6]string{"x", "x", "3", "3", "4", "{5}"},
}

var allowedChats []int64

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	allowedChatIds := os.Getenv("TG_ALLOWED_CHATS")
	if allowedChatIds != "" {
		for _, chatId := range strings.Split(os.Getenv("TG_ALLOWED_CHATS"), ",") {
			n, err := strconv.ParseInt(chatId, 10, 64)
			if err != nil {
				log.Panic(err)
				continue
			}

			allowedChats = append(allowedChats, n)
		}
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if len(allowedChats) > 0 && !slices.Contains(allowedChats, update.Message.Chat.ID) {
				log.Printf("Chat [%d] is not allowed", update.Message.Chat.ID)
				continue
			}

			words := strings.Fields(update.Message.Text)
			if len(words) < 2 {
				reply(update.Message.Chat.ID, update.Message.MessageID, bot, "I'm sorry, I don't understand that")
				continue
			}

			command, ok := availableCommands[words[0]]
			if !ok {
				reply(update.Message.Chat.ID, update.Message.MessageID, bot, "I'm sorry, I don't understand this command")
				continue
			}

			reply(update.Message.Chat.ID, update.Message.MessageID, bot, command(words[1:], getReplyToMessage(&update)))
		}
	}
}

func reply(chatId int64, messageId int, bot *tgbotapi.BotAPI, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyToMessageID = messageId

	bot.Send(msg)
}

func getReplyToMessage(update *tgbotapi.Update) string {
	if update.Message.ReplyToMessage == nil {
		return ""
	}

	return update.Message.ReplyToMessage.Text
}

func roll(params []string, repliedMessage string) string {
	var sequence = params[0]
	result, err := rollDice(sequence)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s\n%s", sequence, strings.Join(result, "\n"))
}

func reroll(params []string, repliedMessage string) string {
	result := strings.Split(repliedMessage, "\n")
	sequence := result[0]

	for _, s := range params {
		index, err := strconv.Atoi(s)
		if err != nil {
			return "Mmmhhhh... Are the dice indexes correct?"
		}

		die, ok := availableDice[string(sequence[index - 1])]
		if !ok {
			return "Mmmhhhh... Are the dice indexes correct?"
		}

		result[index] = fmt.Sprintf("%d. %s", index, rollDie(die))
	}

	return strings.Join(result, "\n")
}

func rollDice(sequence string) ([]string, error) {
	var result []string

	for i, char := range sequence {
		die, ok := availableDice[string(char)]
		if !ok {
			return nil, errors.New("Mmmhhhh... I don't know that die")
		}

		result = append(result, fmt.Sprintf("%d. %s", i + 1, rollDie(die)))
	}

	return result, nil
}

func rollDie(die Die) string {
	faceIndex := rand.Intn(len(die))
	result := die[faceIndex]
	if faceIndex == 5 {
		result += fmt.Sprintf(" -> %s", rollDie(die))
	}
	return result
}
