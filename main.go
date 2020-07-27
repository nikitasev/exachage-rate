package main

import (
	cb "exachage-rate/coinbase"
	db "exachage-rate/database"
	"fmt"
)

func main() {
	dbService := db.InitDBService()
	defer dbService.Con.Close()

	cbConnection := cb.CreateConnection()
	defer cbConnection.Close()

	subscribe := cb.Message{
		Type: "subscribe",
		Channels: []cb.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					cb.EthBtc,
					cb.BtcEur,
					cb.BtcUsd,
				},
			},
		},
	}

	if err := cb.SubscribeMethods(cbConnection, &subscribe); err != nil {
		panic(err.Error())
	}

	msgChan := [3]chan cb.Message{}

	// Create 3 stream to handle responses from websocket
	for i := 0; i < 3; i++ {
		msgChan[i] = make(chan cb.Message)
		go saveExchangeRate(msgChan[i], &dbService)
		defer close(msgChan[i])
	}

	for {
		response := cb.Message{}
		if err := cbConnection.ReadJSON(&response); err != nil || response.Type == "error" {
			println(err.Error())
			continue
		}

		switch response.ProductID {
		case cb.EthBtc:
			msgChan[0] <- response
		case cb.BtcUsd:
			msgChan[1] <- response
		case cb.BtcEur:
			msgChan[2] <- response
		}
	}
}

func saveExchangeRate(msgCh chan cb.Message, dbService *db.DbService) {
	for response := range msgCh {
		dbService.Locker.Lock()
		insert, err := dbService.Con.Query(
			"INSERT INTO ticks (timestamp, symbol, bid, ask)"+
				"VALUE (?, ?, ?, ?)", response.Time.Unix(), response.ProductID, response.BestBid, response.BestAsk)

		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Printf("Symbol : %v\n", response.ProductID)

		insert.Close()
		dbService.Locker.Unlock()
	}

	return
}
