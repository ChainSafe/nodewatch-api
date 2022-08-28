//go:generate go run github.com/99designs/gqlgen generate
// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

// Package graph contans graph related code
package graph

import (
	"eth2-crawler/store/peerstore"
	"eth2-crawler/store/record"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	peerStore    peerstore.Provider
	historyStore record.Provider
}

func NewResolver(peerStore peerstore.Provider, historyStore record.Provider) *Resolver {
	return &Resolver{peerStore: peerStore, historyStore: historyStore}
}
