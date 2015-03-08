package main

import (
	//	gc "code.google.com/p/goncurses"
	"fmt"
	"github.com/fhs/gompd/mpd"
	"github.com/lfiedoro/gmsc"
	"log"
	"os"
)

func getLib(token string, conn *mpd.Client) (out []mpd.Attrs) {
	var err error
	// Get whole library
	if out, err = conn.ListAllInfo(token); err != nil {
		log.Fatal(err)
	}
	log.Println(out)
	return out
}

func clearAndPlay(what string, conn *mpd.Client) {
	if err := conn.Clear(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(what)
	if err := conn.Add(what); err != nil {
		log.Fatal(err)
	}

	if err := conn.Play(-1); err != nil {
		log.Fatal(err)
	}
}

func initMpdConnection() *mpd.Client {
	// Create connection
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}

	return conn
}

func main() {

	conn := initMpdConnection()
	defer conn.Close()

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Start")

	// Get Artists and their Albums
	mpdlib := getLib("/", conn)
	var lib gmsc.Library
	lib.Update(&mpdlib)

	fmt.Println()
	fmt.Println()
	fmt.Println(lib)

	//stdscr, err := gc.Init()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer gc.End()

	// Clear current playlist and play
	//clearAndPlay(artist+" - "+album, conn)
}
