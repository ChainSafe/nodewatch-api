// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package peerstore

import "errors"

var (
	ErrPeerNotFound = errors.New("unable to find the node")
)
