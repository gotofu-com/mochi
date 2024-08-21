/*
Copyright © 2024-present The Mochi Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"log"
	"log/slog"
	"os"

	"gotofu.com/mochi/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mochi",
	Short: "A tool for managing release notes in a monorepo",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
}
