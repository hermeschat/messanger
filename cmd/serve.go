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

	"github.com/amirrezaask/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hermes/pkg/db"
	"hermes/pkg/grpcserver"
	"hermes/pkg/subscription"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves hermes",
	Long:  `serve starts hermes`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			subscription.Clean()
		}()
		logrus.Info("Loading Config")
		config.Init()
		logrus.Info("Initiating DB package")
		db.Init()
		grpcserver.CreateGRPCServer(context.Background())
		logrus.Info("Initializing Hermes")

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Aliases = []string{"serv"}
}
