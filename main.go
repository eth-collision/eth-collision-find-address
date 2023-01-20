package main

import (
	"encoding/hex"
	"fmt"
	tool "github.com/eth-collision/eth-collision-tool"
	"log"
	"math/big"
	"regexp"
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
	totalStr := tool.ReadFile(totalFile)
	n := new(big.Int)
	total, ok := n.SetString(totalStr, 10)
	if !ok {
		total = big.NewInt(0)
	}
	lastTotal := total
	ticker, callback := tool.NewProxyTicker(rollupTime * time.Second)
	go callback()
	for {
		select {
		case <-ticker:
			speed := new(big.Int).Sub(total, lastTotal)
			lastTotal = total
			addresses, err := tool.FileCountLine(accountsFile)
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
			tool.AppendFile(speedFile, text)
			tool.SendMsgText(text)
		case count := <-msg:
			total = bigIntAddMutex(total, count)
			tool.WriteFile(totalFile, total.String())
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

var re = regexp.MustCompile(`0x00000000|0x11111111|0x22222222|0x33333333|0x44444444|0x55555555|0x66666666|0x77777777|0x88888888|0x99999999|0xaaaaaaaa|0xbbbbbbbb|0xcccccccc|0xdddddddd|0xeeeeeeee|0xffffffff`)

func checkAddress(address string) bool {
	if re.MatchString(address) {
		return true
	}
	return false
}

func handleAccount(privateKey string, address string) {
	if checkAddress(address) {
		log.Println("Found: ", privateKey, address)
		text := fmt.Sprintf("%s,%s\n", privateKey, address)
		tool.AppendFile(accountsFile, text)
		tool.SendMsgText(text)
	}
}
