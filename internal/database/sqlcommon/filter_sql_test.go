// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlcommon

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
)

func TestSQLQueryFactory(t *testing.T) {
	s, _ := getMockDB()
	fb := database.MessageQueryFactory.NewFilter(context.Background())
	f := fb.And(
		fb.Eq("namespace", "ns1"),
		fb.Or(
			fb.Eq("id", "35c11cba-adff-4a4d-970a-02e3a0858dc8"),
			fb.Eq("id", "caefb9d1-9fc9-4d6a-a155-514d3139adf7"),
		),
		fb.Gt("sequence", "12345"),
		fb.Eq("confirmed", nil),
	).
		Skip(50).
		Limit(25).
		Sort("namespace").
		Descending()

	sel := squirrel.Select("*").From("mytable")
	sel, err := s.filterSelect(context.Background(), "", sel, f, map[string]string{
		"namespace": "ns",
	})
	assert.NoError(t, err)

	sqlFilter, args, err := sel.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM mytable WHERE (ns = ? AND (id = ? OR id = ?) AND seq > ? AND confirmed IS NULL) ORDER BY ns DESC LIMIT 25 OFFSET 50", sqlFilter)
	assert.Equal(t, "ns1", args[0])
	assert.Equal(t, "35c11cba-adff-4a4d-970a-02e3a0858dc8", args[1])
	assert.Equal(t, "caefb9d1-9fc9-4d6a-a155-514d3139adf7", args[2])
	assert.Equal(t, int64(12345), args[3])
}

func TestSQLQueryFactoryExtraOps(t *testing.T) {

	s, _ := getMockDB()
	fb := database.MessageQueryFactory.NewFilter(context.Background())
	u := fftypes.MustParseUUID("4066ABDC-8BBD-4472-9D29-1A55B467F9B9")
	f := fb.And(
		fb.In("created", []driver.Value{1, 2, 3}),
		fb.NotIn("created", []driver.Value{1, 2, 3}),
		fb.Eq("id", u),
		fb.In("id", []driver.Value{*u}),
		fb.Neq("id", nil),
		fb.Lt("created", "0"),
		fb.Lte("created", "0"),
		fb.Gte("created", "0"),
		fb.Neq("created", "0"),
		fb.Gt("sequence", 12345),
		fb.Contains("context", "abc"),
		fb.NotContains("context", "def"),
		fb.IContains("context", "ghi"),
		fb.NotIContains("context", "jkl"),
	)

	sel := squirrel.Select("*").From("mytable AS mt")
	sel, err := s.filterSelect(context.Background(), "mt", sel, f, nil)
	assert.NoError(t, err)

	sqlFilter, _, err := sel.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM mytable AS mt WHERE (mt.created IN (?,?,?) AND mt.created NOT IN (?,?,?) AND mt.id = ? AND mt.id IN (?) AND mt.id IS NOT NULL AND mt.created < ? AND mt.created <= ? AND mt.created >= ? AND mt.created <> ? AND mt.seq > ? AND mt.context LIKE ? AND mt.context NOT LIKE ? AND mt.context ILIKE ? AND mt.context NOT ILIKE ?) ORDER BY mt.seq DESC", sqlFilter)
}

func TestSQLQueryFactoryFinalizeFail(t *testing.T) {
	s, _ := getMockDB()
	fb := database.MessageQueryFactory.NewFilter(context.Background())
	sel := squirrel.Select("*").From("mytable")
	_, err := s.filterSelect(context.Background(), "ns", sel, fb.Eq("namespace", map[bool]bool{true: false}), nil)
	assert.Regexp(t, "FF10149.*namespace", err.Error())
}

func TestSQLQueryFactoryBadOp(t *testing.T) {

	s, _ := getMockDB()
	sel := squirrel.Select("*").From("mytable")
	_, err := s.filterSelectFinalized(context.Background(), "", sel, &database.FilterInfo{
		Op: database.FilterOp("wrong"),
	}, nil)
	assert.Regexp(t, "FF10150.*wrong", err.Error())
}

func TestSQLQueryFactoryBadOpInOr(t *testing.T) {

	s, _ := getMockDB()
	sel := squirrel.Select("*").From("mytable")
	_, err := s.filterSelectFinalized(context.Background(), "", sel, &database.FilterInfo{
		Op: database.FilterOpOr,
		Children: []*database.FilterInfo{
			{Op: database.FilterOp("wrong")},
		},
	}, nil)
	assert.Regexp(t, "FF10150.*wrong", err.Error())
}

func TestSQLQueryFactoryBadOpInAnd(t *testing.T) {

	s, _ := getMockDB()
	sel := squirrel.Select("*").From("mytable")
	_, err := s.filterSelectFinalized(context.Background(), "", sel, &database.FilterInfo{
		Op: database.FilterOpAnd,
		Children: []*database.FilterInfo{
			{Op: database.FilterOp("wrong")},
		},
	}, nil)
	assert.Regexp(t, "FF10150.*wrong", err.Error())
}
