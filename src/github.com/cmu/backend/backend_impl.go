package backend

import (
	"encoding/json"
	"log"
	"net"
)

type jeopardy struct {
	gameNumber int
	games      map[int]game
	waiting    []player
}

type game struct {
	players []player
	rounds  map[int]struct{} //Keeps track of buzzes
	gameNum int
}

type player struct {
	id         int
	hostport   string
	connection *net.TCPConn
}

func NewJeopardyServer() (Backend, error) {
	return &jeopardy{0, make(map[int]game), nil}, nil
}

/*
* Reads in proposals from the Jeopardy! client and uses type assertions to convert
* the JSON object into a map. From there, we handle the value based on the provided key
 */
func (j *jeopardy) RecvCommit(commitMessage []byte, master bool) error {
	var f interface{}
	err := json.Unmarshal(commitMessage, &f)
	if err != nil {
		log.Println(err)
	}
	m := f.(map[string]interface{})
	for key, value := range m {
		message := value.(map[string]interface{})
		switch key {
		case "Question":
			if master {
				j.handleQA(message, commitMessage)
			}
			break
		case "Buzz":
			j.handleBuzz(message, commitMessage, master)
			break
		case "Join":
			j.handleJoin(message, master)
			break
		case "Answer":
			if master {
				j.handleQA(message, commitMessage)
			}
			break
		default:
			log.Println("Not a supported message type")
		}
	}
	return nil
}

/*
* We let the front-end handle the logic of adding and subtracting from scores
 */
func (j *jeopardy) handleQA(message map[string]interface{}, commitMessage []byte) {
	gameFloat, ok := message["gameID"].(float64)
	if !ok {
		log.Println("Malformed input")
	}
	gameID := int(gameFloat)
	game, ok := j.games[gameID]
	if !ok {
		log.Println("Noexistent game")
	}
	j.echoAll(game.players, commitMessage)
}

/*
* The backend must do the work of determining the first buss for each turn in the given game.
* All buzzes for a round are ignored after the first game.
 */
func (j *jeopardy) handleBuzz(message map[string]interface{}, commitMessage []byte, master bool) {
	gameFloat, ok := message["gameID"].(float64)
	if !ok {
		log.Println("Malformed input")
	}
	gameID := int(gameFloat)
	round, ok := message["turn"].(float64)
	if !ok {
		log.Println("Malformed input")
	}
	roundID := int(round)
	game, ok := j.games[gameID]
	if !ok {
		log.Println("Noexistent game")
	}
	_, ok = game.rounds[roundID]
	if !ok {
		game.rounds[roundID] = struct{}{}
		if master {
			j.echoAll(game.players, commitMessage)
		}
	}
}

/*
* When a client asks to join a game of Jeopardy!, the game determines if there are
* two other waiting players and places them in a game if this is the case. Otherwise,
* they are queued and wait until there are three players.
 */
func (j *jeopardy) handleJoin(message map[string]interface{}, master bool) {
	hostport := message["hostport"].(string)
	if j.waiting == nil {
		j.waiting = make([]player, 0, 3) //3 players per game
	}
	j.waiting = append(j.waiting, player{len(j.waiting), hostport, nil})
	if len(j.waiting) == 3 {
		newGame := game{j.waiting, make(map[int]struct{}), j.gameNumber}
		j.games[j.gameNumber] = newGame
		if master {
			j.sendGame(newGame)
		}
		j.waiting = nil
		j.gameNumber++
	}
}

/*
* Once three players are placed into a game, they are sent a JSON object
* that includes both their gameID and their playerID
 */
func (j *jeopardy) sendGame(newGame game) {
	reply := make(map[string]interface{})
	reply["gameID"] = newGame.gameNum
	for index, player := range newGame.players {
		addr, _ := net.ResolveTCPAddr("tcp", player.hostport)
		connection, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			log.Println(err)
		}
		player.connection = connection
		newGame.players[index] = player
		reply["playerID"] = player.id
		data := make(map[string]interface{})
		data["JoinRep"] = reply
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
		n, err := connection.Write(bytes)
		if len(bytes) != n {
			log.Println("Failed to write the whole message")
		} else if err != nil {
			log.Println(err)
		}
	}
}

/*
* The provided message is sent to each player in the provided array of players
 */
func (j *jeopardy) echoAll(players []player, message []byte) {
	for _, player := range players {
		if player.connection == nil {
			addr, _ := net.ResolveTCPAddr("tcp", player.hostport)
			connection, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				log.Println(err)
			}
			player.connection = connection
		}
		n, err := player.connection.Write(message)
		if err != nil {
			log.Println(err)
		} else if n != len(message) {
			log.Println("Failed to write the whole message")
		}
	}
}
