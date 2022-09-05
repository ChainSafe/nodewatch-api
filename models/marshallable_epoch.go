// Copyright 2022 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package models

import (
	"fmt"
	"strconv"

	"github.com/protolambda/zrnt/eth2/beacon/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Epoch type is needed to support proper BSON marshaling and unmarshaling to/from MongoDB.
type Epoch common.Epoch

func (e *Epoch) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(toHexString(uint64(*e)))
}

func (e *Epoch) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	var container string
	rv := bson.RawValue{Type: t, Value: b}
	err := rv.Unmarshal(&container)
	if err != nil {
		return err
	}
	val, err := fromHexString(container)
	if err != nil {
		return err
	}
	*e = Epoch(val)
	return nil
}

func (e *Epoch) String() string {
	return common.Epoch(*e).String()
}

func toHexString(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

func fromHexString(s string) (uint64, error) {
	return strconv.ParseUint(s[2:], 16, 64)
}
