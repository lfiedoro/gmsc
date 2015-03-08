package main

import (
	gc "code.google.com/p/goncurses"
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
	return out
}

func clearAndPlay(what []string, conn *mpd.Client) {
	if err := conn.Clear(); err != nil {
		log.Fatal(err)
	}

	for _, song := range what {
		log.Println("Adding", song, "to playlist")
		if err := conn.Add(song); err != nil {
			log.Fatal(err)
		}
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

	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	var artists []string
	for _, val := range *lib.ArtistList {
		artists = append(artists, val.Name)
	}

	gmsc.Present(&artists, stdscr)
	artistName := gmsc.Choose(&artists, stdscr)

	var artist gmsc.Artist
	for _, art := range *lib.ArtistList {
		if art.Name == artistName {
			artist = art
			break
		}
	}

	var albums []string
	for _, alb := range *artist.AlbumList {
		albums = append(albums, alb.Name)
	}

	gmsc.Present(&albums, stdscr)
	albumName := gmsc.Choose(&albums, stdscr)

	log.Println("Artist:", artistName, "Album:", albumName)

	var album gmsc.Album
	for _, alb := range *artist.AlbumList {
		if alb.Name == albumName {
			album = alb
			break
		}
	}

	var songs []string
	for _, song := range *album.SongList {
		songs = append(songs, song.File)
	}

	log.Println(songs)

	//Clear current playlist and play
	clearAndPlay(songs, conn)
}
