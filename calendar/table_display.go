package calendar

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	// colors
	todayColor  = "\033[1;33m" // bold yellow for today's date
	headerColor = "\033[1;36m" // bold cyan for month/year header
	resetColor  = "\033[0m"

	// calendar constants
	daysInWeek      = 7
	monthsInYear    = 12
	quartersInYear  = 4
	monthsInQuarter = 3
)


var monthNames = []string{
	"Farvardin", "Ordibehesht", "Khordad", "Tir", "Mordad", "Shahrivar",
	"Mehr", "Aban", "Azar", "Dey", "Bahman", "Esfand",
}

var dayNames = []string{"Shanbe", "Yek", "Do", "Se", "Chahar", "Panj", "Jome"}

// stripANSI removes ANSI color codes from a string for accurate width calculation
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

// createTable creates a new table with common configuration
func createTable() (*tablewriter.Table, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)

	table.SetHeader(dayNames)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	// Set header colors (bright white)
	headerColors := make([]tablewriter.Colors, daysInWeek)
	for i := range headerColors {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiWhiteColor, tablewriter.Bold}
	}
	table.SetHeaderColor(headerColors...)

	return table, buf
}

// formatDay formats a day number with optional highlighting for today
func formatDay(day int, isToday bool) string {
	if day == 0 {
		return ""
	}

	dayStr := strconv.Itoa(day)
	if isToday {
		return todayColor + dayStr + resetColor
	}
	return dayStr
}

// calculateTableWidth calculates the maximum width of table lines (excluding ANSI codes)
func calculateTableWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		cleanLine := stripANSI(line)
		if len(cleanLine) > maxWidth {
			maxWidth = len(cleanLine)
		}
	}
	return maxWidth
}

// centerText centers text within a given width
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

// renderMonthAsLines renders a single month as a slice of strings, with colored header and today highlight
func renderMonthAsLines(year, month int, currentDate JalaliDate) []string {
	calendar := GetMonthCalendar(year, month)

	table, buf := createTable()

	// Add calendar rows
	for _, week := range calendar {
		row := make([]string, daysInWeek)
		for i, day := range week {
			isToday := day == currentDate.Day && month == currentDate.Month && year == currentDate.Year
			row[i] = formatDay(day, isToday)
		}
		table.Append(row)
	}

	table.Render()
	tableLines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")

	// Calculate table width and center month header
	tableWidth := calculateTableWidth(tableLines)
	monthHeader := centerText(monthNames[month-1], tableWidth)
	monthHeaderLine := headerColor + monthHeader + resetColor

	// Compose the final lines
	lines := []string{monthHeaderLine}
	lines = append(lines, tableLines...)
	return lines
}

// DisplayMonthTable displays a single month calendar using tablewriter
func DisplayMonthTable(year, month int, currentDate JalaliDate) {
	calendar := GetMonthCalendar(year, month)

	table, buf := createTable()

	// Add calendar rows
	for _, week := range calendar {
		row := make([]string, daysInWeek)
		for i, day := range week {
			isToday := day == currentDate.Day && month == currentDate.Month && year == currentDate.Year
			row[i] = formatDay(day, isToday)
		}
		table.Append(row)
	}

	table.Render()
	tableOutput := buf.String()
	tableLines := strings.Split(tableOutput, "\n")

	// Calculate table width and center header
	tableWidth := calculateTableWidth(tableLines)
	centeredHeader := fmt.Sprintf("%s %d", monthNames[month-1], year)
	centeredHeader = centerText(centeredHeader, tableWidth)

	fmt.Printf("%s%s%s\n", headerColor, centeredHeader, resetColor)
	fmt.Print(tableOutput)
}

// getAdjacentMonths calculates the previous and next months for a given month/year
func getAdjacentMonths(year, month int) (prevYear, prevMonth, nextYear, nextMonth int) {
	prevMonth = month - 1
	prevYear = year
	if prevMonth < 1 {
		prevMonth = monthsInYear
		prevYear--
	}

	nextMonth = month + 1
	nextYear = year
	if nextMonth > monthsInYear {
		nextMonth = 1
		nextYear++
	}

	return prevYear, prevMonth, nextYear, nextMonth
}

// padMonthLines ensures all month lines have the same height and consistent width
func padMonthLines(monthLines [][]string, maxLines int) {
	for i := range monthLines {
		// Find the maximum width for this month
		maxWidth := 0
		for _, line := range monthLines[i] {
			cleanLine := stripANSI(line)
			if len(cleanLine) > maxWidth {
				maxWidth = len(cleanLine)
			}
		}

		// Pad each line to the maximum width
		for j := range monthLines[i] {
			cleanLine := stripANSI(monthLines[i][j])
			padding := maxWidth - len(cleanLine)
			monthLines[i][j] = monthLines[i][j] + strings.Repeat(" ", padding)
		}

		// Pad to same height
		for len(monthLines[i]) < maxLines {
			monthLines[i] = append(monthLines[i], strings.Repeat(" ", maxWidth))
		}
	}
}

// DisplayThreeMonthsTable displays three months using colored, aligned tables
func DisplayThreeMonthsTable(year, month int) {
	now := time.Now()
	currentJalali := GregorianToJalali(now.Year(), int(now.Month()), now.Day())

	// Calculate previous and next months
	prevYear, prevMonth, nextYear, nextMonth := getAdjacentMonths(year, month)

	// Render three months as lines
	monthLines := make([][]string, 3)
	maxLines := 0

	// Previous month
	monthLines[0] = renderMonthAsLines(prevYear, prevMonth, currentJalali)
	if len(monthLines[0]) > maxLines {
		maxLines = len(monthLines[0])
	}

	// Current month
	monthLines[1] = renderMonthAsLines(year, month, currentJalali)
	if len(monthLines[1]) > maxLines {
		maxLines = len(monthLines[1])
	}

	// Next month
	monthLines[2] = renderMonthAsLines(nextYear, nextMonth, currentJalali)
	if len(monthLines[2]) > maxLines {
		maxLines = len(monthLines[2])
	}

	// Pad months to same height and ensure consistent width
	padMonthLines(monthLines, maxLines)

	// Print side by side with consistent spacing
	for line := 0; line < maxLines; line++ {
		fmt.Printf("%s  %s  %s\n", monthLines[0][line], monthLines[1][line], monthLines[2][line])
	}
}

// calculateQuarterWidth calculates the total width of a quarter (3 months)
func calculateQuarterWidth(allMonthLines [][]string, quarter int) int {
	quarterWidth := 0
	for i := 0; i < monthsInQuarter; i++ {
		monthIdx := quarter*monthsInQuarter + i
		monthWidth := 0
		for _, line := range allMonthLines[monthIdx] {
			cleanLine := stripANSI(line)
			if len(cleanLine) > monthWidth {
				monthWidth = len(cleanLine)
			}
		}
		quarterWidth += monthWidth + 2 // +2 for spacing between months
	}
	return quarterWidth
}

// DisplayYearTable displays the entire year using colored, aligned tables
func DisplayYearTable(year int) {
	now := time.Now()
	currentJalali := GregorianToJalali(now.Year(), int(now.Month()), now.Day())

	// First, render all months to calculate the total width
	allMonthLines := make([][]string, monthsInYear)
	maxLines := 0
	for i := 0; i < monthsInYear; i++ {
		month := i + 1
		lines := renderMonthAsLines(year, month, currentJalali)
		allMonthLines[i] = lines
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	// Calculate total width for centering the year
	totalWidth := 0
	for quarter := 0; quarter < quartersInYear; quarter++ {
		quarterWidth := calculateQuarterWidth(allMonthLines, quarter)
		if quarterWidth > totalWidth {
			totalWidth = quarterWidth
		}
	}

	// Center and print the year
	yearStr := fmt.Sprintf("%d", year)
	yearPadding := (totalWidth - len(yearStr)) / 2
	if yearPadding < 0 {
		yearPadding = 0
	}
	fmt.Printf("%s%s%s%s\n\n", strings.Repeat(" ", yearPadding), headerColor, yearStr, resetColor)

	// Display each quarter
	for quarter := 0; quarter < quartersInYear; quarter++ {
		// Get three months for this quarter
		monthLines := make([][]string, monthsInQuarter)
		for i := 0; i < monthsInQuarter; i++ {
			monthIdx := quarter*monthsInQuarter + i
			monthLines[i] = allMonthLines[monthIdx]
		}

		// Pad months to same height and ensure consistent width
		padMonthLines(monthLines, maxLines)

		// Print side by side with consistent spacing
		for line := 0; line < maxLines; line++ {
			fmt.Printf("%s  %s  %s\n", monthLines[0][line], monthLines[1][line], monthLines[2][line])
		}
		fmt.Println()
	}
}
