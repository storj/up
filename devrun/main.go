// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"
	"net"
	"storj.io/common/identity"
	"storj.io/common/pb"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc"
	"storj.io/common/socket"
	"storj.io/private/process"
	"storj.io/storj/satellite/console/consolewasm"
	"time"
)

var (

	// Commander CLI.
	rootCmd = &cobra.Command{
		Use:   "devrun",
		Short: "CLI to make it easier to create running dev clusters",
	}
	nodeIDCmd = &cobra.Command{
		Use:   "node-id",
		Short: "Generated node id string from identity file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeID, err := identity.NodeIDFromCertPath(args[0])
			if err != nil {
				return err
			}
			fmt.Println(nodeID)
			return nil
		},
	}

	waitForPort = &cobra.Command{
		Use:   "wait-for-port",
		Short: "Wait until ports is opened",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Trying tcp connection to " + args[0])
			for {
				timeout := time.Second
				conn, err := net.DialTimeout("tcp", args[0], timeout)
				if err != nil {
					time.Sleep(300 * time.Millisecond)
					continue
				}
				_ = conn.Close()
				return nil
			}
		},
	}

	waitForSatellite = &cobra.Command{
		Use:   "wait-for-satellite",
		Short: "Wait until satellite can be called and return with the full NodeURL",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			for {
				satellite, err := getSatelliteId(ctx, args[0])
				if err != nil {
					println("Couldn't connect to satellite. Retrying... " + err.Error())
					time.Sleep(1 * time.Second)
					continue
				}

				fmt.Println(satellite)
				return nil
			}
		},
	}

	credentials = &cobra.Command{
		Use:   "credentials",
		Short: "Generate test user with credentials",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)

			satelliteHost := args[0]
			sateliteNodeUrl, err := getSatelliteId(ctx, satelliteHost+":7777")
			if err != nil {
				return err
			}
			console := newConsoleEndpoints(satelliteHost+":10002", args[1])

			err = console.login(ctx)
			if err != nil {
				return err
			}
			projectID, err := console.getOrCreateProject(ctx)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("ProjectID: %s\n", projectID)
			apiKey, err := console.createAPIKey(ctx, projectID)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("API key: %s\n", apiKey)
			grant, err := consolewasm.GenAccessGrant(sateliteNodeUrl, apiKey, "Welcome1", projectID)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("Grant: %s\n", grant)

			return err
		},
	}
)

func getSatelliteId(ctx context.Context, address string) (string, error) {
	tlsOptions, err := getProcessTLSOptions(ctx)
	if err != nil {
		return "", err
	}

	dialer := rpc.NewDefaultDialer(tlsOptions)
	dialer.Pool = rpc.NewDefaultConnectionPool()

	dialer.DialTimeout = 10 * time.Second
	dialContext := socket.BackgroundDialer().DialContext
	dialer.Connector = rpc.NewDefaultTCPConnector(&rpc.ConnectorAdapter{DialContext: dialContext})

	conn, err := dialer.DialAddressInsecure(ctx, address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := pb.GetTimeRequest{}
	client := pb.NewDRPCNodeClient(conn)
	_, err = client.GetTime(ctx, &req)
	if err != nil {
		return "", err
	}
	for _, p := range conn.ConnectionState().PeerCertificates {
		if p.IsCA {
			id, err := identity.NodeIDFromCert(p)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s@%s", id, address), nil
		}
	}
	return "", fmt.Errorf("Couldn't find the right certiticate")
}

func getProcessTLSOptions(ctx context.Context) (*tlsopts.Options, error) {

	ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{
		Difficulty:  0,
		Concurrency: 1,
	})
	if err != nil {
		return nil, err
	}

	tlsConfig := tlsopts.Config{
		UsePeerCAWhitelist: false,
		PeerIDVersions:     "0",
	}

	tlsOptions, err := tlsopts.NewOptions(ident, tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	return tlsOptions, nil
}

func init() {
	rootCmd.AddCommand(nodeIDCmd)
	rootCmd.AddCommand(waitForPort)
	rootCmd.AddCommand(credentials)
	rootCmd.AddCommand(waitForSatellite)
	flag.Parse()
}

func main() {
	process.Exec(rootCmd)
}
