/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/BurntSushi/toml"
	"github.com/hermeschat/engine/config"
	"github.com/hermeschat/engine/monitoring"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models using SqlBoiler",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// generate sqlboiler.toml
		fd, err := os.OpenFile("sqlboiler.toml", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			monitoring.Logger().Fatalf("%s\n", err)
		}
		enc := toml.NewEncoder(fd)
		db := config.C.GetString("database.type")
		host := config.C.GetString("database.host")
		port := config.C.GetString("database.port")
		user := config.C.GetString("database.user")
		password := config.C.GetString("database.password")
		name := config.C.GetString("database.name")

		c := map[string]interface{}{
			db: map[string]interface{}{
				"host":     host,
				"dbname":   name,
				"port":     port,
				"user":     user,
				"password": password,
			},
		}
		err = enc.Encode(c)
		if err != nil {
			monitoring.Logger().Fatalf("%s\n")
		}
		// build sqlboiler command
		cmdArgs := []string{"--wipe"}
		cmdArgs = append(cmdArgs, db)
		// run
		sqlBoilerCmd := exec.Command("sqlboiler", cmdArgs...)
		err = sqlBoilerCmd.Run()
		if err != nil {
			monitoring.Logger().Fatalf("%s\n", err)
		}

		//delete sqlboiler.toml
		err = os.Remove("sqlboiler.toml")
		if err != nil {
			monitoring.Logger().Fatalf("%s\n", err)
		}
	},
}

func init() {
	sqlboilerCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
