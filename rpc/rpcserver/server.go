// Copyright (c) 2015-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package rpcserver implements the RPC API and is used by the main package to
// start gRPC services.
//
// Full documentation of the API implemented by this package is maintained in a
// language-agnostic document:
//
//   https://github.com/btcsuite/btcwallet/blob/master/rpc/documentation/api.md
//
// Any API changes must be performed according to the steps listed here:
//
//   https://github.com/btcsuite/btcwallet/blob/master/rpc/documentation/serverchanges.md
package rpcserver

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/btcsuite/btcd/chaincfg"
	pb "github.com/tuxcanfly/wltd/rpc/walletdrpc"
	"github.com/tuxcanfly/wltd/walletd"
)

// Public API version constants
const (
	semverString = "2.0.1"
	semverMajor  = 2
	semverMinor  = 0
	semverPatch  = 1
)

// versionServer provides RPC clients with the ability to query the RPC server
// version.
type versionServer struct {
}

// walletDaemonServer provides wallet services for RPC clients.
type walletDaemonServer struct {
	walletd   *walletd.WalletDaemon
	activeNet *chaincfg.Params
}

// StartVersionService creates an implementation of the VersionService and
// registers it with the gRPC server.
func StartVersionService(server *grpc.Server) {
	pb.RegisterVersionServiceServer(server, &versionServer{})
}

func (*versionServer) Version(ctx context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		VersionString: semverString,
		Major:         semverMajor,
		Minor:         semverMinor,
		Patch:         semverPatch,
	}, nil
}

// StartWalletDaemonService creates an implementation of the WalletDaemonService and
// registers it with the gRPC server.
func StartWalletDaemonService(server *grpc.Server, walletd *walletd.WalletDaemon,
	activeNet *chaincfg.Params) {
	service := &walletDaemonServer{walletd, activeNet}
	pb.RegisterWalletDaemonServiceServer(server, service)
}

func (s *walletDaemonServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}

func (s *walletDaemonServer) Network(ctx context.Context, req *pb.NetworkRequest) (
	*pb.NetworkResponse, error) {

	return &pb.NetworkResponse{ActiveNetwork: uint32(s.walletd.ChainParams().Net)}, nil
}
