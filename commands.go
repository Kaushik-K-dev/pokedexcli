package main

import ("fmt";"os"; "strconv"; "math/rand"; "strings")

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string)
}

type config struct {
	prevURL *string
	nextURL *string
	cache    *Cache
	PokeCollection map[string]Pokemon
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
		"catch": {
        	name:        "catch",
        	description: "Attempt to catch specified Pokemon.",
        	callback:    commandCatch,
    	},
		"inspect": {
        	name:        "inspect",
        	description: "Inspect a caught Pokemon. Fails if you have not caught it yet.",
        	callback:    commandInspect,
    	},
		"list": {
        	name:        "list",
        	description: "Shows your collection of caught Pokemon.",
        	callback:    commandPokeList,
    	},
		"release": {
        	name:        "release",
        	description: "Releases specified Pokemon. Use release all to clear the Pokedex.",
        	callback:    commandRelease,
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

func commandCatch(cfg *config, args ...string) {
	if len(args[0]) == 0 {fmt.Println("Please enter a Pokemon name to catch.")}

	if _, ok := cfg.PokeCollection[args[0]]; ok {
		fmt.Println("You have already caught "+args[0]+"!")
		return
	}
	ballmult := 1.0
	if len(args) > 1 {
		switch strings.ToLower(args[1]){
		case "great": ballmult = 1.5
		case "ultra": ballmult = 2.0
		case "master": ballmult = 255.0
		case "poke":
		default: fmt.Println("Invalid Pokeball type. Using standard Pokeball")
				args[1] = "poke"
		}
		args[1] = strings.Title(args[1])
		fmt.Println("Throwing a "+args[1]+"ball at "+args[0]+"...")
	} else {
	fmt.Println("Throwing a Pokeball at "+args[0]+"...")
	}

	resp, resp2, err := cfg.PokemonCatch(args[0])
	if err!=nil {fmt.Errorf("err")}

	catchRate := resp.CaptureRate * int(ballmult)
	if rand.Intn(256) > catchRate {
		fmt.Println(args[0]+" escaped!")
		return
	}
	var Pkm Pokemon
	Pkm.Name = args[0]
	Pkm.Ability = strings.Title(resp2.Abilities[rand.Intn(len(resp2.Abilities))].Ability.Name)
	if resp.GenderRate == -1 {
		Pkm.Gender = ""
	} else if rand.Intn(8) < resp.GenderRate {
		Pkm.Gender = "(\u2640)"
	} else {Pkm.Gender = "(\u2642)"}
	Pkm.Height = resp2.Height
	Pkm.Weight = resp2.Weight
	if len(resp2.HeldItems) != 0 {
	Pkm.HeldItem = strings.Title(resp2.HeldItems[rand.Intn(len(resp2.HeldItems))].Item.Name)
	} else {Pkm.HeldItem = "None"}
	Pkm.Types = make([]string, 0)
	for _, typ := range resp2.Types {
		Pkm.Types = append(Pkm.Types, strings.Title(typ.Type.Name))
	}

	Pkm.Stats = make([]string, 0)
	for _, stat := range resp2.Stats {
		statstr := stat.Stat.Name + ": " + strconv.Itoa(stat.BaseStat)+ " (" + strconv.Itoa(rand.Intn(32)) +")"
		Pkm.Stats = append(Pkm.Stats, statstr)
	}

	Pkm.Moves = make([]string, 0)
	for i:=0; i < min(4, len(resp2.Moves)); i++ {
		Pkm.Moves = append(Pkm.Moves, resp2.Moves[i].Move.Name)
	}

	shinytext := ""
	if rand.Intn(4096) == 0 {
		Pkm.Is_shiny = true
		shinytext = "a Shiny "
	} else {Pkm.Is_shiny = false}
	Pkm.NatDexIndex = resp.ID
	DexEntry := resp.FlavorTextEntries[0].FlavorText
	DexEntry = strings.ReplaceAll(DexEntry, "\n", " ")
	DexEntry = strings.ReplaceAll(DexEntry, "\f", " ")

	Pkm.Dex_entry = DexEntry

	cfg.PokeCollection[args[0]] = Pkm

	fmt.Println("You caught "+shinytext+args[0]+"!")
}

func commandPokeList(cfg *config, args ...string) {
	if len(cfg.PokeCollection) == 0{
		fmt.Println("You don't have any Pokemon.")
		return
	}
	for _, pokemon := range cfg.PokeCollection {
		fmt.Println("- " + pokemon.Name)
	}
}

func commandRelease(cfg *config, args ...string) {
	if args[0] == "all" {
		cfg.PokeCollection = make(map[string]Pokemon)
		fmt.Println("All Pokemon have been released!")
		return
	}
	if _, ok := cfg.PokeCollection[args[0]]; ok {
		delete(cfg.PokeCollection, args[0])
		fmt.Println("You have released "+args[0]+"!")
		return
	}
	fmt.Println("You do not have that in your collection. Please check if you have entered the name correctly.")
}

func commandInspect(cfg *config, args ...string) {
	if len(args[0]) == 0 {fmt.Println("Please enter a Pokemon name in your collection.")}

	Pkm, ok := cfg.PokeCollection[args[0]]
	if !ok {
		fmt.Println("You do not have that in your collection!")
		return
	}
	shinytext := ""
	if Pkm.Is_shiny {shinytext += " (Shiny!!)"}
	fmt.Println("Name: "+strings.Title(Pkm.Name)+Pkm.Gender+shinytext)
	fmt.Println("Ability: "+Pkm.Ability)
	fmt.Println("Held Item: "+Pkm.HeldItem)
	fmt.Println("Height: "+strconv.Itoa(Pkm.Height))
	fmt.Println("Weight: "+strconv.Itoa(Pkm.Weight))
	typeline := "Type(s): "+Pkm.Types[0]
	if len(Pkm.Types) > 1 {typeline += " "+Pkm.Types[1]}
	fmt.Println(typeline)
	fmt.Println("\nBase Stats (EVs):")
	for _, stat := range Pkm.Stats {
		fmt.Println(strings.Title(stat))
	}
	fmt.Println("\nMoves: ")
	for _, move := range Pkm.Moves {
		fmt.Println("- "+strings.Title(move))
	}
	fmt.Println("\nEntry No. "+strconv.Itoa(Pkm.NatDexIndex)+":\n"+Pkm.Dex_entry)
}