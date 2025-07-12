package cmd

import (
	"fmt"
	"time"

	"github.com/alizmhdi/shamsi-calendar/calendar"

	"github.com/spf13/cobra"
)

const (
	minYear  = 1
	maxYear  = 9999
	minMonth = 1
	maxMonth = 12
)

var (
	yearFlag     int
	monthFlag    int
	threeFlag    bool
	fullYearFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "scal",
	Short: "Display a Jalali (Shamsi) calendar",
	Long: `A command line tool to display Jalali (Shamsi) calendar, similar to the Unix 'cal' command.

Features:
- Display current month calendar
- Display specific month/year
- Display entire year
- Display three months
- Highlight today's date`,
	RunE: runCalendar,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().IntVarP(&yearFlag, "year", "y", 0, "year to display (default: current year)")
	rootCmd.Flags().IntVarP(&monthFlag, "month", "m", 0, "month to display (1-12, default: current month)")
	rootCmd.Flags().BoolVarP(&threeFlag, "three", "3", false, "display three months spanning the date")
	rootCmd.Flags().BoolVarP(&fullYearFlag, "full-year", "Y", false, "display entire year")
}

func validateInput(year, month int) error {
	if month < minMonth || month > maxMonth {
		return fmt.Errorf("month must be between %d and %d", minMonth, maxMonth)
	}

	if year < minYear || year > maxYear {
		return fmt.Errorf("year must be between %d and %d", minYear, maxYear)
	}

	return nil
}

// getCurrentJalaliDate returns the current date in Jalali calendar
func getCurrentJalaliDate() calendar.JalaliDate {
	now := time.Now()
	return calendar.GregorianToJalali(now.Year(), int(now.Month()), now.Day())
}

// determineDisplayMode determines which display mode to use based on flags
func determineDisplayMode(cmd *cobra.Command) displayMode {
	yearFlagSet := cmd.Flags().Changed("year")
	monthFlagSet := cmd.Flags().Changed("month")

	if fullYearFlag {
		return modeFullYear
	}
	if threeFlag {
		return modeThreeMonths
	}
	if yearFlagSet && !monthFlagSet {
		return modeFullYear
	}
	return modeSingleMonth
}

type displayMode int

const (
	modeSingleMonth displayMode = iota
	modeThreeMonths
	modeFullYear
)

func runCalendar(cmd *cobra.Command, args []string) error {
	// Get current Jalali date for defaults and today highlighting
	currentJalali := getCurrentJalaliDate()

	// Set default values if not provided
	if yearFlag == 0 {
		yearFlag = currentJalali.Year
	}
	if monthFlag == 0 {
		monthFlag = currentJalali.Month
	}

	if err := validateInput(yearFlag, monthFlag); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Determine display mode and execute
	mode := determineDisplayMode(cmd)

	switch mode {
	case modeFullYear:
		calendar.DisplayYearTable(yearFlag)
	case modeThreeMonths:
		calendar.DisplayThreeMonthsTable(yearFlag, monthFlag)
	case modeSingleMonth:
		calendar.DisplayMonthTable(yearFlag, monthFlag, currentJalali)
	default:
		return fmt.Errorf("unknown display mode")
	}

	return nil
}
