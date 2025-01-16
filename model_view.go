package main

func (m Model) View() string {
	// When we exit the program we make sure to render nothing to the screen
	// therefore it will return to normal
	if m.state == Searching {
		return "Start Search!"
	}
	if m.state == Configuration {
		return CreateConfigurationTemplate(m.width, m.height)
	}
	if m.state == Quit {
		return ""
	}

	return m.MainMenuView()
}
