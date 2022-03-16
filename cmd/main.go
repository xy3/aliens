package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens"
	"strconv"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for this simulation program and executes in the main function
var rootCmd = &cobra.Command{
	Use:   "aliens [alien_count]",
	Args:  cobra.ExactArgs(1),
	Short: "Alien Invasion Simulator",
	Long: `Simulate an alien invasion using a configurable amount of aliens and a custom map.

You may use the -h flag to view help for this program. To modify the default config,
modify 'config.json' in the program's working directory.
Author:
Theodore Coyne Morgan | March 2022`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if aliens.Config.DebugMode {
			log.SetLevel(log.DebugLevel)
		}
		_ = aliens.WriteConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = aliens.RunSimulation(count, aliens.Config.MapFile)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(
		&aliens.Config.MapFile,
		"map",
		"m",
		aliens.Config.MapFile,
		"text file with cities and routes for the simulation",
	)
	rootCmd.Flags().StringVarP(
		&aliens.Config.AlienNamesFile,
		"alien-names",
		"n",
		aliens.Config.AlienNamesFile,
		"text file with names of simulation aliens on each line",
	)
	rootCmd.Flags().BoolVarP(
		&aliens.Config.DebugMode,
		"debug",
		"d",
		aliens.Config.DebugMode,
		"debug mode to view more logging",
	)
}

// main loads the config and executes the rootCmd with flags
func main() {
	err := aliens.LoadConfig()
	if err != nil {
		log.Println("Failed to find or read a 'config.json' file, using default config values or flag overrides")
	}
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
