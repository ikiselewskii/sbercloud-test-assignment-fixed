package domain

import (
	"bufio"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"log"
	"os"
	protos "playlist-grpc/presentation"
	"playlist-grpc/player"
	"strings"
	"time"
)

type Track struct {
	Title    string
	Author   string
	Duration time.Duration
	prev     *Track
	next     *Track
}



func New(filename string) *Track {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	dummy := &Track{}
	current := dummy
	for scanner.Scan() {
		line := scanner.Text()
		trackProps := strings.Split(line, " ")
		if len(trackProps) != 3 {
			log.Printf("Invalid track format: %s\n", line)
			continue
		}
		duration, err := time.ParseDuration(trackProps[2])
		if err != nil {
			log.Printf("Invalid track duration: %s\n", trackProps[2])
			continue
		}
		newTrack := &Track{
			Title:    strings.ReplaceAll(trackProps[0], "_", " "),
			Author:   strings.ReplaceAll(trackProps[1], "_", " "),
			Duration: duration,
			prev:     current,
		}
		current.next = newTrack
		current = newTrack
	}
	current.next = nil
	Current := *dummy.next
	Current.prev = nil
	dummy = nil

	return &Current
}

func (p *Track) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	for p.prev != nil {
		p = p.prev
	}
	for node := p; node != nil; node = node.next {
		title := strings.ReplaceAll(node.Title, " ", "_")
		author := strings.ReplaceAll(node.Author, " ", "_")
		duration := node.Duration.String()
		_, err := fmt.Fprintf(writer, "%s %s %s\n", title, author, duration)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *Track) AddSong(song *protos.Song) (error) {
	if song.Duration.AsDuration() <= time.Second*time.Duration(30) {
		return status.Errorf(codes.InvalidArgument, "Invalld song duration")
	}
	for ; p.next != nil; p = p.next {
	}
	p.next = &Track{
		Title:    song.Title,
		Author:   song.Author,
		Duration: song.Duration.AsDuration(),
		prev:     p,
	}
	return status.Errorf(
		codes.OK,
		"Song %s by %s added succesfully",
		song.Title,
		song.Author,
	)
}

func (p *Track) Delete() (*Track,error) {
	if player.IsPlaying {
		return p,status.Errorf(codes.FailedPrecondition, "Can`t delete a playing track")
	}
	if p == nil {
		return nil,status.Errorf(codes.NotFound, "Nothing to delete")
	}
	if p.prev == nil {
		p = p.next
		p.prev = nil
		return p,status.Errorf(codes.OK, "Delete succesful")

	}
	if p.next == nil {
		p = p.prev
		p.next = nil
		return p,status.Errorf(codes.OK, "Delete succesful")

	}
	p.next.prev,p.prev.next = p.prev,p.next
	p = p.prev
	player.PausedOn = 0
	return p,status.Errorf(codes.OK, "Delete succesful")
}

func (p *Track) Next() (*Track, error) {
	if player.IsPlaying{
		player.PauseCh <- struct{}{}
	}
	if p == nil{
		return nil, status.Error(codes.NotFound,"Playlist is empty")
	}
	if p.next != nil {
		player.PausedOn = 0
	} else {
		return p,status.Errorf(codes.NotFound, "This is the last track")
	}

	return p.next,status.Errorf(codes.OK, "Switched to next")
}

func (p *Track) Prev() (*Track, error) {
	if player.IsPlaying{
		player.PauseCh <- struct{}{}
	}
	if p == nil{
		return nil,status.Error(codes.NotFound,"Playlist is empty")
	}
	if p.prev != nil {
		player.PausedOn = 0
	} else {
		return p, status.Errorf(codes.NotFound, "This is the first track")
	}
	return p,status.Errorf(
		codes.OK, "Now playing previous",
	)
}

func (p *Track) Play() (error) {
	if p == nil{
		return status.Error(codes.NotFound,"Playlist is empty")
	}
	if player.IsPlaying {
		return status.Errorf(
			codes.FailedPrecondition,
			"Something is already playing",
		)
	}
	go player.Player(p.Duration)
	return status.Errorf(
		codes.OK,
		"Playing",
	)
}

func (p *Track) Pause() (error) {
	if p == nil{
		return status.Error(codes.NotFound,"Playlist is empty")
	}
	if !player.IsPlaying {
		return status.Errorf(codes.FailedPrecondition, "error")
	}
	player.PauseCh <- struct{}{}
	title := p.Title
	author := p.Author
	pauseTime := player.PausedOn.String()
	return status.Errorf(codes.OK, "Track %v by %v is paused on %v",
		title,
		author,
		pauseTime,
	)
}

func (p *Track) NowPlaying() (*protos.Track, error) {
	if p == nil{
		return nil, status.Error(codes.NotFound,"Playlist is empty")
	}
	if !player.IsPlaying {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"Nothing is on, %s by %s is paused on %s",
			p.Title,
			p.Author,
			player.PausedOn,
		)
	}
	player.PlaytimeReq <- struct{}{}
	playtime := <-player.PlaytimeResp
	return &protos.Track{
			Title:    p.Title,
			Author:   p.Author,
			Duration: durationpb.New(p.Duration),
			Playtime: durationpb.New(playtime),
		},
		status.Errorf(
			codes.OK,
			"Now playing",
		)
}


