package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Exibe logs da aplicação",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		lines, _ := cmd.Flags().GetInt("lines")
		level, _ := cmd.Flags().GetString("level")
		since, _ := cmd.Flags().GetString("since")
		until, _ := cmd.Flags().GetString("until")
		logPath := filepath.Join(os.Getenv("HOME"), ".fg", "logs", "app.log")

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Printf("No log file found at %s\n", logPath)
			return
		}

		content, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		var sinceTime, untilTime time.Time
		if since != "" {
			sinceTime, err = time.Parse("2006-01-02 15:04:05", since)
			if err != nil {
				fmt.Printf("Invalid since time format. Use YYYY-MM-DD HH:MM:SS\n")
				return
			}
		}
		if until != "" {
			untilTime, err = time.Parse("2006-01-02 15:04:05", until)
			if err != nil {
				fmt.Printf("Invalid until time format. Use YYYY-MM-DD HH:MM:SS\n")
				return
			}
		}

		logLines := strings.Split(string(content), "\n")
		filteredLines := make([]string, 0)

		for _, line := range logLines {
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, " ", 3)
			if len(parts) < 3 {
				continue
			}

			timestamp := parts[0] + " " + parts[1]
			logTime, err := time.Parse("2006-01-02 15:04:05", timestamp)
			if err != nil {
				continue
			}

			if !sinceTime.IsZero() && logTime.Before(sinceTime) {
				continue
			}
			if !untilTime.IsZero() && logTime.After(untilTime) {
				continue
			}

			if level != "" {
				levelPart := strings.Split(parts[2], "]")[0]
				if !strings.Contains(levelPart, level) {
					continue
				}
			}

			filteredLines = append(filteredLines, line)
		}

		start := 0
		if len(filteredLines) > lines {
			start = len(filteredLines) - lines
		}
		for _, line := range filteredLines[start:] {
			fmt.Println(line)
		}
	},
}

func init() {
	logCmd.Flags().IntP("lines", "n", 100, "Número de linhas que será exibido")
	logCmd.Flags().StringP("level", "l", "", "Filter logs by level (DEBUG, INFO, WARN, ERROR)")
	logCmd.Flags().String("since", "", "Mostrar logs desde (format: YYYY-MM-DD HH:MM:SS)")
	logCmd.Flags().String("until", "", "Mostrar logs até (format: YYYY-MM-DD HH:MM:SS)")
	RootCmd.AddCommand(logCmd)
}
