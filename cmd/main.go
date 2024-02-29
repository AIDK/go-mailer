package main

import (
	"fmt"
	"log"
	"net/mail"

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

// Model is the main Model for the program
// it contains a slice of text inputs, the index of the currently focused input,
type model struct {
	inputs  []textinput.Model
	focused int
	err     error
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

// validateAddress validates the email address
func (m model) validateAddress() error {

	var err error

	c := m.inputs[m.focused]
	if _, err = mail.ParseAddress(c.Value()); err != nil {
		err = fmt.Errorf("invalid email address")
	}

	return err
}

// func stringValidator(s string) error {

// 	if len(s) == 0 {
// 		return fmt.Errorf("input cannot be empty")
// 	}

// 	return nil
// }

// initialModel returns the initial model for the program
func initialModel() model {
	// we'll create a slice of text inputs (for now just one)
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[to] = textinput.New()
	inputs[to].Placeholder = "Enter to address here..."
	inputs[to].Focus()
	inputs[to].CharLimit = 50
	inputs[to].Width = 50
	inputs[to].Prompt = ""
	// inputs[to].Validate = stringValidator

	inputs[from] = textinput.New()
	inputs[from].Placeholder = "Enter from address here..."
	inputs[from].CharLimit = 50
	inputs[from].Width = 50
	inputs[from].Prompt = ""
	// inputs[from].Validate = stringValidator

	inputs[subject] = textinput.New()
	inputs[subject].Placeholder = "Enter subject here..."
	inputs[subject].CharLimit = 50
	inputs[subject].Width = 50
	inputs[subject].Prompt = ""

	inputs[body] = textinput.New()
	inputs[body].Placeholder = "Send a message..."
	inputs[body].CharLimit = 50
	inputs[body].Width = 50
	inputs[body].Prompt = ""

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

		// we'll handle the enter, tab, and ctrl+n keys to focus the next input
		case tea.KeyEnter, tea.KeyTab, tea.KeyCtrlN:
			// we only really want to check whether the user has provided a To and From address.
			// subject and body can be empty as the email can be sent without them.
			if m.focused == to || m.focused == from {
				m.err = m.validateAddress()
				if m.err != nil {
					return m, nil
				}
			}
			m.nextInput()

		// we'll handle shift+tab to focus the previous input
		case tea.KeyShiftTab:
			m.prevInput()

		// we'll handle ctrl+s to send the message
		case tea.KeyCtrlS:
			// we don't want to send the message if there's an error
			if m.err != nil {
				return m, nil
			}
			m.sendMsg()
			return m, tea.Quit

		// we'll handle ctrl+c to quit the program
		case tea.KeyCtrlC:
			log.Println("Quitting...")
			return m, tea.Quit
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

	if m.err != nil {
		return fmt.Sprintf(`
	%s
	%s
	
	%s
	%s

	%s
	%s

	%s
	%s

	%s`,

			// renders the to header and input
			inputStyle.Width(50).Render("To:"),
			m.inputs[to].View(),

			// renders the from header and input
			inputStyle.Width(50).Render("From:"),
			m.inputs[from].View(),

			// renders the subject header and input
			inputStyle.Width(50).Render("Subject:"),
			m.inputs[subject].View(),

			// renders the body header and input
			inputStyle.Width(50).Render("Body:"),
			m.inputs[body].View(),

			// renders the continue prompt at the bottom of the screen
			continueStyle.Render("(ctrl + c to quit or ctrl + s to send) ->")) + "\n" + m.err.Error() + "\n"
	}

	return fmt.Sprintf(`
	%s
	%s
	
	%s
	%s

	%s
	%s

	%s
	%s

	%s`,

		// renders the to header and input
		inputStyle.Width(50).Render("To:"),
		m.inputs[to].View(),

		// renders the from header and input
		inputStyle.Width(50).Render("From:"),
		m.inputs[from].View(),

		// renders the subject header and input
		inputStyle.Width(50).Render("Subject:"),
		m.inputs[subject].View(),

		// renders the body header and input
		inputStyle.Width(50).Render("Body:"),
		m.inputs[body].View(),

		// renders the continue prompt at the bottom of the screen
		continueStyle.Render("(ctrl + c to quit or ctrl + s to send) ->")) + "\n"
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

func (m *model) sendMsg() {
	log.Println("Sending message...")
}
