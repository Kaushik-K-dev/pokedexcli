package main

import ("fmt";"os")

type cliCommand struct {
	name        string
	description string
	callback    func()
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
	}
}

func commandExit() {
	fmt.Println("Exiting Pokedex...")
	os.Exit(0)
}

func commandHelp() {
	fmt.Println("Here are the list of available commands.")
	commandList := commandList()
	for _, cmd := range commandList {
		fmt.Printf("- %v: %v\n", cmd.name, cmd.description)
	}
}