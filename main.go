package main

import (
	cb "exachage-rate/coinbase"
	e "exachage-rate/entity"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql"
)

func main() {
	if err := initDb(); err != nil {
		panic(err.Error())
	}

	connection, err := cb.CreateConnection()

	if err != nil {
		panic(err.Error())
	}

	defer connection.Close()

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

	if err := cb.SubscribeMethods(connection, &subscribe); err != nil {
		panic(err.Error())
	}

	// Create 3 stream to handle responses from websocket
	ch1 := make(chan cb.Message)
	ch2 := make(chan cb.Message)
	ch3 := make(chan cb.Message)

	go saveExchangeRate(ch1)
	go saveExchangeRate(ch2)
	go saveExchangeRate(ch3)

	for {
		response := cb.Message{}
		if err := connection.ReadJSON(&response); err != nil || response.Type == "error" {
			println(err.Error())
			continue
		}

		switch response.ProductID {
		case cb.EthBtc:
			ch1 <- response
		case cb.BtcUsd:
			ch2 <- response
		case cb.BtcEur:
			ch3 <- response
		}
	}
}

func initDb() error {
	configs := make(map[string]*gosql.Config)
	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     "root@tcp(127.0.0.1:3306)/test",
		ShowSql: true,
	}

	if err := gosql.Connect(configs); err != nil {
		return err
	}

	return nil
}

func saveExchangeRate(ch chan cb.Message) {
	for {
		response := <-ch

		if _, err := gosql.Model(
			&e.Ticks{
				Timestamp: response.Time.Unix(),
				Symbol:    response.ProductID,
				Bid:       response.BestBid,
				Ask:       response.BestAsk,
			}).Create(); err != nil {
			fmt.Print(err.Error())
			continue
		}
	}
}
