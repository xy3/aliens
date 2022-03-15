package main

import (
	"github.com/xy3/aliens"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var mapFile = "map.txt"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aliens [alien_count]",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = aliens.RunSimulation(count, mapFile)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&mapFile, "map", "m", "map.txt", "map file with cities for the simulation")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
