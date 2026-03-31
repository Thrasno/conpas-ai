package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the ASCII art for the Conpas AI stylized N logo.
var logoLines = []string{
	"                                                      ",
	"         ███╗   ██╗                                   ",
	"         ████╗  ██║                                   ",
	"         ██╔██╗ ██║                                   ",
	"         ██║╚██╗██║                                   ",
	"         ██║ ╚████║                                   ",
	"         ╚═╝  ╚═══╝                                   ",
	"                                                      ",
}

// gradientColors defines the top-to-bottom gradient for the logo.
// Distributed across rows: Green/Cyan → Teal → Blue → Lavender → Mauve (cyan to purple).
var gradientColors = []lipgloss.Color{
	ColorGreen,    // band 1 (cyan/green top)
	ColorTeal,     // band 2
	ColorBlue,     // band 3 (middle blue)
	ColorLavender, // band 4
	ColorMauve,    // band 5 (purple bottom)
}

// RenderLogo returns the braille ASCII logo with a top-to-bottom gradient.
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		if i < total-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
