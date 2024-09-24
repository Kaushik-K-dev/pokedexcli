package main

import ("fmt";"os"; "strconv")

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string)
}

type config struct {
	prevURL *string
	nextURL *string
	cache    *Cache
}

func commandList() map[string]cliCommand {
	return map[string]cliCommand{
    	"help": {
        	name:        "help",
        	description: "Displays all available commands.",
        	callback:    commandHelp,
    	},
    	"exit": {
        	name:        "exit",
        	description: "Exit the Pokedex.",
        	callback:    commandExit,
    	},
		"map": {
        	name:        "map",
        	description: "Shows the first 20 Locations. Repeat the command to see next 20. Use 'mapb' to go back.",
        	callback:    commandMap,
    	},
		"mapb": {
        	name:        "mapb",
        	description: "Shows the previous result (if any). Repeat the command to go further back.",
        	callback:    commandMapBack,
    	},
		"explore": {
        	name:        "explore",
        	description: "Shows the Pokemon in specified Area.",
        	callback:    commandExplore,
    	},
	}
}

func commandExit(cfg *config, args ...string) {
	fmt.Println("Exiting Pokedex...")
	os.Exit(0)
}

func commandHelp(cfg *config, args ...string) {
	fmt.Println("Here are the list of available commands.")
	commandList := commandList()
	for _, cmd := range commandList {
		fmt.Printf("- %v: %v\n", cmd.name, cmd.description)
	}
}

func commandMap(cfg *config, args ...string) {
	repeatCount := 1
	if len(args) !=0 {
		repeatCount, err := strconv.Atoi(args[0])
		if err != nil || repeatCount <= 0 {
			fmt.Println("Invalid number of repetitions. Please provide a positive integer.")
			return
		}
	}
	resp, err := cfg.LocationListReq(cfg.nextURL)
	if err!=nil {fmt.Errorf("err")}
	for i:=1; i < repeatCount; i++ {
		resp, err = cfg.LocationListReq(cfg.nextURL)
		if err!=nil {fmt.Errorf("err")}
	}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
}

func commandMapBack(cfg *config, args ...string) {
	if cfg.prevURL == nil {
		fmt.Println("Cannot go back anymore.")
		return
	}
	resp, err := cfg.LocationListReq(cfg.prevURL)
	if err!=nil {fmt.Errorf("err")}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
}

func commandExplore(cfg *config, args ...string) {
	if len(args) != 0 {
		resp, err := cfg.LocationDataReq(args[0])
		if err!=nil {
			fmt.Println("Incorrect Location Name provided")
			return
		}
		fmt.Println("Exploring "+args[0]+"...")
		for _, encounter := range resp.PokemonEncounters {
			fmt.Println(encounter.Pokemon.Name)
		}
	} else {
		fmt.Println("No Area Name or incorrect Name provided. (Please check if you have typed the exact name. There should be hyphens instead of spaces.)")
	}
}