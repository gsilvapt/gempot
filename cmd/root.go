package cmd

import (
	"github.com/gsilvapt/gempot/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VERSION string = "0.1"

var (
	output string
	Logger *logger.Log
)

var rootCmd = &cobra.Command{
	Use:   `gempot`,
	Short: "gempot is a time tracking tool to track time spent on specific projects or clients",
	Long: `gempot was designed for those interesting in keeping track of time spent on projects and/or clients, allowing 
			them to filter based on project, time filters and also report paid amounts so they are no longer considered.`,
	Version: VERSION,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	Logger = logger.InitLogger()

	rootCmd.PersistentFlags().StringVar(&output, "output", "gempot.csv", "output specifies location of tracker, default is output.csv in home directory")

	rootCmd.InitDefaultVersionFlag()

	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	// TODO: validate whether output exists and create header in it
}
