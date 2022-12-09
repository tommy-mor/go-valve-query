package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"container/list"

	"github.com/babashka/pod-babashka-go-sqlite3/babashka"
	"github.com/russolsen/transit"

	"github.com/NewPage-Community/go-source-server-query"
)

func debug(v interface{}) {
	fmt.Fprintf(os.Stderr, "debug: %+q\n", v)
}


type ExecResult = map[transit.Keyword]int64

func listToSlice(l *list.List) []interface{} {
	slice := make([]interface{}, l.Len())
	cnt := 0
	for e := l.Front(); e != nil; e = e.Next() {
		slice[cnt] = e.Value
		cnt++
	}

	return slice
}

func parseQuery(args string) (string, []interface{}, error) {
	reader := strings.NewReader(args)
	decoder := transit.NewDecoder(reader)
	value, err := decoder.Decode()
	if err != nil {
		return "", nil, err
	}

	argSlice := listToSlice(value.(*list.List))
	addr := argSlice[0].(string)
	rest := argSlice[1:]
	return addr, rest, nil

	//if len(argSlice) == 1 {
	//addr := argSlice[0].(string)
	//return addr, "", "", nil
	//} else if len(argSlice) == 3 {
	//addr := argSlice[0].(string)
	//pw := argSlice[1].(string)
	//cmd := argSlice[2].(string)
	//return addr, pw, cmd, nil
	//} else {
	//return "", "", "", errors.New("invalid args")
	//}
		
}

func respond(message *babashka.Message, response interface{}) {
	buf := bytes.NewBufferString("")
	encoder := transit.NewEncoder(buf, false)

	errors.New("unexpected query type, expected a string or a vector")

	if err := encoder.Encode(response); err != nil {
		babashka.WriteErrorResponse(message, err)
	} else {
		babashka.WriteInvokeResponse(message, string(buf.String()))
	}
}

func encodeResult(res *babashka.Message) (ExecResult, error) {
	return nil, errors.New("woops")
}



type InfoResult = map[transit.Keyword]interface{}


func processMessage(message *babashka.Message) {
	switch message.Op {
	case "describe":
		babashka.WriteDescribeResponse(
			&babashka.DescribeResponse{
				Format: "transit+json",
				Namespaces: []babashka.Namespace{
					{
						Name: "tommy-mor.go-valve-query",
						Vars: []babashka.Var{
							{
								Name: "ping",
							},
							{
								Name: "info",
							},
							{
								Name: "players",
							},
							{
								Name: "rcon",
							},
						},
					},
				},
			})
	case "invoke":
		addr, args, err := parseQuery(message.Args)
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			return
		}

		var options *steam.ConnectOptions = nil
		if len(args) != 0 {
			options = &steam.ConnectOptions{RCONPassword: args[0].(string)}
		}  else {
			options = &steam.ConnectOptions{}
		}

		conn, err := steam.Connect(addr, options)
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			return
		}

		defer conn.Close()

		switch message.Var {
		case "tommy-mor.go-valve-query/info":

			info, err := conn.Info()

			if err != nil {
				babashka.WriteErrorResponse(message, err)
				return
			}

			res := InfoResult{
				transit.Keyword("protocol"):          info.Protocol,
				transit.Keyword("name"):              info.Name,
				transit.Keyword("map"):               info.Map,
				transit.Keyword("folder"):            info.Folder,
				transit.Keyword("game"):              info.Game,
				transit.Keyword("id"):                info.ID,
				transit.Keyword("players"):           info.Players,
				transit.Keyword("max-players"):       info.MaxPlayers,
				transit.Keyword("bots"):              info.Bots,
				transit.Keyword("server-type"):       info.ServerType,
				transit.Keyword("environment"):       info.Environment,
				transit.Keyword("visibility"):        info.Visibility,
				transit.Keyword("vac"):               info.VAC,
				transit.Keyword("version"):           info.Version,
				transit.Keyword("port"):              info.Port,
				transit.Keyword("steam-id"):          info.SteamID,
				transit.Keyword("sourcetv-port"):     info.SourceTVPort,
				transit.Keyword("sourcetv-name"):     info.SourceTVName,
				transit.Keyword("keywords"):          info.Keywords,
				transit.Keyword("game-id"):           info.GameID,
			}

			respond(message, res)

		case "tommy-mor.go-valve-query/players":

			players, err := conn.PlayersInfo()

			if err != nil {
				babashka.WriteErrorResponse(message, err)
				return
			}

			res := make([]InfoResult, len(players.Players))
			for playerIdx, player := range players.Players {
				res[playerIdx] = InfoResult{
					transit.Keyword("index"): player.Index,
					transit.Keyword("name"):  player.Name,
					transit.Keyword("score"): player.Score,
					transit.Keyword("duration"):  player.Duration,
				}
			}

			respond(message, res)

		case "tommy-mor.go-valve-query/ping":

			duration, err := conn.Ping()

			if err != nil {
				babashka.WriteErrorResponse(message, err)
				return
			}

			res := ExecResult{
				transit.Keyword("ms"): duration.Milliseconds(),
			}

			respond(message, res)

		case "tommy-mor.go-valve-query/rcon":
			if(len(args) != 2) {
				babashka.WriteErrorResponse(message, errors.New("rcon requires 2 arguments"))
				return
			}

			resp, err := conn.Send(args[1].(string))

			if err != nil {
				babashka.WriteErrorResponse(message, err)
				return
			}

			respond(message, resp)


		default:
			babashka.WriteErrorResponse(message, fmt.Errorf("Unknown var %s", message.Var))
		}
	default:
		babashka.WriteErrorResponse(message, fmt.Errorf("Unknown op %s", message.Op))
	}
}

func main() {
	for {
		message, err := babashka.ReadMessage()
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			continue
		}

		processMessage(message)
	}
}
