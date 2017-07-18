// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package kv_test

import (
	"github.com/juju/errors"
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb"
)

type testFaultInjectionSuite struct{}

var _ = Suite(testFaultInjectionSuite{})

func (s testFaultInjectionSuite) TestFaultInjectionBasic(c *C) {
	var cfg kv.InjectionConfig
	err := errors.New("foo")
	cfg.SetGetError(err)

	store, _, err := newStoreWithBootstrap()
	c.Assert(err, IsNil)
	storage := kv.NewInjectedStore(store, &cfg)
	txn, err := storage.Begin()
	c.Assert(err, IsNil)
	_, err = storage.BeginWithStartTS(0)
	c.Assert(err, IsNil)
	ver := kv.Version{1}
	snap, err := storage.GetSnapshot(ver)
	c.Assert(err, IsNil)
	b, err := txn.Get([]byte{'a'})
	c.Assert(err.Error(), Equals, errors.New("foo").Error())
	c.Assert(b, IsNil)
	b, err = snap.Get([]byte{'a'})
	c.Assert(err.Error(), Equals, errors.New("foo").Error())
	c.Assert(b, IsNil)
}

func newStoreWithBootstrap() (kv.Storage, *domain.Domain, error) {
	store, err := tidb.NewStore(tidb.EngineGoLevelDBMemory)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	do, err := tidb.BootstrapSession(store)
	return store, do, errors.Trace(err)
}
