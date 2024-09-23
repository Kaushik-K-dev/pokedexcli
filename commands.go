package main

import ("fmt";"os")

type cliCommand struct {
	name        string
	description string
	callback    func(*config)
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
        	description: "Displays a help message",
        	callback:    commandHelp,
    	},
    	"exit": {
        	name:        "exit",
        	description: "Exit the Pokedex",
        	callback:    commandExit,
    	},
		"map": {
        	name:        "map (next page)",
        	description: "Shows the first 20 Locations. Repeat the command to see next 20. Use 'mapb' to go back.",
        	callback:    commandMap,
    	},
		"mapb": {
        	name:        "map prev page",
        	description: "Shows the previous result. Repeat the command to go further back.",
        	callback:    commandMapBack,
    	},
	}
}

func commandExit(cfg *config) {
	fmt.Println("Exiting Pokedex...")
	os.Exit(0)
}

func commandHelp(cfg *config) {
	fmt.Println("Here are the list of available commands.")
	commandList := commandList()
	for _, cmd := range commandList {
		fmt.Printf("- %v: %v\n", cmd.name, cmd.description)
	}
}

func commandMap(cfg *config) {
	resp, err := cfg.LocationReq(cfg.nextURL)
	if err!=nil {fmt.Errorf("err")}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
}

func commandMapBack(cfg *config) {
	if cfg.prevURL == nil {
		fmt.Println("Cannot go back anymore.")
		return
	}
	resp, err := cfg.LocationReq(cfg.prevURL)
	if err!=nil {fmt.Errorf("err")}
	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}
}