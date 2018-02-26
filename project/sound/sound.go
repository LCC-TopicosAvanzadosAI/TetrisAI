package sound

import (
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

func Play() {

	//tetris audio
	f, _ := os.Open("./../../resources/korobeiniki.mp3")
	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(s)
	//s.Loop(-1, s)

}
