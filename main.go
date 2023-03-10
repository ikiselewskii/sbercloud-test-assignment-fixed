package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"playlist-grpc/player"
	protos "playlist-grpc/presentation"
	"playlist-grpc/server"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	exit := make(chan os.Signal,1)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
	gs := grpc.NewServer()
	
	playlist := server.New(log.Default(),"playlist.txt")
	
	protos.RegisterPlaylistServer(gs,playlist)

	reflection.Register(gs)

	l,err := net.Listen("tcp",":9092")
	if err != nil{
		panic(err)
	}
	go func(){
		gs.Serve(l)
	}()
	for{
		select{
		case <- exit:
			gs.Stop()
			l.Close()
			err := playlist.Save()
			if err != nil{
				log.Panic(err)
			}
			return
		case <- player.FinishedTrack:
			playlist.Next(context.Background(),&emptypb.Empty{})
		default:
			continue
		}
	}
}
