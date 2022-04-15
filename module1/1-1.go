package main

import "fmt"

func main() {
	me := []string{"I", "am", "stupid", "and", "weak"}

	fmt.Println("me: ", me)

	for index, value := range me {
		if value == "stupid" {
			me[index] = "smart"
		}
		if value == "weak" {
			me[index] = "strong"
		}
	}

	fmt.Println("me: ", me)
}
