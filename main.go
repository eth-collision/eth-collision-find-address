package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	msg := make(chan int)
	for i := 0; i < 4; i++ {
		go generateAccountJob(msg)
	}
	total := 0
	lastTotal := 0
	tick := time.Tick(1 * time.Hour)
	for {
		select {
		case <-tick:
			speed := total - lastTotal
			lastTotal = total
			addresses, err := fileCountLine("accounts.txt")
			if err != nil {
				log.Println(err)
			}
			text := fmt.Sprintf("Total: %d, Speed: %d, Addresses: %d\n", total, speed, addresses)
			appendFile("speed.txt", text)
			sendMsgText(text)
		case count := <-msg:
			total += count
		}
	}
}

func generateAccountJob(msg chan int) {
	count := 0
	tick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-tick:
			msg <- count
			count = 0
		default:
			generateAccount()
			count++
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
	if strings.HasPrefix(address, "0x8888") && strings.HasSuffix(address, "8888") {
		return true
	}
	return false
}

func handleAccount(privateKey string, address string) {
	if checkAddress(address) {
		log.Println("Found: ", privateKey, address)
		text := fmt.Sprintf("%s,%s\n", privateKey, address)
		appendFile("accounts.txt", text)
	}
}
