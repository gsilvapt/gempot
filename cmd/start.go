package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tracking time",
	Long: `This command lets users start tracking time to a specific project. It contains flags to track time to specific 
	projects, or just fall into a bucket of all. Ideally, use a "project" in any case to facilitate later filtering and 
	searching. Example:
	
	gempot start -p "my-project"

	To stop recording time, just hit ctrl-c to leave the tracker and it will add the entry to the chosen output 
	(default is gempot.csv in your home directory)

	`,
	Run: startTracker,
}

type timeDiff struct {
	hours   int
	minutes int
	seconds int
}

func startTracker(cmd *cobra.Command, args []string) {
	filename := viper.GetString("output")
	projectId, _ := cmd.Flags().GetString("project")

	Logger.Info(fmt.Sprintf("Starting time tracker for project %s", projectId))
	startTs := time.Now()
	currentTs := time.Now()
	run := true

	go func() {
		// Setup listener for SIGINT, which will be used by the timer.
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		Logger.Info(fmt.Sprintf("Shutting down tracker to project %s", projectId))
		run = false
	}()

	for range time.Tick(1 * time.Second) {
		currentTs = time.Now()
		diff := getTimeDiff(startTs, currentTs)

		Logger.Info(fmt.Sprintf("Hours: %02d; Minutes: %02d; Seconds: %02d", diff.hours, diff.minutes, diff.seconds))
		if !run {
			break
		}
	}

	totalDiff := getTimeDiff(startTs, currentTs)
	writeToCsv(filename, projectId, startTs, currentTs, totalDiff)
}

func getTimeDiff(start time.Time, end time.Time) *timeDiff {
	diff := end.Sub(start)
	total := int(diff.Seconds())
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return &timeDiff{
		hours:   hours,
		minutes: minutes,
		seconds: seconds,
	}
}

func prepCmdFlags() {
	startCmd.Flags().String("project", "", "provide a project unique identifier to start tracking task time")
}

func init() {
	rootCmd.AddCommand(startCmd)
	prepCmdFlags()
}

func writeToCsv(filename, projectId string, startTs, currentTs time.Time, diff *timeDiff) {
	homedir, _ := os.UserHomeDir()
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", homedir, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Error("failed opening output file to write to it")
	}

	startTs.Format(time.RFC3339)
	currentTs.Format(time.RFC3339)

	diffHours := diff.hours
	if diff.minutes > 30 {
		diffHours += 1
	}

	w := csv.NewWriter(file)
	if err := w.Write([]string{projectId, startTs.String(), currentTs.String(), strconv.Itoa(diffHours)}); err != nil {
		Logger.Error(fmt.Sprintf("failed writing to output file: %e\n", err))
	}
}
