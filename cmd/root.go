// Copyright © 2017 Michael Shindle <mshindle@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/mshindle/logext"
	"github.com/mshindle/objstore/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
	"syscall"
)

const (
	cfgEngine = "engine"
	cfgLogFile = "log"
	cfgPort = "port"
	cfgProcs   = "procs"
)

var (
	cfgFile  string
	settings *server.Settings
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "objstore",
	Short: "store arbitrary objects for a given key",
	Long: `Objstore is a simplified HTTP ops interface
which runs on top of various ops implementations like local file system,
S3, or SwiftStack.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		runtime.GOMAXPROCS(viper.GetInt(cfgProcs))
		initLogging()
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.objstore.yaml)")

	RootCmd.PersistentFlags().StringP(cfgLogFile, "l", "", "send logging information to file instead of stdout")
	viper.BindPFlag(cfgLogFile, RootCmd.PersistentFlags().Lookup(cfgLogFile))

	RootCmd.PersistentFlags().Int(cfgProcs, runtime.NumCPU(), "max number of processors to use")
	viper.BindPFlag(cfgProcs, RootCmd.PersistentFlags().Lookup(cfgProcs))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".objstore") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	settings = &server.Settings{}
	err := viper.Unmarshal(settings)
	if err != nil {
		logrus.WithError(err).Fatal("error processing settings")
	}
}

func initLogging() {
	// set the logging level
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	// build our log writer
	logFile := viper.GetString(cfgLogFile)
	if logFile != "" {
		cw := logext.NewCycleWriter(logFile)
		if cw == nil {
			logrus.WithField("logFile", logFile).Error("could not open for writing - using stdout instead")
			return
		}
		cw.OnSignal(syscall.SIGHUP)
		logrus.SetOutput(cw)
	}
}
