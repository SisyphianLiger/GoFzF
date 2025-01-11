package main

import (
        "os"
	"os/exec"
)


// Figure out windows
func clearScreen() {
    clear := exec.Command("clear")
    clear.Stdout = os.Stdout
    clear.Run()
}
