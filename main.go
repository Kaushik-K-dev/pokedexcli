package main

import ("fmt"; "bufio"; "os"; "strings")

func main() {
	fmt.Println("Pokemon started. Enter command.")
	scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("pokedex > ")
        if !scanner.Scan() {break}

		cmd := scanner.Text()
		cmd = strings.ToLower(strings.TrimSpace(cmd))
		if cmd == "" {continue}
		
		commandList := commandList()
		command, ok := commandList[cmd]
		if !ok {
			fmt.Println("Invalid Command. Use 'help' to see all available commands.")
			continue
		}
		command.callback()
    }
}