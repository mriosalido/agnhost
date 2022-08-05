/*
Copyright 2015 The Kubernetes Authors.

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

package resconsumer

import (
	"fmt"
	"log"
	"net/http"
	"github.com/spf13/cobra"
)

var (
	port	int
)

var CmdResourceConsumer = &cobra.Command{
	Use:   "resource-consumer",
	Short: "Starts a HTTP server that spreads requests around resource consumers",
	Long:  "Starts a HTTP server that spreads requests around resource consumers. The HTTP server has the same endpoints and usage as the one spawned by the \"resource-consumer\" subcommand.",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

func main(cmd *cobra.Command, args []string) {
	resourceConsumerHandler := NewResourceConsumerHandler()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), resourceConsumerHandler))
}

func init() {
	CmdResourceConsumer.Flags().IntVar(&port, "port", 8080, "Port number.")
}