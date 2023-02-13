package model

const SHARED_SECRET = "prototype"

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

// Criteria for determining suspicion levels  that can be changed to your liking
const (
	C_BHOP_MAG        = 285.93 // Velocity magnitude
	C_BHOP_MAX        = 4      // If this number is hit consistently throughout, start to suspect player
	C_BHOP_SUS        = 2      // used to check consecutive bhops in relation to rounds
	C_BHOP_HIGHLY_SUS = 2      // used to check consecutive bhops in relation to rounds
	C_KILLS           = 3.0    // Start to suspect if player has this number of kills every round (on average) throughout the match
	C_PRCT_WALL       = 0.7    // Become suspicious if player has this percentage for kills of this type
	C_PRCT_SMOKE      = 0.7    // Become suspicious if player has this percentage for kills of this type
	C_PRCT_HS_SUS     = 0.9    // Only relevant if already suspicious

)
