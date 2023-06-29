package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/astaxie/beego/logs"
)

var roominfo RoomInfo
var playcard PlayCard
var playerInfo PlayerInfo
var UnPlayCardNum = 0

type Card struct {
	Value       string `json:"value"`       // 储存牌的数值，如3,4,5...J,Q,K,A,2
	Suit        string `json:"suit"`        // 储存牌的花色，如Spades, Hearts, Diamonds, Clubs
	ValueWeight int    `json:"ValueWeight"` // 储存牌的数值的排序权重
	SuitWeight  int    `json:"SuitWeight"`  // 储存牌的花色的排序权重
}

type Deck struct {
	Cards     []Card `json:"cards"`
	LordCards []Card `json:"lord_cards"`
}

type PlayCard struct {
	Cards []Card `json:"cards"`
}

type PlayerInfo struct {
	SeatNum   int    `json:"seatNum"`
	CardNum   int    `json:"cardNum"`
	PlayCards []Card `json:"playCards"`
}

// 接收出牌的数据，并把数据发送给发送客户端的下一个客户端
func receivePlayCards(msg []byte, room *Room, client *Client) {
	fmt.Println("client seatnum =", client.SeatNum)
	err := json.Unmarshal(msg, &playerInfo)
	if err != nil {
		logs.Error("Error parsing JSON:", err)
		return
	}

	fmt.Println("cardinfo=", playerInfo.PlayCards[0].Value)
	if len(playerInfo.PlayCards) == 0 {
		fmt.Println("PlayCards is empty or null")
	}

	if client.SeatNum < 3 {
		sendResponse(playerInfo, 42, room.clients[client.SeatNum]) //把出牌数据发送给下一个客户端(因为client从0开始，seatNum从一开始，所以这里相当于+1了)
	} else {
		sendResponse(playerInfo, 42, room.clients[0]) //把出牌数据发送给第一个客户端
	}
}

// 接收到不出的消息
func UnreceivePlayCards(msg []byte, room *Room, client *Client) {
	fmt.Println("client seatnum =", client.SeatNum)
	err := json.Unmarshal(msg, &playcard)
	if err != nil {
		logs.Error("Error parsing JSON:", err)
		return
	}

	fmt.Println("cardinfo=", playcard.Cards[0].Value)

	UnPlayCardNum++

	if client.SeatNum < 3 {
		if UnPlayCardNum == 2 {
			sendResponse(playcard, 44, room.clients[client.SeatNum]) //发送不出的消息，下一个客户端必须要出牌
			UnPlayCardNum = 0
		}
		sendResponse(playcard, 42, room.clients[client.SeatNum]) //把出牌数据发送给下一个客户端(因为client从0开始，seatNum从一开始，所以这里相当于+1了)
	} else {
		if UnPlayCardNum == 2 {
			sendResponse(playcard, 44, room.clients[0]) //把出牌数据发送给第一个客户端
			UnPlayCardNum = 0
		}
		sendResponse(playcard, 42, room.clients[0]) //把出牌数据发送给第一个客户端
	}
}

func createAndShuffleDeck() Deck {
	fmt.Println("createAndShuffleDeck")
	// 初始化牌库
	values := map[string]int{
		"3":            1,
		"4":            2,
		"5":            3,
		"6":            4,
		"7":            5,
		"8":            6,
		"9":            7,
		"10":           8,
		"J":            9,
		"Q":            10,
		"K":            11,
		"A":            12,
		"2":            13,
		"little_joker": 14,
		"big_joker":    15,
	}
	suits := map[string]int{
		"Spades":   1, //黑桃
		"Hearts":   2, //红心
		"Diamonds": 3, //方块
		"Clubs":    4, //梅花
		"Joker":    5,
	}
	deck := make([]Card, 54)

	index := 0
	for value, valueWeight := range values {
		if value == "little_joker" || value == "big_joker" {
			deck[index] = Card{Value: value, Suit: "Joker", ValueWeight: valueWeight, SuitWeight: suits["Joker"]}
			index++
		} else {
			for suit, suitWeight := range suits {
				if suit != "Joker" {
					deck[index] = Card{Value: value, Suit: suit, ValueWeight: valueWeight, SuitWeight: suitWeight}
					index++
				}
			}
		}
	}

	// 洗牌
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return Deck{
		Cards:     deck[:51],
		LordCards: deck[51:],
	}
}

func DealCards(room *Room, client []*Client) {
	// 创建和洗牌的逻辑
	deck := createAndShuffleDeck()

	// 分牌的逻辑，这里简单地分配每个玩家17张牌
	for i := range client {
		client[i].Hand = deck.Cards[i*17 : (i+1)*17]

		// 对每个玩家的手牌进行排序
		sort.Slice(client[i].Hand, func(j, k int) bool {
			if client[i].Hand[j].ValueWeight == client[i].Hand[k].ValueWeight {
				return client[i].Hand[j].SuitWeight < client[i].Hand[k].SuitWeight
			}
			return client[i].Hand[j].ValueWeight < client[i].Hand[k].ValueWeight
		})
	}

	// 存储地主牌
	room.LordCards = deck.LordCards

	// 打印地主牌
	fmt.Println("Lord Cards: ", room.LordCards)

	//剩余的三张牌作为地主牌
	//client[0].LordCards = deck.LordCards
}

// 决定谁是地主
func DealLord(msg []byte, room *Room, client *Client) {
	fmt.Println("At the start of DealLord, room = ", room)

	err := json.Unmarshal(msg, &roominfo)
	if err != nil {
		logs.Error("Error parsing JSON:", err)
		return
	}

	fmt.Println("roomID=", roominfo.ID)
	fmt.Println("isCalled=", roominfo.IsCalled)

	if roominfo.IsCalled {
		client.LordPoint += 3
	}
	fmt.Println("client.LordPoint=", client.LordPoint)

	room.CallLordCount += 1 // 增加叫地主的客户数量,用来判断有几名用户叫地主了
	// 假设房间里最初没有客户，LordPoint 最高的客户为 nil
	if room.CallLordCount == len(room.clients) {
		var maxLordPointClient *Client = nil

		if room.clients == nil {
			logs.Error("Room clients list is nil")
			return
		}

		// 遍历房间里的所有客户
		for _, c := range room.clients {
			// 如果这是第一个客户，或者他的 LordPoint 比当前最高的还要高
			if maxLordPointClient == nil || c.LordPoint > maxLordPointClient.LordPoint {
				// 则更新最高 LordPoint 的客户
				maxLordPointClient = c
			}
		}

		// 打印最高 LordPoint 的客户信息
		if maxLordPointClient != nil {
			fmt.Println("The client with the highest LordPoint is:", maxLordPointClient.UserInfo.Username)
		} else {
			fmt.Println("There are no clients in the room.")
		}

		room.sendInformation(maxLordPointClient.SeatNum, 36) //把成为地主的seatNum发送给客户端
		room.sendInformation(room.LordCards, 38)             //发送地主牌消息
	}
}
