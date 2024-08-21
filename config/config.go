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

package config

import (
	"errors"
	"fmt"

	"gotofu.com/mochi/domain"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Types      []domain.ChangeType
	Targets    []domain.Target
	BaseBranch string
}

var Configuration *Config

func InitConfig() {
	viper.AddConfigPath("./.mochi")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("baseBranch", "main")
	viper.SetDefault("types", []domain.ChangeType{
		{Id: "feature", Name: "Feature", Title: "Features"},
		{Id: "bugfix", Name: "Bug fix", Title: "Bug Fixes"},
		{Id: "doc", Name: "Documentation", Title: "Documentation"},
		{Id: "removal", Name: "Removal", Title: "Removals"},
		{Id: "misc", Name: "Miscellaneous", Title: "Miscellaneous"},
	})
	viper.SetDefault("targets", []domain.Target{})

	viper.BindEnv("githubToken", "GITHUB_TOKEN")

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		cobra.CheckErr(err)
	}

	if err := viper.Unmarshal(&Configuration); err != nil {
		cobra.CheckErr(fmt.Errorf("unable to decode into struct, %v", err))
	}
}
