package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/d4c5d1e0/webtoons"
	"github.com/d4c5d1e0/webtoons/internal/helpers"
	"github.com/d4c5d1e0/webtoons/mail"
)

const (
	ThreadNumber = 1
	// Change the proxy
	// Proxy = ""
)

var promos = func() *os.File {
	file, err := os.OpenFile("promos.txt", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}()

func main() {
	defer promos.Close()
	mailers := []mail.Mailer{
		mail.NewTidalMailer("OUTLOOK"),
	}

	// for i := 0; i < ThreadNumber; i++ {
	// 	time.Sleep(25 * time.Millisecond)

	// 	go func() {
	// 		for {
	mailer := mailers[rand.Intn(len(mailers))]
	creator, err := webtoons.NewCreator(mailer)
	if err != nil {
		log.Printf("[x] ERROR %v\n", err)
		// continue
	}
	email, id := mailer.RandomAddress()
	err = creator.Create(email, helpers.RandString(8), id)
	if err != nil {
		log.Printf("[x] ERROR %v\n", err)
		// continue
	}
	time.Sleep(4 * time.Second)
	code, err := creator.RedeemCode()
	if err != nil {
		log.Printf("[x] ERROR %v\n", err)
		// continue
	}

	log.Printf("[*] Got code (%s)\n", code)

	io.WriteString(promos, code+"\n")
	// 		}
	// 	}()
	// }

	// s := make(chan os.Signal, 1)
	// signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	// <-s

}

// func main() {
// 	// defer promos.Close()
// 	mailers := []mail.Mailer{
// 		mail.NewTidalMailer("OUTLOOK"),
// 	}

// 	// for i := 0; i < ThreadNumber; i++ {
// 	// 	time.Sleep(25 * time.Millisecond)

// 	// 	go func() {
// 	// for {
// 	mailer := mailers[rand.Intn(len(mailers))]
// 	creator, err := webtoons.NewCreator(mailer)
// 	if err != nil {
// 		log.Printf("[x] ERROR %v\n", err)
// 		// continue
// 	}
// 	email, id := mailer.RandomAddress()
// 	err = creator.Create(email, helpers.RandString(8), id)
// 	if err != nil {
// 		log.Printf("[x] ERROR %v\n", err)
// 		// continue
// 	}

// 	code, err := creator.RedeemCode()
// 	if err != nil {
// 		log.Printf("[x] ERROR %v\n", err)
// 		// continue
// 	}

// 	// log.Printf("[*] Got code (%s)\n", code)

// 	// io.WriteString(promos, code+"\n")
// 	// 		}
// 	// 	}()
// 	// }

// 	// s := make(chan os.Signal, 1)
// 	// signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
// 	// <-s

// }
