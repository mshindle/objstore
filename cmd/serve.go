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
	"github.com/mshindle/objstore/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "listen and serve requests",
	Long: `Run the webserver to handle ops requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server.ListenAndServe(settings)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP(cfgEngine, "e", "local", "ops engine to use")
	viper.BindPFlag(cfgEngine, serveCmd.Flags().Lookup(cfgEngine))

	serveCmd.Flags().IntP(cfgPort, "p", 8080, "port to listen at")
	viper.BindPFlag(cfgPort, serveCmd.Flags().Lookup(cfgPort))
}
