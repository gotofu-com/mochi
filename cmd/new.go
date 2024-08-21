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
	"fmt"

	"gotofu.com/mochi/change"
	"gotofu.com/mochi/change_type"
	"gotofu.com/mochi/config"
	"gotofu.com/mochi/domain"
	"gotofu.com/mochi/target"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [type] [target] [message]",
	Short: "Create a new release note entry",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string

		switch len(args) {
		case 0:
			for _, t := range config.Configuration.Types {
				comps = append(comps, t.Id)
			}
		case 1:
			for _, t := range config.Configuration.Targets {
				comps = append(comps, t.Id)
			}
		case 2:
			break
		case 3:
			comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
		default:
			comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
		}

		return comps, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			c   domain.Change
			err error
		)

		if len(args) > 0 {
			if c.Type, err = change_type.Get(args[0]); err != nil {
				return err
			}
		} else {
			prompt := promptui.Select{
				Label:     "Choose the type of change",
				Items:     config.Configuration.Types,
				Templates: namedItemPromptTemplate,
			}

			if index, _, err := prompt.Run(); err != nil {
				return err
			} else {
				c.Type = &config.Configuration.Types[index]
			}
		}

		if len(args) > 1 {
			if c.Target, err = target.Get(args[1]); err != nil {
				return err
			}
		} else {
			prompt := promptui.Select{
				Label:     "Choose the target affected by the change",
				Items:     config.Configuration.Targets,
				Templates: namedItemPromptTemplate,
			}

			if index, _, err := prompt.Run(); err != nil {
				return err
			} else {
				c.Target = &config.Configuration.Targets[index]
			}
		}

		if len(args) > 2 {
			c.Message = args[2]
		} else {
			prompt := promptui.Prompt{
				Label: "Enter the message for the change",
			}

			if c.Message, err = prompt.Run(); err != nil {
				return err
			}
		}

		if err := change.Commit(&c); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

var namedItemPromptTemplate = &promptui.SelectTemplates{
	Label:    fmt.Sprintf("%s {{ .Name }}: ", promptui.IconInitial),
	Active:   fmt.Sprintf("%s {{ .Name | underline }}", promptui.IconSelect),
	Inactive: "  {{ .Name }}",
	Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}`, promptui.IconGood),
}
