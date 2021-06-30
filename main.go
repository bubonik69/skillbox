package main

import (
	"encoding/json"
	"errors"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type bnResp struct{
	Price float64`json:"price,string"`
	Code int64 `json:"code"`
}
type wallet map[string]float64
var db =map[int64]wallet{}

func main() {

	bot, err := tgbotapi.NewBotAPI("1881313779:AAGCCDiMrcv48Ood8NJcMYhS7WZ0vsfED3Y")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Split(update.Message.Text, " ")
		switch command[0] {
		case "ADD":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "неверная команда"))
			}
				amount, err := strconv.ParseFloat(command[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,"неверная команда"))
				}
				if _,ok:=db[update.Message.Chat.ID];!ok {
					db[update.Message.Chat.ID] = wallet{}
				}
				db[update.Message.Chat.ID][command[1]] += amount
				textBalance:=fmt.Sprintf("Баланс %s:  %f",command[1],db[update.Message.Chat.ID][command[1]] )
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,textBalance))
//
		case "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "неверная команда"))
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,"неверная команда"))
			}
			if _,ok:=db[update.Message.Chat.ID];!ok {
				db[update.Message.Chat.ID] = wallet{}
			}
			db[update.Message.Chat.ID][command[1]] -= amount
			textBalance:=fmt.Sprintf("Баланс :  %f",db[update.Message.Chat.ID][command[1]] )
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,textBalance))

		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "неверная команда"))
			}
			delete(db[update.Message.Chat.ID], command[1])
		case "SHOW":
			msg := ""
			var sum float64
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getPrice(key)
				sum+= value*price
				msg += fmt.Sprintf("%s : %f [%f] summ\n", key, value, price*value)
			}
			msg += fmt.Sprintf("TOTAL : %f \n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "команда не найдена"))

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, command[0])
			//bot.Send(msg)
		}
	}
}


//https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT
func getPrice(symbol string) (price float64, err error){
	resp,err:=http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT",symbol))
	defer resp.Body.Close()
	var jsonResp bnResp
	err= json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err!=nil{
		return
	}
	if jsonResp.Code!=0 {
		err = errors.New("Неверный символ")
	}
	price=jsonResp.Price
	return
}
