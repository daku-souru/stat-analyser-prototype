package model

import (
	"fmt"
	"log"
	"math"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

func checkBhop(consecutive_bhops int) int {

	if float64(rounds)/float64(consecutive_bhops) <= C_BHOP_SUS {
		if float64(rounds)/float64(consecutive_bhops) <= C_BHOP_HIGHLY_SUS {
			return V_HIGHLY_SUSPICIOUS
		}
		return V_SUSPICIOUS
	}
	return V_OK
}

func DetectAutomation(demoPath string, steamdID uint64, rounds int) (int, error) {
	f, err := os.Open(demoPath)
	if err != nil {
		log.Panic("failed to open demo file: ", err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	// first jump does not count towads bhop
	jumpCount := -1
	var lastFrameAirborne bool

	var consecutive_bhop_count int
	var bhop_count int

	for {
		// Get player entity
		for _, player := range p.GameState().Participants().All() {
			if player.SteamID64 == steamdID {

				// Bunnyhop detection

				// Get the magnitude of the player's velocity
				velocity := player.Velocity()
				velocityMag := math.Sqrt(math.Pow(velocity.X, 2) + math.Pow(velocity.Y, 2) + math.Pow(velocity.Z, 2))
				nowAirborne := player.IsAirborne()

				// Player jumps
				if nowAirborne && !lastFrameAirborne {
					jumpCount++

					if velocityMag > C_BHOP_MAG {
						bhop_count++

						if bhop_count > C_BHOP_MAX {
							consecutive_bhop_count++
						}
					}

				} else if lastFrameAirborne && !nowAirborne {
					if velocityMag < C_BHOP_MAG {
						bhop_count = -1
					}
				}
				lastFrameAirborne = nowAirborne
			}
		}
		moreFrames, err := p.ParseNextFrame()
		if !moreFrames {
			fmt.Printf("Parser ended after less than %d frames", p.CurrentFrame())
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	return checkBhop(consecutive_bhop_count), nil
}
