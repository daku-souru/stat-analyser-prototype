package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
)

var matchID = ""
var userID = 0
var rank = 0
var killCount = 0
var headshotKills = 0
var wallbangKills = 0
var flashedKills = 0
var smokeKills = 0
var headshotPercentage = float32(0)

var duration = 0.0
var rounds = 0

func checkSmokeKills() bool {
	return float64(smokeKills)/float64(killCount) >= C_PRCT_SMOKE
}

func checkWall() bool {
	return float64(wallbangKills)/float64(killCount) >= C_PRCT_WALL
}

func checkHeadshots() bool {
	return float64(headshotKills)/float64(killCount) >= C_PRCT_HS_SUS
}

// check if various conditions are met
func checkStatisticalAnomalies() int {
	if float64(killCount)/float64(rounds) >= C_KILLS {
		if checkHeadshots() || checkWall() || checkSmokeKills() {
			return V_HIGHLY_SUSPICIOUS
		}
		return V_SUSPICIOUS
	} else if checkSmokeKills() {
		return V_SUSPICIOUS
	} else if checkWall() {
		return V_SUSPICIOUS
	} else {
		return V_OK
	}
}

// Demofiles do not contain match IDs, however this would not be needed for official games probably
// So I am just generating an ID based on the hashsum of the .dem file
func generateMatchID(demoPath string) (string, error) {
	f, err := os.Open(demoPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func AnalyseDemo(demoPath string, steamID uint64) (Verdict, error) {
	verdict := Verdict{}

	genID, err := generateMatchID(demoPath)
	if err != nil {
		fmt.Printf("Failed to generate match ID")
		return verdict, err
	}
	verdict.MatchID = genID

	f, err := os.Open(demoPath)
	if err != nil {
		fmt.Printf("Failed to open demo file")
		return verdict, err
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	// Register Events

	// Get user ID
	p.RegisterEventHandler(func(e events.MatchStart) {
		for _, player := range p.GameState().Participants().All() {
			if player.SteamID64 == steamID {
				userID = player.UserID
				verdict.PlayerID = userID
			}
		}
	})

	// Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		// Sometimes a player dies on their own, so check if there is actually a killer
		if e.Killer == nil || e.Killer.UserID != userID {
			return
		}
		killCount++
		if e.IsHeadshot {
			headshotKills++
		}
		if e.PenetratedObjects > 0 {
			wallbangKills++
		}
		if e.ThroughSmoke {
			smokeKills++
		}
	})

	p.RegisterEventHandler(func(e events.RankUpdate) {
		if e.Player.UserID != userID {
			return
		}
		// Rank from 1 (Silver 1) to 18 (Global Elite)
		rank = e.RankNew
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		// Demo is likely corrupted
		verdict.StatResult = V_CORRUPTED
	}

	duration = p.Header().PlaybackTime.Seconds()
	rounds = p.GameState().TotalRoundsPlayed()

	headshotPercentage = (float32(headshotKills) / float32(killCount) * 100)
	verdict.StatResult = checkStatisticalAnomalies()

	automationResult, err := DetectAutomation(demoPath, steamID, rounds)
	if err != nil {
		fmt.Printf("Auomation detection failed: %s", err)
	} else {
		verdict.AutomationResult = automationResult
	}

	// Set final verdict values
	verdict.SteamID = steamID
	verdict.Data = VerdictData{
		killCount,
		headshotPercentage,
		wallbangKills,
		smokeKills,
		flashedKills,
		rank,
		rounds,
		duration}
	return verdict, nil
}
