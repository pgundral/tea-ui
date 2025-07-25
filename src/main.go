package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TeaType string

const (
	Black   TeaType = "Black"
	Oolong  TeaType = "Oolong"
	Green   TeaType = "Green"
	White   TeaType = "White"
	Herbal  TeaType = "Herbal"
)

type BrewStep int

const (
	SelectTea BrewStep = iota
	AdditionsChoice
	PrepMaterials
	BoilWater
	StartBrewing
	Brewing
	AddMixins
	SelectMusic
	Complete
)

type model struct {
	step           BrewStep
	selectedTea    TeaType
	choices        []TeaType
	cursor        int
	addMilk       bool
	addSugar      bool
	brewTimer     time.Duration
	brewStartTime time.Time
}

func initialModel() model {
	return model{
		step:    SelectTea,
		choices: []TeaType{Black, Oolong, Green, White, Herbal},
		brewTimer: 3 * time.Minute, // Default brew time
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.step {
			case SelectTea:
				m.selectedTea = m.choices[m.cursor]
				m.step = AdditionsChoice
			case AdditionsChoice:
				m.step = PrepMaterials
			case PrepMaterials:
				m.step = BoilWater
			case BoilWater:
				m.step = StartBrewing
			case StartBrewing:
				m.brewStartTime = time.Now()
				m.step = Brewing
				return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			case Brewing:
				if time.Since(m.brewStartTime) >= m.brewTimer {
					m.step = AddMixins
				}
			case AddMixins:
				m.step = SelectMusic
			case SelectMusic:
				m.step = Complete
			}
		case "y":
			if m.step == AdditionsChoice {
				m.addMilk = true
				m.addSugar = true
			}
		case "n":
			if m.step == AdditionsChoice {
				m.step = PrepMaterials
			}
		}
	case tickMsg:
		if m.step == Brewing && time.Since(m.brewStartTime) >= m.brewTimer {
			m.step = AddMixins
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

type tickMsg time.Time

func (m model) View() string {
	s := "\n🫖 Tea Brewing Assistant 🫖\n\n"

	switch m.step {
	case SelectTea:
		s += "Select your tea type:\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

	case AdditionsChoice:
		s += fmt.Sprintf("You selected: %s\n\n", m.selectedTea)
		s += "Would you like to add milk and/or sugar? (y/n)\n"

	case PrepMaterials:
		s += "Please prepare:\n"
		s += "- Tea/teabag\n"
		s += "- Mug\n"
		s += "- Water\n"
		if m.addMilk || m.addSugar {
			s += "- Milk/Creamer and/or Sugar\n"
		}
		s += "\nPress ENTER when ready\n"

	case BoilWater:
		temp := getIdealTemp(m.selectedTea)
		s += fmt.Sprintf("Boiling water to %d°C...\n", temp)
		s += "\nPress ENTER when water is ready\n"

	case StartBrewing:
		s += "Ready to start brewing!\n"
		s += "Put teabag in mug and pour water.\n"
		s += "\nPress ENTER to start timer\n"

	case Brewing:
		elapsed := time.Since(m.brewStartTime)
		remaining := m.brewTimer - elapsed
		if remaining < 0 {
			remaining = 0
		}
		s += fmt.Sprintf("Brewing... %s remaining\n", remaining.Round(time.Second))
		s += renderProgressBar(float64(elapsed)/float64(m.brewTimer), 40)

	case AddMixins:
		s += "🎉 Tea is ready!\n\n"
		if m.addMilk || m.addSugar {
			s += "Add your milk/sugar now\n"
		}
		s += "\nPress ENTER to continue\n"

	case SelectMusic:
		s += "Select your playlist:\n"
		s += "(Feature coming soon!)\n"
		s += "\nPress ENTER to finish\n"

	case Complete:
		s += "Enjoy your tea! 🫖\n"
		s += "\nPress 'q' to quit\n"
	}

	s += "\n(Press 'q' to quit at any time)\n"
	return s
}

func getIdealTemp(teaType TeaType) int {
	switch teaType {
	case Black:
		return 100
	case Oolong:
		return 85
	case Green:
		return 80
	case White:
		return 75
	case Herbal:
		return 100
	default:
		return 90
	}
}

func renderProgressBar(percent float64, width int) string {
	filled := int(percent * float64(width))
	if filled > width {
		filled = width
	}
	return fmt.Sprintf("[%s%s]",
		strings.Repeat("=", filled),
		strings.Repeat("-", width-filled))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
