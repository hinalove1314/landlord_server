package controllers

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Card struct {
	Value       string `json:"value"` // 储存牌的数值，如3,4,5...J,Q,K,A,2
	Suit        string `json:"suit"`  // 储存牌的花色，如Spades, Hearts, Diamonds, Clubs
	ValueWeight int    // 储存牌的数值的排序权重
	SuitWeight  int    // 储存牌的花色的排序权重
}

type Deck struct {
	Cards     []Card `json:"cards"`
	LordCards []Card `json:"lord_cards"`
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
		"Spades":   1,
		"Hearts":   2,
		"Diamonds": 3,
		"Clubs":    4,
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
