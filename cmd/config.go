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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "dump out the configuration being used",
	Long:  `Prints out the configuration settings used by objstore at runtime`,
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := yaml.Marshal(settings)
		if err != nil {
			logrus.WithError(err).Error("could not marshal settings")
			return err
		}
		fmt.Printf("--- settings dump:\n%s\n\n", string(d))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
