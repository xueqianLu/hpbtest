package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AccountNonce struct {
	Addr  string
	Nonce uint64
}

/*
[{"jsonrpc":"2.0","id":68,"result":{"raw":"0xf86780850430e2340083015f90949d882e29357b8fda9ad232760ab8ea763c7484908203e88026a008e1f63b9d2ceb607216d9b4f8ff7662c797ea5da84c59611519d7bf71fafc5ba07dd7f052492ed54194881accda1c22a104fd8191d156fda0cf3f9546277a2d98","tx":{"nonce":"0x0","gasPrice":"0x430e23400","gas":"0x15f90","to":"0x9d882e29357b8fda9ad232760ab8ea763c748490","value":"0x3e8","input":"0x","exdata":{"txversion":0,"txtype":0,"vmversion":0,"txflag":0,"reserve":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},"v":"0x26","r":"0x8e1f63b9d2ceb607216d9b4f8ff7662c797ea5da84c59611519d7bf71fafc5b","s":"0x7dd7f052492ed54194881accda1c22a104fd8191d156fda0cf3f9546277a2d98","hash":"0x54e20f016901ca4e074e75ff17ad2d7795f08bc444c7af60062748a9056e9a2e"}}}]
*/

type TxInfo struct {
	Nonce string `json:"nonce"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type SignTxResult struct {
	Raw    string `json:"raw"`
	Txinfo TxInfo `json:"tx"`
}

func (result SignTxResult) InValid() bool{
	return result.Raw == "" 
}

func (result SignTxResult) String() string {
	var s = "raw:" + result.Raw + "\r\n"
	s += "nonce:" + result.Txinfo.Nonce + "\r\n"
	return s
}

type TxArgs struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Nonce string `json:"nonce,omitempty"`
}

type SendSignData struct {
	Txs     []TxArgs `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
}

type RespSignData struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      int          `json:"id"`
	Signtx  SignTxResult `json:"result"`
}

type SendRawData struct {
	Txs     []string `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
}

type RespRawData struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Txhash  string `json:"result"`
}

/*
[{"jsonrpc":"2.0","id":68,"error":{"code":-32000,"message":"authentication needed: password or unlock"}}]
*/
type ErrMsg struct {
	Msg	string `json:"message"`
}
type RespErrData struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Err     ErrMsg `json:"error"`
}

var (
	fromAddr [251]string
	toAddr   [251]string
	fromInfo []AccountNonce
	toInfo   []AccountNonce
)

var (
	sendurl  = [2]string{"http://192.168.1.26:18545","http://192.168.1.28:18545"}
	url = "http://192.168.1.26:18545" // sign url
)

func buildSignString(from string, to string, value int, nonce uint64) string {
	datas := make([]interface{}, 0)

	ss := SendSignData{Method: "hpb_signTransaction", Jsonrpc: "2.0", Id: 68}
	_tx := TxArgs{From: from, To: to}
	_tx.Value = "0x" + strconv.FormatInt(int64(value), 16)
	if nonce != 0 {
		_tx.Nonce = "0x" + strconv.FormatUint(nonce, 16)
	}

	ss.Txs = append(ss.Txs, _tx)
	datas = append(datas, ss)

	body, _ := json.Marshal(datas)
	return string(body)
}

/*
[{"jsonrpc":"2.0","id":68,"result":{"raw":"0xf86780850430e2340083015f90949d882e29357b8fda9ad232760ab8ea763c7484908203e88026a008e1f63b9d2ceb607216d9b4f8ff7662c797ea5da84c59611519d7bf71fafc5ba07dd7f052492ed54194881accda1c22a104fd8191d156fda0cf3f9546277a2d98","tx":{"nonce":"0x0","gasPrice":"0x430e23400","gas":"0x15f90","to":"0x9d882e29357b8fda9ad232760ab8ea763c748490","value":"0x3e8","input":"0x","exdata":{"txversion":0,"txtype":0,"vmversion":0,"txflag":0,"reserve":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},"v":"0x26","r":"0x8e1f63b9d2ceb607216d9b4f8ff7662c797ea5da84c59611519d7bf71fafc5b","s":"0x7dd7f052492ed54194881accda1c22a104fd8191d156fda0cf3f9546277a2d98","hash":"0x54e20f016901ca4e074e75ff17ad2d7795f08bc444c7af60062748a9056e9a2e"}}}]
*/
func doSignTx(url string, body string, client *http.Client) (*SignTxResult, error) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	//fmt.Println("doSignTx >>>>>>:",body)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http Send HPB Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("doSignTx <<<<<<:", string(content))

	var results = make([]RespSignData, 0)
	//err = json.NewDecoder(resp.Body).Decode(&results)
	err = json.Unmarshal(content, &results)
	if err != nil {
		return nil, err
	}
	return &results[0].Signtx, nil
}

func buildSendRawString(raw string) string {
	datas := make([]interface{}, 0)
	ss := SendRawData{Method: "hpb_sendRawTransaction", Jsonrpc: "2.0", Id: 69, Txs: make([]string, 0)}
	ss.Txs = append(ss.Txs, raw)
	datas = append(datas, ss)

	body, _ := json.Marshal(datas)
	return string(body)
}

func doSendRaw(url string, body string, client *http.Client) ([]string, error) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	//fmt.Println("doSendRaw >>>>>>:",string(body))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http Send HPB Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("doSignTx <<<<<<:", string(content))

	var results = make([]RespRawData, 0)
	var txhashs = make([]string, 0)
	//err = json.NewDecoder(resp.Body).Decode(&results)
	err = json.Unmarshal(content, &results)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		txhashs = append(txhashs, result.Txhash)
	}

	return txhashs, nil
}

func parseUint(str string) uint64 {
	//println("str=",str)
	var nstr string
	if strings.Compare(str[0:2], "0x") == 0 {
		if len(str) == 2 {
			return 0
		} else {
			nstr = str[2:len(str)]
		}
	} else {
		nstr = str
	}
	s, _ := strconv.ParseUint(nstr, 16, 64)
	return s
}

type Queue struct {
	queue    chan string
}

func main() {
	var (
		txRaw = make([]string, 0)
	)

	runtime.GOMAXPROCS(runtime.NumCPU())
	count := flag.Int("c", 1, "count to send 100")
	pool  := flag.Int("p", 100, "routine pool number")
	us    := flag.Int("s", 0, "url index start")
	ue    := flag.Int("e", 0, "url index end")
	usurl := make([]string, 0)

	flag.Parse()
	if *us >= 0 && *us <= *ue && *ue < len(sendurl) {
		for i := *us ; i <= *ue; i++ {
			usurl = append(usurl, sendurl[i])
			println("use send url ", sendurl[i])
		}
	} else {
		fmt.Printf("us(%d) and ue(%d) invalid.\n", *us, *ue)
	}
	cli := http.Client{}
	for cnt := 0; cnt < *count; cnt++ {
		for i := 0; i < len(fromInfo); i++ {
			body := buildSignString(fromInfo[i].Addr, toInfo[i].Addr, 1001, fromInfo[i].Nonce)
			signedTx, err := doSignTx(url, body, &cli)

			if err != nil || signedTx.InValid() {
				fmt.Println("doSignTx failed, cnt =", cnt, "i =", i)
				break
			}
			fromInfo[i].Nonce = parseUint(signedTx.Txinfo.Nonce) + 1
			txRaw = append(txRaw, signedTx.Raw)
		}
	}

	var wg sync.WaitGroup 
	st := time.Now().UnixNano()
	l := len(txRaw)
	qpool := make([]*Queue, *pool)
	for i,_ := range qpool {
		qpool[i] = &Queue{queue:make(chan string, 100)}
		go func (c chan string) {
			surl := usurl[i%len(usurl)]
			for {
				txraw := <- c

				body := buildSendRawString(txraw)
				hashs, err := doSendRaw(surl, body, &cli)

				if err != nil {
					fmt.Println("SendRaw failed:", err)
				} else {
					for _, hash := range hashs {
						fmt.Println("SendTxHash:", hash)
					}
				}
				wg.Done()
			}
		}(qpool[i].queue)
	}
	for i, txraw := range txRaw {
		// hpb_sendRawTransaction 只支持一次发送一笔交易.
		wg.Add(1)
		q := qpool[i%len(qpool)]
		q.queue <- txraw
	}
	wg.Wait()
	se := time.Now().UnixNano()
	fmt.Println("SendRawTransaction cost time ", (se - st)/1000/1000, "ms, total ", l)
}
