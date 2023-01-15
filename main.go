package main

import (
	"encoding/hex"
	"fmt"
	tool "github.com/eth-collision/eth-collision-tool"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

var totalFile = "total.txt"
var accountsFile = "accounts.txt"
var speedFile = "speed.txt"

var locker = sync.Mutex{}

// second
var rollupTime time.Duration = 1 * 60 * 60
var submitTime time.Duration = 1 * 60

func main() {
	msg := make(chan *big.Int)
	for i := 0; i < 20; i++ {
		go generateAccountJob(msg)
	}
	totalStr := readFile(totalFile)
	n := new(big.Int)
	total, ok := n.SetString(totalStr, 10)
	if !ok {
		total = big.NewInt(0)
	}
	lastTotal := total
	tick := time.Tick(rollupTime * time.Second)
	for {
		select {
		case <-tick:
			speed := new(big.Int).Sub(total, lastTotal)
			lastTotal = total
			addresses, err := fileCountLine(accountsFile)
			if err != nil {
				log.Println(err)
			}
			totalStr := tool.FormatBigInt(*total)
			speedStr := tool.FormatBigInt(*speed)
			addrsStr := tool.FormatInt(int64(addresses))
			text := fmt.Sprintf(""+
				"[ETH Collision Find Address]\n"+
				"Total: %s\n"+
				"Speed: %s\n"+
				"Addrs: %s\n",
				totalStr, speedStr, addrsStr)
			appendFile(speedFile, text)
			tool.SendMsgText(text)
		case count := <-msg:
			total = bigIntAddMutex(total, count)
			writeFile(totalFile, total.String())
		}
	}
}

func bigIntAddMutex(a, b *big.Int) *big.Int {
	locker.Lock()
	defer locker.Unlock()
	c := new(big.Int)
	c.Add(a, b)
	return c
}

func generateAccountJob(msg chan *big.Int) {
	count := big.NewInt(0)
	tick := time.Tick(submitTime * time.Second)
	for {
		select {
		case <-tick:
			msg <- count
			count = big.NewInt(0)
		default:
			generateAccount()
			count = count.Add(count, big.NewInt(1))
		}
	}
}

func generateAccount() {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Println(err)
	}
	privateKey := hex.EncodeToString(key.D.Bytes())
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	handleAccount(privateKey, address)
}

func checkAddress(address string) bool {
	if strings.HasPrefix(address, "0x0000000") {
		return true
	}
	return false
}

func handleAccount(privateKey string, address string) {
	if checkAddress(address) {
		log.Println("Found: ", privateKey, address)
		text := fmt.Sprintf("%s,%s\n", privateKey, address)
		appendFile(accountsFile, text)
	}
}
