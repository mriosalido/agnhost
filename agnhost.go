/*
Copyright 2019 The Kubernetes Authors.

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

package main

import (
	"os"

	"github.com/spf13/cobra"

	"k8s.io/component-base/cli"
	auditproxy "agnhost/audit-proxy"
	"agnhost/connect"
	crdconvwebhook "agnhost/crd-conversion-webhook"
	"agnhost/dns"
	"agnhost/entrypoint-tester"
	"agnhost/fakegitserver"
	grpchealthchecking "agnhost/grpc-health-checking"
	"agnhost/guestbook"
	"agnhost/inclusterclient"
	"agnhost/liveness"
	logsgen "agnhost/logs-generator"
	"agnhost/mounttest"
	"agnhost/net"
	"agnhost/netexec"
	"agnhost/nettest"
	nosnat "agnhost/no-snat-test"
	nosnatproxy "agnhost/no-snat-test-proxy"
	"agnhost/openidmetadata"
	"agnhost/pause"
	portforwardtester "agnhost/port-forward-tester"
	"agnhost/porter"
	resconsumerctrl "agnhost/resource-consumer-controller"
	servehostname "agnhost/serve-hostname"
	testwebserver "agnhost/test-webserver"
	"agnhost/webhook"
	"agnhost/echoserver"
	resconsumer "agnhost/resource-consumer"
	"agnhost/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "app",
		Version: version.Version,
	}

	rootCmd.AddCommand(auditproxy.CmdAuditProxy)
	rootCmd.AddCommand(connect.CmdConnect)
	rootCmd.AddCommand(crdconvwebhook.CmdCrdConversionWebhook)
	rootCmd.AddCommand(dns.CmdDNSSuffix)
	rootCmd.AddCommand(dns.CmdDNSServerList)
	rootCmd.AddCommand(dns.CmdEtcHosts)
	rootCmd.AddCommand(entrypoint.CmdEntrypointTester)
	rootCmd.AddCommand(fakegitserver.CmdFakeGitServer)
	rootCmd.AddCommand(guestbook.CmdGuestbook)
	rootCmd.AddCommand(inclusterclient.CmdInClusterClient)
	rootCmd.AddCommand(liveness.CmdLiveness)
	rootCmd.AddCommand(logsgen.CmdLogsGenerator)
	rootCmd.AddCommand(mounttest.CmdMounttest)
	rootCmd.AddCommand(net.CmdNet)
	rootCmd.AddCommand(netexec.CmdNetexec)
	rootCmd.AddCommand(nettest.CmdNettest)
	rootCmd.AddCommand(nosnat.CmdNoSnatTest)
	rootCmd.AddCommand(nosnatproxy.CmdNoSnatTestProxy)
	rootCmd.AddCommand(pause.CmdPause)
	rootCmd.AddCommand(porter.CmdPorter)
	rootCmd.AddCommand(portforwardtester.CmdPortForwardTester)
	rootCmd.AddCommand(resconsumerctrl.CmdResourceConsumerController)
	rootCmd.AddCommand(servehostname.CmdServeHostname)
	rootCmd.AddCommand(testwebserver.CmdTestWebserver)
	rootCmd.AddCommand(webhook.CmdWebhook)
	rootCmd.AddCommand(openidmetadata.CmdTestServiceAccountIssuerDiscovery)
	rootCmd.AddCommand(grpchealthchecking.CmdGrpcHealthChecking)
	rootCmd.AddCommand(echoserver.CmdEchoServer)
	rootCmd.AddCommand(resconsumer.CmdResourceConsumer)

	// NOTE(claudiub): Some tests are passing logging related flags, so we need to be able to
	// accept them. This will also include them in the printed help.
	code := cli.Run(rootCmd)
	os.Exit(code)
}
