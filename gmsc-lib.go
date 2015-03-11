package gmsc

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
)

// Song
type Song struct {
	Track string
	File  string
	Name  string
}

func (s Song) String() string {
	return fmt.Sprintf("\n%v %v %v", s.Track, s.File, s.Name)
}

// Album
type Album struct {
	Name     string
	SongList *[]Song
}

func (a Album) String() string {
	return fmt.Sprintf("\n%v: %v", a.Name, a.SongList)
}

func (lib *Album) AddSong(songName string, songTrack string, songFile string) {
	if lib.SongList == nil {
		lib.SongList = new([]Song)
	}
	*lib.SongList = append(*lib.SongList, Song{songTrack, songFile, songName})
}

// Artist
type Artist struct {
	Name      string
	AlbumList *[]Album
}

func (a Artist) String() string {
	return fmt.Sprintf("\n%v: %v", a.Name, a.AlbumList)
}

func (artist *Artist) ContainsAlbum(that string) (bool, *Album) {
	for _, this := range *artist.AlbumList {
		if this.Name == that {
			return true, &this
		}
	}
	return false, nil
}

func (lib *Artist) AddAlbum(albumName string) *Album {
	if lib.AlbumList == nil {
		lib.AlbumList = new([]Album)
	}
	*lib.AlbumList = append(*lib.AlbumList, Album{albumName, nil})
	return &((*lib.AlbumList)[len(*lib.AlbumList)-1])
}

// Library
type Library struct {
	ArtistList *[]Artist
}

func (l Library) String() string {
	return fmt.Sprintf("%v", l.ArtistList)
}

func (lib *Library) ContainsArtist(that string) (bool, *Artist) {
	if lib.ArtistList == nil {
		return false, nil
	}
	for _, this := range *lib.ArtistList {
		if this.Name == that {
			return true, &this
		}
	}
	return false, nil
}

func (lib *Library) AddArtist(artistName string) *Artist {
	if lib.ArtistList == nil {
		lib.ArtistList = new([]Artist)
	}
	*lib.ArtistList = append(*lib.ArtistList, Artist{artistName, nil})
	return &((*lib.ArtistList)[len(*lib.ArtistList)-1])
}

func assignIfEmpty(strin string) string {
	if strin == "" {
		return "<empty>"
	}

	return strin
}

func (lib *Library) Update(attrs *[]mpd.Attrs) {

	for _, entry := range *attrs {
		artistName := assignIfEmpty(entry["Artist"])
		albumName := assignIfEmpty(entry["Album"])
		titleName := assignIfEmpty(entry["Title"])
		fileName := assignIfEmpty(entry["file"])
		trackName := assignIfEmpty(entry["Track"])

		if b, cart := lib.ContainsArtist(artistName); !b {
			// add artist to lib
			artist := lib.AddArtist(artistName)
			// add album to artist
			album := artist.AddAlbum(albumName)
			// add song to album
			album.AddSong(titleName, trackName, fileName)
		} else if b, calb := cart.ContainsAlbum(albumName); !b {
			// add album to artist
			album := cart.AddAlbum(albumName)
			// add song to album
			album.AddSong(titleName, trackName, fileName)
		} else {
			// add song to album
			calb.AddSong(titleName, trackName, fileName)
		}
	}
}
