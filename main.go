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

func parseQuery(args string) (string, string, []interface{}, error) {
	reader := strings.NewReader(args)
	decoder := transit.NewDecoder(reader)
	value, err := decoder.Decode()
	if err != nil {
		return "", "", nil, err
	}

	argSlice := listToSlice(value.(*list.List))
	addr := argSlice[0].(string)

	switch queryArgs := argSlice[1].(type) {
	case string:
		return addr, queryArgs, make([]interface{}, 0), nil
	case []interface{}:
		return addr, queryArgs[0].(string), queryArgs[1:], nil
	default:
		return "", "", nil, errors.New("unexpected query type, expected a string or a vector")
	}
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
								Name: "connect",
							},
							{
								Name: "query",
							},
						},
					},
				},
			})
	case "invoke":
		addr, _, _, err := parseQuery(message.Args)
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			return
		}

		conn, err := steam.Connect(addr)
		if err != nil {
			babashka.WriteErrorResponse(message, err)
			return
		}

		defer conn.Close()

		switch message.Var {
		case "tommy-mor.go-valve-query/connect":

			duration, err := conn.Ping()

			if err != nil {
				babashka.WriteErrorResponse(message, err)
				return
			}

			res := ExecResult{
				transit.Keyword("ms"): duration.Milliseconds(),
			}

			respond(message, res)

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
