package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

/*
Define a model:
*/





type model struct {
	Directories []string 
	cursor	int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		Directories: fakeData(),		
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}


func main() {
	fmt.Println("Hello WOrld")
}
