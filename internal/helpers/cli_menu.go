package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// displays a CLI menu with multiple options and returns the selected option
func DisplayMenu(prompt string, options []string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Display the prompt and options
		fmt.Println(prompt)
		for i, option := range options {
			fmt.Printf("%d. %s\n", i+1, option)
		}

		fmt.Print("Enter the number of your choice: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var choice int
		_, err := fmt.Sscanf(input, "%d", &choice)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		return options[choice-1]
	}
}
