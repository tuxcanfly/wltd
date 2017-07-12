// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package walletd

import (
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
)

type WalletDaemon struct {
	chainParams *chaincfg.Params
	wg          sync.WaitGroup

	started bool
	quit    chan struct{}
	quitMu  sync.Mutex
}

func NewWalletDaemon(activeNet *chaincfg.Params) *WalletDaemon {
	return &WalletDaemon{
		chainParams: activeNet,
	}
}

func (w *WalletDaemon) Start() {
	w.quitMu.Lock()
	select {
	case <-w.quit:
		// Restart the walletd goroutines after shutdown finishes.
		w.WaitForShutdown()
		w.quit = make(chan struct{})
	default:
		// Ignore when the walletd is still running.
		if w.started {
			w.quitMu.Unlock()
			return
		}
		w.started = true
	}
	w.quitMu.Unlock()
}

// quitChan atomically reads the quit channel.
func (w *WalletDaemon) quitChan() <-chan struct{} {
	w.quitMu.Lock()
	c := w.quit
	w.quitMu.Unlock()
	return c
}

// Stop signals all wallet goroutines to shutdown.
func (w *WalletDaemon) Stop() {
	w.quitMu.Lock()
	quit := w.quit
	w.quitMu.Unlock()

	select {
	case <-quit:
	default:
		close(quit)
	}
}

// ShuttingDown returns whether the wallet is currently in the process of
// shutting down or not.
func (w *WalletDaemon) ShuttingDown() bool {
	select {
	case <-w.quitChan():
		return true
	default:
		return false
	}
}

// WaitForShutdown blocks until all wallet goroutines have finished executing.
func (w *WalletDaemon) WaitForShutdown() {
	w.wg.Wait()
}

// ChainParams returns the network parameters for the blockchain the wallet
// belongs to.
func (w *WalletDaemon) ChainParams() *chaincfg.Params {
	return w.chainParams
}
