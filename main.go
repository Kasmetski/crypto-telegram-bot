package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kasmetski/cmcAPI"
)

//GetCoinPrice - Get price information from Coin Market Cap and make a message
func GetCoinPrice(s string) (msg string, err error) {
	coin, err := cmcAPI.GetCoinInfo(s)
	if err != nil {
		msg = "Can't find anything. Type the full name of the coin"
		return
	}

	name := fmt.Sprintf("Name: %s | %s | #%d\n", coin.Name, coin.Symbol, coin.Rank)
	price := fmt.Sprintf("PriceBTC: %f\nPriceUSD: %.2f\n", coin.PriceBtc, coin.PriceUsd)
	change := fmt.Sprintf("Change 1H/24H/7d: %.2f | %.2f | %.2f\n", coin.PercentChange1H, coin.PercentChange24H, coin.PercentChange7D)

	msg = name + price + change

	return
}

//GetCoinInfo - Get full infomartion from Coin Market Cap and make a message
func GetCoinInfo(s string) (msg string, err error) {
	coin, err := cmcAPI.GetCoinInfo(s)
	if err != nil {
		msg = "Can't find anything. Type the full name of the coin"
		return
	}

	name := fmt.Sprintf("Name: %s | %s | #%d\n", coin.Name, coin.Symbol, coin.Rank)
	supply := fmt.Sprintf("Available Supply: %d\n", int(coin.AvailableSupply))
	mcap := fmt.Sprintf("MarketCap: %d\n", int(coin.MarketCapUsd))
	volume := fmt.Sprintf("Volume: %d\n", int(coin.Two4HVolumeUsd))
	price := fmt.Sprintf("PriceBTC: %f\nPriceUSD: %.2f\n", coin.PriceBtc, coin.PriceUsd)
	change := fmt.Sprintf("Change 1H/24H/7d: %.2f | %.2f | %.2f\n", coin.PercentChange1H, coin.PercentChange24H, coin.PercentChange7D)

	msg = name + supply + mcap + volume + price + change

	return
}

func main() {
	bot, err := tgbotapi.NewBotAPI("TOKEN HERE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			commandArgs := update.Message.CommandArguments()

			switch update.Message.Command() {
			case "help":
				msg.Text = "there is no help for the people here\nbut you can try /status /info /price"
			case "status":
				msg.Text = "I'm ok. Thanks for checking"
			case "info":
				msg.Text, err = GetCoinInfo(commandArgs)
				if err != nil {
					log.Println(err)
				}
			case "price":
				msg.Text, err = GetCoinPrice(commandArgs)
				if err != nil {
					log.Println(err)
				}
			default:
				msg.Text = "I don't know that command, try help"
			}
			bot.Send(msg)
		}
	}
}
