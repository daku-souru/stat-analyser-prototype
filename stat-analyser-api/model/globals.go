package model

const SHARED_SECRET = "prototype"

// Mirrored from stat-analyser

type Verdict struct {
	MatchID          string      `json:"match_id"`
	SteamID          uint64      `json:"steam_id"`
	PlayerID         int         `json:"player_id"`
	StatResult       int         `json:"stat_result"`
	AutomationResult int         `json:"automation_result"`
	Data             VerdictData `json:"data"`
}

type VerdictData struct {
	Kills              int     `json:"kills"`    // number of all kills
	HeadshotPercentage float32 `json:"hs"`       // number of all kills
	Wall               int     `json:"wall"`     // number of kills through objects/walls
	Smoke              int     `json:"smoke"`    // number of kills through smoke
	Flashed            int     `json:"flashed"`  // number of kills while flashed
	Rank               int     `json:"rank"`     // rank of suspect
	Rounds             int     `json:"rounds"`   // number of played rounds
	Duration           float64 `json:"duration"` // duration of the match in seconds
}

// VerdictTypes
const (
	V_CORRUPTED         = -1
	V_OK                = 0
	V_SUSPICIOUS        = 1
	V_HIGHLY_SUSPICIOUS = 2
)
