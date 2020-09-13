/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"github.com/hermeschat/engine/monitoring"
	"github.com/hermeschat/engine/transport/grpc"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves hermes",
	Long:  `serve starts hermes`,
	Run: func(cmd *cobra.Command, args []string) {

		monitoring.Logger().Info("Loading Config")
		monitoring.Logger().Info("Initiating DB package")
		grpc.CreateGRPCServer(context.Background())
		monitoring.Logger().Info("Initializing Hermes")

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Aliases = []string{"serv"}
}
