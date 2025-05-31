package Game

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	Aries       = 0
	Taurus      = 1
	Gemini      = 2
	Cancer      = 3
	Leo         = 4
	Virgo       = 5
	Libra       = 6
	Scorpio     = 7
	Sagittarius = 8
	Capricorn   = 9
	Aquarius    = 10
	Pisces      = 11
)

var currentZodiac atomic.Int32

var zodiacState atomic.Bool

func StopZodiacCycle() {
	zodiacState.Store(false)
}
func StartZodiacCycle() {
	if zodiacState.Load() {
		return
	}

	zodiacState.Store(true)

	go func() {
		for zodiacState.Load() {
			currentZodiac.Store(int32(rand.Intn(11)))

			fmt.Printf("Set new blessed zodiac to %s\n", GetZodiacName(int(currentZodiac.Load())))
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
			time.Sleep(time.Until(nextMidnight))
		}
	}()
}

func GetZodiacName(sign int) string {
	var zodiacName string

	switch sign {
	case Aries:
		zodiacName = "Aries"
		break
	case Taurus:
		zodiacName = "Taurus"
		break
	case Gemini:
		zodiacName = "Gemini"
		break
	case Cancer:
		zodiacName = "Cancer"
		break
	case Leo:
		zodiacName = "Leo"
		break
	case Virgo:
		zodiacName = "Virgo"
		break
	case Libra:
		zodiacName = "Libra"
		break
	case Scorpio:
		zodiacName = "Scorpio"
		break
	case Sagittarius:
		zodiacName = "Sagittarius"
		break
	case Capricorn:
		zodiacName = "Capricorn"
		break
	case Aquarius:
		zodiacName = "Aquarius"
		break
	case Pisces:
		zodiacName = "Pisces"
		break
	default:
		zodiacName = "Unknown"
		break
	}

	return zodiacName
}

func GetZodiacSign() int {
	return int(currentZodiac.Load())
}
