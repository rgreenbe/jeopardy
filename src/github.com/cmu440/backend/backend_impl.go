package backend

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
	return jeopardy{0, new(map[int]struct{}), nil}
}

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
			j.handleQA(message, commitMessage)
		case "Buzz":
			j.handleBuzz(message, commitMessage)
		case "Join":
			j.handleJoin(message)
		case "Answer":
			j.handleQA(message, commitMessage)
		default:
			log.Println("Not a supported message type")
		}
	}
}

func (j *jeopardy) handleQA(message map[string]interface{}, commitMessage []byte) {
	gameID, ok := int(message["gameID"])
	if !ok {
		log.Println("Malformed input")
	}
	game, ok := j.games[gameID]
	if !ok {
		log.Println("Noexistent game")
	}
	j.echoAll(game.players, commitMessage)
}

func (j *jeopardy) handleBuzz(message map[string]interface{}, commitMessage []byte) {
	gameID, ok := int(message["gameID"])
	if !ok {
		log.Println("Malformed input")
	}
	round, ok := int(message["turn"])
	if !ok {
		log.Println("Malformed input")
	}
	game, ok := j.games[gameID]
	if !ok {
		log.Println("Noexistent game")
	}
	_, ok = game.rounds[round]
	if !ok {
		game.rounds[round] = struct{}{}
		j.echoAll(game.players, commitMessage)
	}
}

func (j *jeopardy) handleJoin(message map[string]interface{}) {
	hostport := string(message["hostport"])
	if j.waiting == nil {
		j.waiting = make([]player, 0, 3) //3 players per game
	}
	j.waiting = append(j.waiting, player{len(j.waiting), hostport, nil})
	if len(j.waiting) == 3 {
		newGame := game{j.waiting, make(map[int]struct{}), j.gameNumber}
		j.games[j.gameNumber] = newGame
		j.sendGame(newGame)
		j.waiting = nil
		j.gameNumber++
	}
}

func (j *jeopardy) sendGame(newGame game) {
	reply := make(map[string]interface{})
	reply["gameID"] = newGame.gameNum
	for _, player := range newGame.players {
		connection, err := net.DialTCP("tcp", nil, player.hostport)
		if err != nil {
			log.Println(err)
		}
		player.connection = connection
		reply["playerID"] = player.id
		data := make(map[string]interface{})
		data["JoinRep"] = reply
		bytes, err := json.Marshal(data)
		n, err := connection.Write(bytes)
		if len(bytes) != n {
			log.Println("Failed to write the whole message")
		}
	}
}

func (j *jeopardy) echoAll(players []player, message []byte) {
	for _, player := range players {
		n, err := player.connection.Write(message)
		if err != nil {
			log.Println(err)
		} else if n != len(message) {
			log.Println("Failed to write the whole message")
		}
	}
}
