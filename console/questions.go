package console

import (
	"fmt"
	"strings"
)

//QyORn asks yes or no and returns bool
func QyORn(question string) bool {
	var answer string

	fmt.Printf("    %s\n", Cyan(question))
	fmt.Println("    ─────")
	fmt.Println("(1) Yes")
	fmt.Println("(2) No")
	fmt.Println("    ─────")

	for true {
		fmt.Printf("  > ")

		fmt.Scan(&answer)
		answer = strings.ToLower(answer)
		if answer == "y" || answer == "yes" || answer == "1" {
			return true
		} else if answer == "n" || answer == "no" || answer == "2" {
			return false
		}
	}
	return false
}
