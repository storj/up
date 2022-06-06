// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package up

import (
	"context"
	"fmt"
	"time"

	"storj.io/common/identity"
	"storj.io/common/pb"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc"
	"storj.io/common/socket"
)

// GetSatelliteID retrieves node identity from SSL endpoint.
// Only for testing. Using identified node id is not reliable.
func GetSatelliteID(ctx context.Context, address string) (string, error) {
	tlsOptions, err := getProcessTLSOptions(ctx)
	if err != nil {
		return "", err
	}

	dialer := rpc.NewDefaultDialer(tlsOptions)
	dialer.Pool = rpc.NewDefaultConnectionPool()

	dialer.DialTimeout = 10 * time.Second
	dialContext := socket.BackgroundDialer().DialContext

	//lint:ignore SA1019 it's safe to use TCP here instead of QUIC + TCP
	dialer.Connector = rpc.NewDefaultTCPConnector(&rpc.ConnectorAdapter{DialContext: dialContext}) //nolint:staticcheck

	conn, err := dialer.DialAddressInsecure(ctx, address)
	if err != nil {
		return "", err
	}
	defer func() { _ = conn.Close() }()

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
	return "", fmt.Errorf("couldn't find the right certiticate")
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
