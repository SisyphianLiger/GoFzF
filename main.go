package main

import (
	"fmt"
	"os"

	// "os/exec"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


var (
    selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)


// Need to be able to highlight selected option...

type mapandindex struct {
    FilandDirMap map[int]struct{}   // which to-do items are selected
    idx int // current location of cursor
}


type model struct {
    FilesAndDirectories  []string           // items on the to-do list
    mapandindex // embedded to reduce clutter
}

type Data func() []string
/*
    fakeData() is to be replaced with data from the pipeline
*/
func initialModel(fn Data) model {
    mapoffile := mapandindex{ 
                    FilandDirMap:make(map[int]struct{}),
                    idx: len(fn()) - 1, 
    }
	return model{
            // Our to-do list is a grocery list
            FilesAndDirectories:  fn(),

            // A map which indicates which FilesAndDirectories are FilandDirMap. We're using
            // the map like a mathematical set. The keys refer to the indexes
            // of the `FilesAndDirectories` slice, above.
            mapandindex: mapoffile,
        }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.idx > 0 {
                m.idx--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.idx < len(m.FilesAndDirectories)-1 {
                m.idx++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the FilandDirMap state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.FilandDirMap[m.idx]
            if ok {
                delete(m.FilandDirMap, m.idx)
            } else {
                m.FilandDirMap[m.idx] = struct{}{}
            }
        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m model) View() string {
    // The header
    s := "Files and Directories\n\n"

    // Iterate over our FilesAndDirectories
    for i, fileordir := range m.FilesAndDirectories {

        // Is the cursor pointing at this fileordir?
        cursor := " " // no cursor
        if m.idx == i {
            cursor = ">" // cursor!
            selectedItemStyle.Render(cursor + " " + fileordir)
        }

        // Render the row
        s += fmt.Sprintf("%s %s\n", cursor, fileordir)
    }

    // The footer
    s += "\nPress q or ctrl+c to quit.\n"

    // Send the UI for rendering
    return s
}

func main() {

    clearScreen()
    p := tea.NewProgram(initialModel(fakeData))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error found: %v", err)
        os.Exit(1)
    }
}
