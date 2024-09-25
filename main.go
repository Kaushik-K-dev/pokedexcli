package main

import ("fmt"; "bufio"; "os"; "strings"; "time")

func main() {
	cache := newCache(5 * time.Minute)
	cfg :=  &config{
		cache: cache,
		PokeCollection: make(map[string]Pokemon),
	}
	fmt.Println("Pokedex started. Enter command.")
	scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("pokedex > ")
        if !scanner.Scan() {break}

		cmd := scanner.Text()
		cmd = strings.ToLower(strings.TrimSpace(cmd))
		if cmd == "" {continue}

		cmdslice := strings.Fields(cmd)
		
		commandList := commandList()
		command, ok := commandList[cmdslice[0]]
		if !ok {
			fmt.Println("Invalid Command. Use 'help' to see all available commands.")
			continue
		}
		args := cmdslice[1:]
		command.callback(cfg, args...)
    }
}