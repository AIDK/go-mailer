package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

// we'll use these constants to keep track of which input we're focused on
// and to make it easier to update the model
const (
	to = iota // iota is a special golang constant that starts at 0 and increments by 1 for each const it's used in
	from
	subject
	body
)

const (
	hotPink  = lipgloss.Color("#FF0687") // a nice hot pink
	darkGrey = lipgloss.Color("#767676") // a dark grey
)

// we'll use these styles to render the inputs and the continue prompt
var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGrey)
)

// model is the main model for the program
// it contains a slice of text inputs, the index of the currently focused input,
type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

// toAddressValidator is a simple validator for email addresses
func toAddressValidator(s string) error {
	// TODO: Implement email to address validation
	return nil
}

// fromAddressValidator is a simple validator for email addresses
func fromAddressValidator(s string) error {
	// TODO: Implement email from address validation
	return nil
}

// subjectValidator is a simple validator for email subjects
func subjectValidator(s string) error {
	// TODO: Implement email subject validation
	return nil
}

// bodyValidator is a simple validator for email bodies
func bodyValidator(s string) error {
	// TODO: Implement email body validation
	return nil
}

// initialModel returns the initial model for the program
func initialModel() model {
	//TODO: Add more inputs for the email body, subject, and from address

	// we'll create a slice of text inputs (for now just one)
	var inputs []textinput.Model = make([]textinput.Model, 1)
	inputs[to] = textinput.New()
	inputs[to].Placeholder = "johndoe@domain.com"
	inputs[to].Focus()
	inputs[to].CharLimit = 50
	inputs[to].Width = 50
	inputs[to].Prompt = ""
	inputs[to].Validate = toAddressValidator

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

// Init initializes the model with a command to blink the cursor
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// we'll need to update each of the text inputs, so we'll create a slice
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	// we'll handle the messages for each input and update them accordingly
	switch msg := msg.(type) {

	// KeyMsg is sent when a key is pressed while the component is in focus
	case tea.KeyMsg:

		// we want to handle the key presses for the inputs ourselves
		switch msg.Type {

		// enter will move the focus to the next input
		case tea.KeyEnter:
			// if we're on the last input, we'll quit the program
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()

		// esc and ctrl+c will quit the program
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		// tab and shift+tab will move the focus to the next and previous inputs
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()

		// tab and shift+tab will move the focus to the next and previous inputs
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}

		// we blur all the inputs so we can focus the one we want
		for i := range m.inputs {
			m.inputs[i].Blur()
		}

		// we focus the input we want
		m.inputs[m.focused].Focus()

	// errMsg is sent when an error is returned from a text input's Validate function
	case errMsg:

		// we'll set the error on the model so we can display it in the view
		m.err = msg
		return m, nil
	}

	// we loop through the inputs and update them with the message we received
	for i := range m.inputs {
		// we update the input and store the command it returns,
		// so we can return a batch of all the commands
		// we also store the updated input in the inputs slice
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// we return the updated model and a batch of all the commands we received
	return m, tea.Batch(cmds...)
}

// View renders the model to the screen
func (m model) View() string {
	return fmt.Sprintf(`
	%s
	%s
	
	%s`,
		inputStyle.Width(50).Render("To:"),
		m.inputs[to].View(),
		continueStyle.Render("Continue ->")) + "\n"
}

// nextInput focuses on the next input
func (m *model) nextInput() {
	// we want to focus on the next input by incrementing the focused index
	// and wrapping around to the beginning if we're at the end
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses on the previous input
func (m *model) prevInput() {
	// we want to focus on the previous input by decrementing the focused index
	// and wrapping around to the end if we're at the beginning
	m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
}
