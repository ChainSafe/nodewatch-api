// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package graph contans graph related code
package graph

import "eth2-crawler/store"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	peerStore store.Provider
}

func NewResolver(peerStore store.Provider) *Resolver {
	return &Resolver{peerStore: peerStore}
}
