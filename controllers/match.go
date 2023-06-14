package controllers

var matchQueue []*Client

type MatchQueue struct {
	clients []*Client
}

// type Room struct {
// 	ID      int
// 	clients []*Client
// }

func NewRoom(clients []*Client) *Room {
	return &Room{
		clients: clients,
	}
}

func (q *MatchQueue) Enqueue(client *Client) {
	q.clients = append(q.clients, client)
}

func (q *MatchQueue) Dequeue() *Client {
	client := q.clients[0]
	q.clients = q.clients[1:]
	return client
}

func (q *MatchQueue) Size() int {
	return len(q.clients)
}

// func AddToQueue(client *Client) {
// 	matchQueue = append(matchQueue, client)

// 	//打印用户信息
// 	fmt.Println("Username: ", client.UserInfo.Username)
// 	fmt.Println("Password: ", client.UserInfo.UserId)
// }

// func startMatchService(matchQueue *MatchQueue) {
// 	if matchQueue.Size() >= 3 {
// 		clients := make([]*Client, 3)
// 		for i := 0; i < 3; i++ {
// 			clients[i] = matchQueue.Dequeue()
// 		}
// 		room := NewRoom(clients)
// 		// 这里添加将 room 分配到一个用于处理房间游戏逻辑的地方
// 	}
// }
