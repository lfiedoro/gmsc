package main

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"log"
	"sort"
)

func getLib(token string, conn *mpd.Client) (out []mpd.Attrs) {
	var err error
	// Get whole library
	if out, err = conn.ListAllInfo(token); err != nil {
		log.Fatal(err)
	}
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

func contains(arr []string, that string) bool {
	for _, this := range arr {
		if this == that {
			return true
		}
	}
	return false
}

func getKeys(m *map[string][]string) []string {
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}

func getAlbumArtist(attrs *[]mpd.Attrs) map[string][]string {

	mapping := make(map[string][]string)

	for _, entry := range *attrs {
		key := entry["Artist"]
		val := entry["Album"]

		if _, ok := mapping[key]; ok {
			if !contains(mapping[key], val) {
				mapping[key] = append(mapping[key], val)
			}
		} else {
			mapping[key] = append(mapping[key], val)
		}
	}
	return mapping
}

func present(list *[]string) {
	for i, v := range *list {
		fmt.Println("[", i, "] ", v)
	}
}

func choose(list *[]string) string {
    // If there is only one element, choose it by default
	if len(*list) == 1 {
		return (*list)[0]
	}

	var input int
	if _, err := fmt.Scanf("%d", &input); err != nil {
		log.Fatal(err)
	}

	var element string
	if input >= 0 && input < len(*list) {
		element = (*list)[input]
	}

	return element
}

func main() {
	// Create connection
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Get Artists and their Albums
	lib := getLib("/", conn)
	artistAlbum := getAlbumArtist(&lib)

	artists := getKeys(&artistAlbum)
	sort.Strings(artists)

	// Choose Artist
	present(&artists)
	artist := choose(&artists)

	albums := artistAlbum[artist]

	// Choose Album
	present(&albums)
	album := choose(&albums)

	// Clear current playlist and play
	fmt.Println("Playing", album, "by", artist)
	clearAndPlay(artist+" - "+album, conn)
}
