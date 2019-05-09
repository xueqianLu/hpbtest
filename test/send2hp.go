package main

import (
	"fmt"
	"flag"
	"time"
	"strconv"
	"math/rand"
	"encoding/json"
	"net/http"
	"strings"
)

type Tx struct {
    From  string   `json:"from"`
    To    string   `json:"to"`
    Value string   `json:"value"`
}

type Data struct {
    Txs      []Tx   `json:"params"`
    Id       int    `json:"id"`
    Jsonrpc  string `json:"jsonrpc"`
    Method   string `json:"method"`
}


const (
    synurl1 = "http://127.0.0.1:28581"
    synurl2 = "http://127.0.0.1:28582"
    synurl3 = "http://127.0.0.1:28583"
    synurl4 = "http://127.0.0.1:28584"
    synurl5 = "http://127.0.0.1:28585"
    synurl6 = "http://127.0.0.1:28586"
    synaddr1 = "0x75c9feb3a21b88b42f7f1455041f1916af411e5a"
    synaddr2 = "0x485e6ae66b2b916a1c4258bc82c6e8e33c3cc05a"
    synaddr3 = "0xa5fccef6da8ba4fe4e9129b9ae49380c2cfa609a"
    synaddr4 = "0x1617bd386d2fb34280db46aa81747feb3e537a0a"
    synaddr5 = "0x2ca0a43bdb7f30af8e7d3a7d902a1e1737470c5a"
    synaddr6 = "0xa6bab593741805615dada5e0f2091febf930221a"

	hpurl1 = "http://127.0.0.1:28511"
	hpurl2 = "http://127.0.0.1:28512"
	hpurl3 = "http://127.0.0.1:28513"
	hpurl4 = "http://127.0.0.1:28514"
    hpaddr1 = "0xa0908bc7f19d0e1fce0c9198d7ca72d105088600"
    hpaddr2 = "0x9d3c60ee03bcd028ee82c728dbeb087b4aacffc8"
    hpaddr3 = "0x170633e5cb45bb2b1ed61e4ea7db2f5dab3a6f78"
    hpaddr4 = "0xc65ff5de55c3ca33efbf10cde0fc6ed23cd0e2f8"
    hpaddr5 = "0xd5729a0eca5e674da48692d32c86c97c24630758"
    hpaddr6 = "0x5096d44cc486d20d758b3e37322a4f940b8eb158"
    hpaddr7 = "0x766f9a78575114bd22d6c5db1c8eb89efa5e57b8"
    hpaddr8 = "0x4817c85a565a92aef9782f4705461bb06a96fff8"
    hpaddr9 = "0x20e9005b5c6cd4788b9a66e62e219518dfb91418"

    preaddr1 = "0x3768feb27566f3725cb175209cf9e8fe2918a099"
    preaddr2 = "0x9a8acd088d91b134abe30c4631829f0bd40cd709"
    preaddr3 = "0xe7b435c93c3bc62202caceda438f5a62f02beb69"
    preaddr4 = "0x96d957a8b412a671674662354acdd87ad3026999"
    preaddr5 = "0x757e17786cd90c63f61f694424db7254bdaef369"
    preaddr6 = "0x4c95f3bdc9ef0bb19656e050514d78ba9299ed79"
    preaddr7 = "0xfd9f49db44458b1fa4cfa4df83303ce536e6a349"
    preaddr8 = "0x10b3351f018fe43fc8a36921856d219637a50639"
    preaddr9 = "0x2df2cab7e57ff4d10a13de9306e49b7a054d2129"

)
var URL     =[]string{synurl1,  synurl2, synurl3}
var FROMADR =[]string{synaddr1, synaddr2,synaddr3}

//var URL     =[]string{hpurl1,  hpurl2, hpurl3,  hpurl4}
//var FROMADR =[]string{hpaddr1, hpaddr2,hpaddr3, hpaddr4}

func main() {

	index:= flag.Int("t", 0,   "target syn node index 0")
	count:= flag.Int("c", 100, "count to send 100")
	sleep:= flag.Int("s", 200, "sleep betown intervel millisecond 200 ")
	bodycount:= flag.Int("b",1000, "tx count by one body")
	flag.Parse()

	cli := http.Client{}
	start := time.Now()
	txs := 0.0

	if *index <0 || *index >= len(URL) {
		panic("target syn is out of rang")
	}
	url  := URL[*index]
	from := FROMADR[*index]
	toadr :=[]string{preaddr1,preaddr2,preaddr3,preaddr4,preaddr5,preaddr6,preaddr7,preaddr8}

	rand.Seed(time.Now().UnixNano())

	for i:=0; i< *count; i++{
		//stb  := time.Now()
		body := buildBodyString(from, toadr[rand.Intn(len(toadr))], *bodycount)
		//fmt.Println("build body elapse (ms)",time.Now().Sub(stb).Nanoseconds()/1000000)
		stb  := time.Now()
		httpSendHPB(url, body, cli)
		//fmt.Println("send body elapse (ms)",time.Now().Sub(stb).Nanoseconds()/1000000,"size",len(body))
		time.Sleep(time.Millisecond*time.Duration(*sleep))
		txs = txs + float64(*bodycount)
		t :=float64(time.Now().Sub(start).Seconds())
		if t < 1 {continue}
		fmt.Println("total txs:",txs, " rate(tx/s):",txs/t, " elapse(ms)",time.Now().Sub(stb).Nanoseconds()/1000000," size",len(body))
	}

	elapse :=time.Now().Sub(start).Seconds()
	fmt.Println("elapse(s):",elapse)
}



func buildBodyString(from string, to string, times int) string {
	hpb := 100000000000000000 // 0.1 HPB
	//if from == "0x9d3c60ee03bcd028ee82c728dbeb087b4aacffc8" {
	//	hpb = 200000000000000000 // 0.2 HPB
	//}
	//if from == "0x170633e5cb45bb2b1ed61e4ea7db2f5dab3a6f78" {
	//	hpb = 300000000000000000 // 0.3 HPB
	//}
	//if from == "0xc65ff5de55c3ca33efbf10cde0fc6ed23cd0e2f8" {
	//	hpb = 400000000000000000 // 0.4 HPB
	//}

	datas := make([]interface{},0)

	for i:=0; i < times; i++{
		ss := Data{Method:"hpb_sendTransaction",Jsonrpc:"2.0",Id:67}
		_tx := Tx{From:from,To:to}
		_tx.Value = "0x"+ strconv.FormatInt(int64(hpb), 16)
		ss.Txs = append(ss.Txs,_tx)
		datas = append(datas, ss)
	}


	body,_:= json.Marshal(datas)

	return string(body)
}

func httpSendHPB(url string, body string, client http.Client)  {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http Send HPB Error:", err)
	}
	defer resp.Body.Close()
}
