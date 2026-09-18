package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"sapasora/internal/app"
	"sapasora/internal/modules/account"
	"sapasora/internal/modules/permission"
	"sapasora/internal/modules/role"
	"sapasora/platform/support/console"
	"sapasora/platform/support/hash"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// AdminUserCmd represents the admin user creation command
var AdminUserCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin user management tools",
}

// CreateUserCmd represents the command to create a new admin user
var CreateUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new admin user interactively",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			console.Error("Error running program: %s", err.Error())
			return
		}
	},
}

func init() {
	AdminUserCmd.AddCommand(CreateUserCmd)
}

// Styles for the form
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
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
}

// initialModel returns the initial form model
func initialModel() model {
	m := model{
		inputs: make([]InputField, 3),
	}

	// Name input
	nameInput := textinput.New()
	nameInput.Placeholder = "John Doe"
	nameInput.Focus()
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle
	m.inputs[0] = InputField{
		textInput: nameInput,
		label:     "Name",
	}

	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "johndoe"
	usernameInput.PromptStyle = noStyle
	usernameInput.TextStyle = noStyle
	m.inputs[1] = InputField{
		textInput: usernameInput,
		label:     "Username",
	}

	// Password input
	passwordInput := textinput.New()
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Placeholder = "password"
	passwordInput.PromptStyle = noStyle
	passwordInput.TextStyle = noStyle
	m.inputs[2] = InputField{
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
			return m.createAdminUser()
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
		m.message = fmt.Sprintf("Admin user '%s' created successfully!", msg.username)
		return m, tea.Tick(time.Duration(2000)*time.Millisecond, func(_ time.Time) tea.Msg {
			return timeoutMsg{}
		})

	case createErrorMsg:
		m.creating = false
		m.success = false
		m.message = fmt.Sprintf("Error creating admin user: %s", msg.err.Error())
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
		return "\nCreating admin user...\n"
	}

	if m.message != "" {
		if m.success {
			return fmt.Sprintf("\n✓ %s\n", m.message)
		} else {
			return fmt.Sprintf("\n✗ %s\n", m.message)
		}
	}

	var s strings.Builder
	s.WriteString("\nCreate New Admin User (Super Admin Role)\n\n")

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
	username string
}

type createErrorMsg struct {
	err error
}

type timeoutMsg struct{}

// createAdminUser handles the actual creation of the admin user
func (m model) createAdminUser() (model, tea.Cmd) {
	// Validate inputs
	name := strings.TrimSpace(m.inputs[0].textInput.Value())
	username := strings.TrimSpace(m.inputs[1].textInput.Value())
	password := m.inputs[2].textInput.Value()

	// Basic validation
	hasError := false
	if name == "" {
		m.inputs[0].err = fmt.Errorf("name is required")
		hasError = true
	} else {
		m.inputs[0].err = nil
	}

	if username == "" {
		m.inputs[1].err = fmt.Errorf("username is required")
		hasError = true
	} else {
		m.inputs[1].err = nil
	}

	if password == "" {
		m.inputs[2].err = fmt.Errorf("password is required")
		hasError = true
	} else {
		m.inputs[2].err = nil
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
		var accService account.AccountService
		var roleService role.RoleService
		var createErr error

		app.Bootstrap(
			fx.Populate(&shutdown, &accService, &roleService),
			app.Module,
			app.FxLogger,
			app.Start(func(as account.AccountService, rs role.RoleService) {
				accService = as
				roleService = rs

				// Check if user already exists
				existingAcc, err := accService.GetAccountByUsername(context.Background(), username)
				if err == nil && existingAcc != nil {
					createErr = fmt.Errorf("user with username '%s' already exists", username)

					if err := shutdown.Shutdown(); err != nil {
						console.Error("Error shutting down application: %s", err.Error())
					}

					return
				}

				// Get or create super admin role
				var superAdminRole *role.Role
				// Try to find existing super admin role (adjust the criteria based on your role structure)
				// For example, you might look for a role with name "Super Admin" or with specific permissions
				existingRole, err := roleService.GetRoleByName(context.Background(), "Super Admin")
				if err != nil || existingRole == nil {
					// Create super admin role if it doesn't exist
					superAdminRole = &role.Role{
						PublicID:    hash.NanoID(),
						Name:        "Super Admin",
						Description: "Super Administrator with all permissions",
						// Add any other role fields that indicate super admin status
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
					}
					if err := roleService.CreateRole(context.Background(), superAdminRole); err != nil {
						createErr = fmt.Errorf("error creating super admin role: %w", err)

						if err := shutdown.Shutdown(); err != nil {
							console.Error("Error shutting down application: %s", err.Error())
						}

						return
					}
					console.Info("Super Admin role created successfully!")
				} else {
					superAdminRole = existingRole
				}

				// Hash the password
				hashedPassword, err := hash.Argon2(password)
				if err != nil {
					createErr = fmt.Errorf("error hashing password: %w", err)

					if err := shutdown.Shutdown(); err != nil {
						console.Error("Error shutting down application: %s", err.Error())
					}

					return
				}

				// Create the new account with super admin role
				newAccount := &account.Account{
					PublicID:       hash.NanoID(),
					Name:           name,
					Username:       username,
					HashedPassword: hashedPassword,
					RoleID:         superAdminRole.ID,
				}

				if err := accService.CreateAccount(context.Background(), newAccount); err != nil {
					createErr = fmt.Errorf("error creating account: %w", err)

					if err := shutdown.Shutdown(); err != nil {
						console.Error("Error shutting down application: %s", err.Error())
					}

					return
				}

				console.Info("Account created successfully!")

				if err := shutdown.Shutdown(); err != nil {
					console.Error("Error shutting down application: %s", err.Error())

					os.Exit(1)
				}
			}),
		)

		if createErr != nil {
			if err := shutdown.Shutdown(); err != nil {
				console.Error("Error shutting down application: %s", err.Error())
			}
			return createErrorMsg{err: createErr}
		}

		if err := shutdown.Shutdown(); err != nil {
			console.Error("Error shutting down application: %s", err.Error())
			return createErrorMsg{err: err}
		}

		return createCompleteMsg{username: username}
	}
}
