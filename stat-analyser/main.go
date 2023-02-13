package main

import (
	"flag"
	"fmt"
	"os"

	actions "github.com/daku-souru/stat-analyser/controller"
	"github.com/daku-souru/stat-analyser/model"
)

// Flags (start arguments)
const (
	FLAG_USE_API   = "prod"
	FLAG_DEMO_PATH = "demo"
	FLAG_SUSPECT   = "suspect"
)

// How to run: go run main.go -demo "{C:/users/full/path/to/match.dem}" -suspect {steamID64}
func main() {
	// Check flags
	demoPath := flag.String(FLAG_DEMO_PATH, "", "Path to demo file")
	suspectID := flag.Uint64(FLAG_SUSPECT, 0, "Steam ID (decimal steamID64) of the suspect")
	prod := flag.Bool(FLAG_USE_API, false, "If production is wanted")
	flag.Parse()

	// Check for errors
	if *demoPath == "" || *suspectID == 0 {
		fmt.Println("Both -demo and -suspect flags are required")
		os.Exit(1)
	}

	verdict, err := model.AnalyseDemo(*demoPath, *suspectID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Verdict: \n %s", actions.ConvertVerdictToJSON(verdict))

	if !*prod {
		return
	}
	if verdict.StatResult > 0 || verdict.AutomationResult > 0 {
		err = actions.SendVerdict(verdict)
		if err != nil {
			// API not reachable
			panic(err)
		}
	}
}
