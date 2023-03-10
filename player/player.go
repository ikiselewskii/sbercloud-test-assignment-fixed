package player

import (
	"time"
)

var (
	IsPlaying    bool
	PausedOn     time.Duration
	PauseCh      chan struct{}      = make(chan struct{})
	PlaytimeReq  chan struct{}      = make(chan struct{})
	PlaytimeResp chan time.Duration = make(chan time.Duration)
	FinishedTrack chan struct{} 	= make(chan struct{}, 1)
)


func Player(duratuion time.Duration){
	IsPlaying = true
	start := time.Now()
	timer := time.NewTimer(duratuion - PausedOn)
	for{
		select{
		case <- timer.C:
			PausedOn = 0
			IsPlaying = false
			FinishedTrack <- struct{}{}
			return
		case <-PauseCh:
			PausedOn = time.Since(start)
			IsPlaying = false
			return
		case <-PlaytimeReq:
			playtime := PausedOn + time.Since(start)
			PlaytimeResp <- playtime
			continue
		default:
			continue
		}
	}
}