package main

import (
	"fmt"
	"os"
	"path/filepath"

	walletpb "github.com/btcsuite/btcwallet/rpc/walletrpc"
	flags "github.com/btcsuite/go-flags"
	wltdpb "github.com/tuxcanfly/wltd/rpc/walletdrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/btcsuite/btcutil"
)

type config struct {
	Create     string `short:"c" long:"create" description:"Create new Wallet"`
	GetBalance string `short:"b" long:"balance" description:"Get Wallet Balance"`
}

func createWallet(pass string) (string, error) {
	certificateFile := filepath.Join(btcutil.AppDataDir("wltd", false), "rpc.cert")
	creds, err := credentials.NewClientTLSFromFile(certificateFile, "localhost")
	if err != nil {
		return "", err
	}
	conn, err := grpc.Dial("localhost:18335", grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := wltdpb.NewWalletDaemonServiceClient(conn)
	req := &wltdpb.CreateWalletRequest{pass}
	resp, err := c.CreateWallet(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Uuid, nil
}

func getBalance(wid string) (int64, error) {
	certificateFile := filepath.Join(btcutil.AppDataDir("btcwallet", false), "rpc.cert")
	creds, err := credentials.NewClientTLSFromFile(certificateFile, "localhost")
	if err != nil {
		return 0, err
	}
	conn, err := grpc.Dial("localhost:18332", grpc.WithTransportCredentials(creds))
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	wltDir := filepath.Join(btcutil.AppDataDir("wltd", false), "wallets", wid, "testnet")
	c := walletpb.NewWalletLoaderServiceClient(conn)
	req := &walletpb.OpenWalletRequest{Path: wltDir}
	_, err = c.OpenWallet(context.Background(), req)
	if err != nil {
		return 0, err
	}

	w := walletpb.NewWalletServiceClient(conn)
	breq := &walletpb.BalanceRequest{}
	bresp, err := w.Balance(context.Background(), breq)
	if err != nil {
		return 0, err
	}

	creq := &walletpb.CloseWalletRequest{Path: wltDir}
	_, err = c.CloseWallet(context.Background(), creq)
	if err != nil {
		return 0, err
	}

	return bresp.Total, nil
}

func main() {
	cfg := config{}
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return
	}

	if cfg.Create == "" && cfg.GetBalance == "" {
		fmt.Fprintln(os.Stderr, "no cmd specified")
		os.Exit(1)
	}

	if cfg.Create != "" {
		wid, err := createWallet(cfg.Create)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Wallet id: %v\n", wid)
	}

	if cfg.GetBalance != "" {
		balance, err := getBalance(cfg.GetBalance)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Balance: %v\n", balance)
	}
}
