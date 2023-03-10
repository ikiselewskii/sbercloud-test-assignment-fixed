package domain

import (
	__ "playlist-grpc/presentation"
	"testing"
	"time"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestNext(t *testing.T) {
	first := &Track{
		Title: "The Fool",
		Author: "Neutral Milk Hotel",
		Duration: time.Minute + time.Second*53,
		}
	second := &Track{
		Title: "1969",
		Author: "Boards of Canada",
		Duration: time.Minute*4 + time.Second*20,
		prev:	first,
	}
	first.next = second
	first,err := first.Next()
	if err != nil{
		t.Fatal(err)
	}
	if first.Title != second.Title || first.Author != second.Author || first.Duration != second.Duration{
		t.Fatal("Failed, tracks are different")
	}
	
}

func TestPrev(t *testing.T){
	second := &Track{
		Title: "The Fool",
		Author: "Neutral Milk Hotel",
		Duration: time.Minute + time.Second*53,
		}
	first := &Track{
		Title: "1969",
		Author: "Boards of Canada",
		Duration: time.Minute*4 + time.Second*20,
		prev:	second,
	}
	second.prev = first
	second,err := first.Prev()
	if err != nil{
		t.Fatal(err)
	}
	if first.Title != second.Title || first.Author != second.Author || first.Duration != second.Duration{
		t.Fatal("Failed, tracks are different")
	}
}

func TestDelete(t *testing.T){
	first := &Track{
		Title: "The Fool",
		Author: "Neutral Milk Hotel",
		Duration: time.Minute + time.Second*53,
		}
	second := &Track{
		Title: "1969",
		Author: "Boards of Canada",
		Duration: time.Minute*4 + time.Second*20,
		prev:	first,
	}
	first.next = second
	third := &Track{Title: "Cornish Acid",
		Author: "Aphex Twin",
		Duration: time.Minute*2 + time.Second*15,
		prev:	second,}
	second.next = third
	track,err := second.Delete()
	if err != nil{
		t.FailNow()
	}
	if track != first{
		t.FailNow()
	}
	
}

func TestAddSong(t *testing.T){
	first := &Track{
		Title: "The Fool",
		Author: "Neutral Milk Hotel",
		Duration: time.Minute + time.Second*53,
		}
	second := &Track{
		Title: "1969",
		Author: "Boards of Canada",
		Duration: time.Minute*4 + time.Second*20,
	}
	first.AddSong(&__.Song{
		Title: second.Title,
		Author: second.Author,
		Duration: durationpb.New(second.Duration),
	})
	if first.next.Title != second.Title || first.next.Author != second.Author || first.next.Duration != second.Duration{
		t.Fatal("Failed, tracks are different")
	}

}