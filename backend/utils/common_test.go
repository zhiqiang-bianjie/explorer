package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
	"log"
	"testing"
)

func TestParseInt(t *testing.T) {
	_, ok := ParseInt("1")
	assert.True(t, ok)

	_, ok = ParseUint("-1")
	assert.False(t, ok)

	_, ok = ParseInt("sd")
	assert.False(t, ok)
}

const EventNewBlock = "tm.event='NewBlock'"
const EventNewTx = "tm.event='Tx'"

func TestTendermintWs(t *testing.T) {
	var wsUrl = "ws://192.168.150.7:30657/websocket"
	var origin = "http://192.168.150.7:30657/"
	ws, err := websocket.Dial(wsUrl, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var req = RpcReq{
		Jsonrpc: "2.0",
		Method:  "subscribe",
		ID:      "0",
		Params: Params{
			Query: EventNewTx,
		},
	}

	bz, _ := json.Marshal(req)

	if _, err := ws.Write(bz); err != nil {
		log.Fatal(err)
	}

	for {
		var msg = make([]byte, 10*1024)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", msg[:n])
	}

}

type RpcReq struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      string `json:"id"`
	Params  Params `json:"params"`
}

type Params struct {
	Query string `json:"query"`
}

type RpcResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Query string `json:"query"`
		Data  struct {
			Type  string `json:"type"`
			Value struct {
				TxResult struct {
					Height string `json:"height"`
					Index  int    `json:"index"`
					Tx     string `json:"tx"`
					Result struct {
						Log       string `json:"log"`
						GasWanted string `json:"gas_wanted"`
						GasUsed   string `json:"gas_used"`
						Tags      []struct {
							Key   string `json:"key"`
							Value string `json:"value"`
						} `json:"tags"`
					} `json:"result"`
				} `json:"TxResult"`
			} `json:"value"`
		} `json:"data"`
	} `json:"result"`
}
