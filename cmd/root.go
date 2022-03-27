package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xy3/aliens/config"
	"github.com/xy3/aliens/parser"
	"github.com/xy3/aliens/simulation"
	"os"
	"strconv"
	"sync"
)

// Channels for messages
// Packaging
// Using IO writer for files and logs
// Controlling the simulation in the cmd package
// TODO Update tests
// TODO Update Comments
// New for command

func runSimulation(count int, logger *log.Logger) error {
	mapFile, err := os.Open(config.Config.MapFile)
	if err != nil {
		return err
	}
	defer mapFile.Close()

	worldMap, err := parser.ParseMap(mapFile)
	if err != nil {
		return err
	}

	logger.Info("PARSED SIMULATION WORLD MAP")
	fmt.Println(worldMap.Serialize())

	namesFile, err := os.Open(config.Config.AlienNamesFile)
	if err != nil {
		return err
	}
	defer namesFile.Close()

	allAliens, err := simulation.RandomAliens(count, namesFile)
	if err != nil {
		return err
	}
	logger.Infof("CREATED %d RANDOM ALIENS SUCCESSFULLY", count)

	sim := simulation.New(allAliens, worldMap, config.Config.MaxAlienMoves)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sim.LogWorker(ctx, logger)
	}()
	sim.DeployAliens()
	result := sim.Run()
	cancel()
	wg.Wait()

	result.Display(logger)

	logger.Info("The world map that still remains is:")
	fmt.Println(result.WorldMap.Serialize())

	return nil
}

func NewRootCmd(logger *log.Logger) *cobra.Command {
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
			if config.Config.DebugMode {
				logger.SetLevel(log.DebugLevel)
			}
			configPath := config.Path()
			configFile, err := os.Open(configPath)
			if err != nil {
				logger.WithError(err).Warnf("Failed to open the config file [%s] for writing", configPath)
			}
			_ = config.Write(configFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			count, err := strconv.Atoi(args[0])
			if err != nil {
				logger.Fatal(err)
			}
			err = runSimulation(count, logger)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringVarP(
		&config.Config.MapFile,
		"map",
		"m",
		config.Config.MapFile,
		"text file with cities and routes for the simulation",
	)
	rootCmd.Flags().StringVarP(
		&config.Config.AlienNamesFile,
		"alien-names",
		"n",
		config.Config.AlienNamesFile,
		"text file with names of simulation aliens on each line",
	)
	rootCmd.Flags().BoolVarP(
		&config.Config.DebugMode,
		"debug",
		"d",
		config.Config.DebugMode,
		"debug mode to view more logging",
	)
	return rootCmd
}
