package main

import (
	"github.com/charmbracelet/bubbles/list"
)

type FilesAndDir struct {
	file string
	dir string
	folder int
}

func (fnd FilesAndDir) FilterValue() string { return string(fnd.folder) }

func (fnd FilesAndDir) fileName() string { return fnd.file }
func (fnd FilesAndDir) dirName() string { return fnd.dir }
func (fnd FilesAndDir) folderNumber() int { return fnd.folder }

func fakeData() []list.Item {
	return []list.Item{
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
		FilesAndDir{ file: "", dir: "docs/project_notes.txt", folder: 1 },
	}
}


func MainMenu() []list.Item {
	return []list.Item{
		MenuOption{ title: "Start Searching"},
		MenuOption{ title: "Configuration Options"},
		MenuOption{ title: "Color Schemes"},
		MenuOption{ title: "Quit"},
	}
}
