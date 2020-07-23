package coinbase

import (
	"github.com/gorilla/websocket"
	"time"
)

type Message struct {
	Type       string           `json:"type"`
	Channels   []MessageChannel `json:"channels"`
	ProductIds []string         `json:"product_ids"`
	ProductID  string           `json:"product_id"`
	Message    string           `json:"message"`
	Time       time.Time        `json:"time,string"`
	BestBid    float64          `json:"best_bid,string"`
	BestAsk    float64          `json:"best_ask,string"`
}

type MessageChannel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

const EthBtc string = "ETH-BTC"
const BtcUsd string = "BTC-USD"
const BtcEur string = "BTC-EUR"

func SubscribeMethods(connection *websocket.Conn, message *Message) error {
	var response Message

	if err := connection.WriteJSON(message); err != nil {
		return err
	}

	for {
		if err := connection.ReadJSON(&response); err != nil {
			return err
		}

		if response.Type != "subscriptions" {
			break
		}
	}

	return nil
}

func CreateConnection() (*websocket.Conn, error) {
	var dialer websocket.Dialer
	connection, _, err := dialer.Dial("wss://ws-feed.pro.coinbase.com", nil)

	return connection, err
}
