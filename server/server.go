package server

import (
	"context"
	"log"
	"playlist-grpc/domain"
	protos "playlist-grpc/presentation"
	"sync"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PlaylistServer struct{
	l *log.Logger
	current	*domain.Track
	filename string
	mu sync.RWMutex
}

func New(l*log.Logger,filename string)*PlaylistServer{
	track := domain.New(filename)
	return &PlaylistServer{l:l, current: track,filename: filename}
}

func(p*PlaylistServer) Save()(error){
	p.mu.RLock()
	err := p.current.Save(p.filename)
	p.mu.RUnlock()
	return err
}

func (p*PlaylistServer)AddSong(ctx context.Context, song *protos.Song)(*emptypb.Empty, error){
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.current.AddSong(song)
	return &emptypb.Empty{},err
}

func (p*PlaylistServer)Delete(ctx context.Context, void *emptypb.Empty)(*emptypb.Empty,error){
	p.mu.Lock()
	defer p.mu.Unlock()
	new,err := p.current.Delete()
	p.current = new
	return void,err
}

func (p*PlaylistServer)Next(ctx context.Context, void *emptypb.Empty)(*emptypb.Empty,error){
	p.mu.Lock()
	next,err := p.current.Next()
	p.current = next
	p.mu.Unlock()
	p.Play(ctx, void)
	return void,err
}

func (p*PlaylistServer)Prev(ctx context.Context, void *emptypb.Empty)(*emptypb.Empty,error){
	p.mu.Lock()
	prev,err := p.current.Prev()
	p.current = prev
	p.mu.Unlock()
	p.Play(ctx,void)
	return void,err
}

func (p*PlaylistServer)Play(ctx context.Context, void *emptypb.Empty)(*emptypb.Empty,error){
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.current.Play()
	return void,err
}

func (p*PlaylistServer)Pause(ctx context.Context, void *emptypb.Empty)(*emptypb.Empty,error){
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.current.Pause()
	return void,err
}

func (p*PlaylistServer)NowPlaying(ctx context.Context, void *emptypb.Empty)(*protos.Track, error){
	p.mu.RLock()
	defer p.mu.RUnlock()
	track, err := p.current.NowPlaying()
	return track, err
}