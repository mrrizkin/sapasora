// Package apikeys provides the APIKeys command
package apikeys

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sapasora/internal/app"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/permission"

	"sapasora/internal/modules/apikey"
	"sapasora/platform/support/console"
	"sapasora/platform/support/hash"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// APIKeysCmd represents the apikeys command
var APIKeysCmd = &cobra.Command{
	Use:   "apikeys",
	Short: "API keys management tools",
}

// CreateAPIKeyCmd represents the command to create a new API key
var CreateAPIKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new super API key for authenticated user",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			console.Error("Error running program: %s", err.Error())
			return
		}
	},
}

func init() {
	APIKeysCmd.AddCommand(CreateAPIKeyCmd)
}

// Styles for the form
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

// InputField represents a single input field in the form
type InputField struct {
	textInput textinput.Model
	label     string
	err       error
}

// Model represents the form model
type model struct {
	focusIndex int
	inputs     []InputField
	creating   bool
	success    bool
	message    string
	apiKey     string
}

// initialModel returns the initial form model
func initialModel() model {
	m := model{
		inputs: make([]InputField, 2),
	}

	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "johndoe"
	usernameInput.Focus()
	usernameInput.PromptStyle = focusedStyle
	usernameInput.TextStyle = focusedStyle
	m.inputs[0] = InputField{
		textInput: usernameInput,
		label:     "Username",
	}

	// Password input
	passwordInput := textinput.New()
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Placeholder = "password"
	passwordInput.PromptStyle = noStyle
	passwordInput.TextStyle = noStyle
	m.inputs[1] = InputField{
		textInput: passwordInput,
		label:     "Password",
	}

	return m
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit the form
			return m.createAPIKey()
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevFocus()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextFocus()
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}

		// Handle character input and other keys for the focused input
		cmd := m.updateInputs(msg)
		return m, cmd

	case createCompleteMsg:
		m.creating = false
		m.success = true
		m.apiKey = msg.apiKey
		m.message = "Super API key created successfully!"
		return m, tea.Tick(time.Duration(5000)*time.Millisecond, func(_ time.Time) tea.Msg {
			return timeoutMsg{}
		})

	case createErrorMsg:
		m.creating = false
		m.success = false
		m.message = fmt.Sprintf("Error creating API key: %s", msg.err.Error())
		return m, tea.Tick(time.Duration(3000)*time.Millisecond, func(_ time.Time) tea.Msg {
			return timeoutMsg{}
		})

	case timeoutMsg:
		return m, tea.Quit

	}

	return m, nil
}

// View renders the form UI
func (m model) View() string {
	if m.creating {
		return "\nAuthenticating and creating API key...\n"
	}

	if m.message != "" {
		if m.success {
			var s strings.Builder
			s.WriteString("\n")
			s.WriteString(successStyle.Render("✓ " + m.message))
			s.WriteString("\n\n")
			s.WriteString(lipgloss.NewStyle().Bold(true).Render("Your Super API Key:"))
			s.WriteString("\n")
			s.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("43")).
				Background(lipgloss.Color("235")).
				Padding(0, 1).
				Render(m.apiKey))
			s.WriteString("\n\n")
			s.WriteString(helpStyle.Render("⚠️  Save this key securely. It won't be shown again."))
			s.WriteString("\n\n")
			return s.String()
		} else {
			return fmt.Sprintf("\n%s\n", errorStyle.Render("✗ "+m.message))
		}
	}

	var s strings.Builder
	s.WriteString("\nCreate New Super API Key\n\n")
	s.WriteString(helpStyle.Render("Authenticate with your credentials to generate a new API key"))
	s.WriteString("\n\n")

	for i := range m.inputs {
		input := &m.inputs[i]
		focused := i == m.focusIndex

		if focused {
			input.textInput.PromptStyle = focusedStyle
			input.textInput.TextStyle = focusedStyle
		} else {
			input.textInput.PromptStyle = noStyle
			input.textInput.TextStyle = noStyle
		}

		if input.err != nil {
			input.textInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			input.textInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		}

		fmt.Fprintf(&s, "%s: %s", input.label, input.textInput.View())
		if input.err != nil {
			fmt.Fprintf(&s, "\n  %s", helpStyle.Render(input.err.Error()))
		}
		s.WriteString("\n")
	}

	s.WriteString("\nPress Tab to navigate, Enter to submit, Esc/Ctrl+C to quit.\n")

	return s.String()
}

// Helper methods
func (m *model) nextFocus() {
	m.focusIndex++
	if m.focusIndex >= len(m.inputs) {
		m.focusIndex = 0
	}
	m.updateFocus()
}

func (m *model) prevFocus() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs) - 1
	}
	m.updateFocus()
}

func (m *model) updateFocus() {
	for i := range m.inputs {
		m.inputs[i].textInput.Blur()
	}
	m.inputs[m.focusIndex].textInput.Focus()
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.inputs[m.focusIndex].textInput, cmd = m.inputs[m.focusIndex].textInput.Update(msg)
	return cmd
}

// Message types for async operations
type createCompleteMsg struct {
	apiKey string
}

type createErrorMsg struct {
	err error
}

type timeoutMsg struct{}

// createAPIKey handles the actual creation of the API key
func (m model) createAPIKey() (model, tea.Cmd) {
	// Validate inputs
	username := strings.TrimSpace(m.inputs[0].textInput.Value())
	password := m.inputs[1].textInput.Value()

	// Basic validation
	hasError := false
	if username == "" {
		m.inputs[0].err = fmt.Errorf("username is required")
		hasError = true
	} else {
		m.inputs[0].err = nil
	}

	if password == "" {
		m.inputs[1].err = fmt.Errorf("password is required")
		hasError = true
	} else {
		m.inputs[1].err = nil
	}

	if hasError {
		return m, nil
	}

	// Set creating state
	newM := m
	newM.creating = true

	// Run the creation in a goroutine
	return newM, func() tea.Msg {
		// Bootstrap the app to access services
		var shutdown fx.Shutdowner
		var createErr error
		var generatedKey string

		app.Bootstrap(
			fx.Populate(
				&shutdown,
			),
			app.Module,
			app.FxLogger,
			app.Start(
				func(accService account.AccountService,
					apiKeyService apikey.APIKeyService,
				) {
					// Authenticate user
					existingAcc, err := accService.GetAccountByUsername(
						context.Background(),
						username,
					)
					if err != nil || existingAcc == nil {
						createErr = fmt.Errorf("invalid username or password")
						if err := shutdown.Shutdown(); err != nil {
							console.Error("Error shutting down application: %s", err.Error())
						}

						return
					}

					if !hash.Argon2Verify(password, existingAcc.HashedPassword) {
						createErr = fmt.Errorf("invalid username or password")
						if err := shutdown.Shutdown(); err != nil {
							console.Error("Error shutting down application: %s", err.Error())
						}

						return
					}

					generatedKey = "sk-dak-" + hash.GenerateNanoID(
						"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
						42,
					)

					newAPIKey := apikey.APIKey{
						PublicID: hash.NanoID(),
						Name:     "Super API Key",
						Key:      generatedKey,
						Status:   apikey.APIKeyStatusActive,
						Permissions: &permission.Permission{
							// devices
							CanListDevice:   true,
							CanGetDevice:    true,
							CanUpdateDevice: true,
							CanStoreDevice:  true,
							CanDeleteDevice: true,

							// accounts
							CanListAccount:           true,
							CanGetAccount:            true,
							CanUpdateAccount:         true,
							CanUpdatePasswordAccount: true,
							CanStoreAccount:          true,
							CanDeleteAccount:         true,

							// roles
							CanListRole:             true,
							CanGetRole:              true,
							CanUpdateRole:           true,
							CanUpdatePermissionRole: true,
							CanStoreRole:            true,
							CanDeleteRole:           true,

							// apikeys
							CanListAPIKey:   true,
							CanGetAPIKey:    true,
							CanUpdateAPIKey: true,
							CanStoreAPIKey:  true,
							CanDeleteAPIKey: true,

							// gateways
							CanCheckUserGateway:        true,
							CanConnectGateway:          true,
							CanDisconnectGateway:       true,
							CanGetAvatarGateway:        true,
							CanGetContactsGateway:      true,
							CanGetQRGateway:            true,
							CanGetStatusGateway:        true,
							CanGetUserGateway:          true,
							CanLogoutGateway:           true,
							CanSendAudioGateway:        true,
							CanSendButtonGateway:       true,
							CanSendChatPresenceGateway: true,
							CanSendContactGateway:      true,
							CanSendDocumentGateway:     true,
							CanSendImageGateway:        true,
							CanSendListGateway:         true,
							CanSendLocationGateway:     true,
							CanSendStickerGateway:      true,
							CanSendTextGateway:         true,
							CanSendVideoGateway:        true,

							// devicetokens
							CanListDeviceToken:   true,
							CanGetDeviceToken:    true,
							CanStoreDeviceToken:  true,
							CanDeleteDeviceToken: true,
						},
						UserID: existingAcc.ID,
					}

					if err := apiKeyService.CreateAPIKey(context.Background(), &newAPIKey); err != nil {
						createErr = fmt.Errorf("error creating API key: %w", err)
						if err := shutdown.Shutdown(); err != nil {
							console.Error("Error shutting down application: %s", err.Error())
						}
						return
					}

					console.Info("API key created successfully for user: %s", username)

					if err := shutdown.Shutdown(); err != nil {
						console.Error("Error shutting down application: %s", err.Error())
					}
				},
			),
		)

		// Shutdown the app
		defer func() {
			if err := shutdown.Shutdown(); err != nil {
				console.Error("Error shutting down application: %s", err.Error())
			}
		}()

		if createErr != nil {
			return createErrorMsg{err: createErr}
		}

		return createCompleteMsg{apiKey: generatedKey}
	}
}
