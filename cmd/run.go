/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"study/pkg/config"
	"study/pkg/k8s2eureka"
	"syscall"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server running")
		PrintFlags(cmd.Flags())
		stop := make(chan struct{})

		controller := k8s2eureka.CreateController()
		err := controller.InitController(*runConf)
		if err != nil {
			panic(err)
		}
		err = controller.Run(stop)
		if err != nil {
			panic(err)
		}
		WaitSignal(stop)
	},
}

func WaitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

var runConf *config.Config

func init() {
	RootCmd.AddCommand(runCmd)
	runConf = &config.Config{
		KubeConfigPath: "/root/.kube",
	}
	runCmd.Flags().StringVar(&conf.KubeConfigPath, "kubeconfig", "", "kubeconfig name")
}
