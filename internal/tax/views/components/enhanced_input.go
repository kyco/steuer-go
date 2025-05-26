package components

import (
	"fmt"
	"strconv"
	"strings"

	"tax-calculator/internal/tax/views/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ValidationRule func(string) error

type EnhancedInput struct {
	textinput.Model
	Label           string
	Description     string
	Placeholder     string
	ValidationRules []ValidationRule
	ErrorMessage    string
	IsValid         bool
	Required        bool
}

func NewEnhancedInput(label, description, placeholder string) EnhancedInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.FgColor)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.AccentColor)

	return EnhancedInput{
		Model:       ti,
		Label:       label,
		Description: description,
		Placeholder: placeholder,
		IsValid:     true,
		Required:    false,
	}
}

func (ei *EnhancedInput) AddValidation(rule ValidationRule) {
	ei.ValidationRules = append(ei.ValidationRules, rule)
}

func (ei *EnhancedInput) Validate() {
	ei.ErrorMessage = ""
	ei.IsValid = true

	value := strings.TrimSpace(ei.Value())

	if ei.Required && value == "" {
		ei.ErrorMessage = fmt.Sprintf("%s is required", ei.Label)
		ei.IsValid = false
		return
	}

	for _, rule := range ei.ValidationRules {
		if err := rule(value); err != nil {
			ei.ErrorMessage = err.Error()
			ei.IsValid = false
			return
		}
	}
}

func (ei *EnhancedInput) Update(msg tea.Msg) (EnhancedInput, tea.Cmd) {
	var cmd tea.Cmd
	ei.Model, cmd = ei.Model.Update(msg)

	if _, ok := msg.(tea.KeyMsg); ok {
		ei.Validate()
	}

	return *ei, cmd
}

func (ei *EnhancedInput) View() string {
	var builder strings.Builder

	labelStyle := styles.SubtitleStyle
	if !ei.IsValid {
		labelStyle = labelStyle.Foreground(styles.DangerColor)
	}

	builder.WriteString(labelStyle.Render(ei.Label))
	builder.WriteString("\n")

	if ei.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(styles.NeutralColor).
			Italic(true)
		builder.WriteString(descStyle.Render(ei.Description))
		builder.WriteString("\n")
	}

	inputStyle := styles.InputFieldStyle
	if ei.Focused() {
		if ei.IsValid {
			inputStyle = styles.ActiveInputStyle
		} else {
			inputStyle = inputStyle.BorderForeground(styles.DangerColor)
		}
	} else if !ei.IsValid {
		inputStyle = inputStyle.BorderForeground(styles.DangerColor)
	}

	builder.WriteString(inputStyle.Render(ei.Model.View()))

	if !ei.IsValid && ei.ErrorMessage != "" {
		builder.WriteString("\n")
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.DangerColor).
			Bold(true)
		builder.WriteString(errorStyle.Render("⚠ " + ei.ErrorMessage))
	}

	return builder.String()
}

func ValidateIncome(value string) error {
	if value == "" {
		return nil
	}

	income, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("Please enter a valid number")
	}

	if income < 0 {
		return fmt.Errorf("Income cannot be negative")
	}

	if income > 10000000 {
		return fmt.Errorf("Income seems unrealistically high")
	}

	return nil
}

func ValidateYear(value string) error {
	if value == "" {
		return nil
	}

	year, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("Please enter a valid year")
	}

	if year < 2020 || year > 2030 {
		return fmt.Errorf("Year must be between 2020 and 2030")
	}

	return nil
}
