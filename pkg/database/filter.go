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

package database

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/kaleido-io/firefly/internal/i18n"
)

// Filter is the output of the builder
type Filter interface {
	// Sort adds a set of sort conditions (all in a single sort order)
	Sort(...string) Filter

	// Ascending sort order
	Ascending() Filter

	// Descending sort order
	Descending() Filter

	// Skip for pagination
	Skip(uint64) Filter

	// Limit for pagination
	Limit(uint64) Filter

	// Finalize completes the filter, and for the plugin to validated output structure to convert
	Finalize() (*FilterInfo, error)

	// Builder returns the builder that made it
	Builder() FilterBuilder
}

// MultiConditionFilter gives convenience methods to add conditions
type MultiConditionFilter interface {
	Filter
	// Add adds filters to the condition
	Condition(...Filter) MultiConditionFilter
}

type AndFilter interface{ MultiConditionFilter }

type OrFilter interface{ MultiConditionFilter }

// FilterOp enum of filter operations that must be implemented by plugins - the string value is
// used in the core string formatting method (for logging etc.)
type FilterOp string

const (
	FilterOpAnd      FilterOp = "&&"
	FilterOpOr       FilterOp = "||"
	FilterOpEq       FilterOp = "=="
	FilterOpNe       FilterOp = "!="
	FilterOpIn       FilterOp = "IN"
	FilterOpNotIn    FilterOp = "NI"
	FilterOpGt       FilterOp = ">"
	FilterOpLt       FilterOp = "<"
	FilterOpGte      FilterOp = ">="
	FilterOpLte      FilterOp = "<="
	FilterOpCont     FilterOp = "%="
	FilterOpNotCont  FilterOp = "%!"
	FilterOpICont    FilterOp = "^="
	FilterOpNotICont FilterOp = "^!"
)

// FilterBuilder is the syntax used to build the filter, where And() and Or() can be nested
type FilterBuilder interface {
	// Fields is the list of available fields
	Fields() []string
	// And requires all sub-filters to match
	And(and ...Filter) AndFilter
	// Or requires any of the sub-filters to match
	Or(and ...Filter) OrFilter
	// Eq equal
	Eq(name string, value driver.Value) Filter
	// Neq not equal
	Neq(name string, value driver.Value) Filter
	// In one of an array of values
	In(name string, value []driver.Value) Filter
	// NotIn not one of an array of values
	NotIn(name string, value []driver.Value) Filter
	// Lt less than
	Lt(name string, value driver.Value) Filter
	// Gt greater than
	Gt(name string, value driver.Value) Filter
	// Gte greater than or equal
	Gte(name string, value driver.Value) Filter
	// Lte less than or equal
	Lte(name string, value driver.Value) Filter
	// Contains allows the string anywhere - case sensitive
	Contains(name string, value driver.Value) Filter
	// NotContains disallows the string anywhere - case sensitive
	NotContains(name string, value driver.Value) Filter
	// IContains allows the string anywhere - case sensitive
	IContains(name string, value driver.Value) Filter
	// INotContains disallows the string anywhere - case sensitive
	NotIContains(name string, value driver.Value) Filter
}

// FilterInfo is the structure returned by Finalize to the plugin, to serialize this filter
// into the underlying database mechanism's filter language
type FilterInfo struct {
	Sort       []string
	Skip       uint64
	Limit      uint64
	Descending bool
	Field      string
	Op         FilterOp
	Values     []FieldSerialization
	Value      FieldSerialization
	Children   []*FilterInfo
}

func valueString(f FieldSerialization) string {
	v, _ := f.Value()
	switch tv := v.(type) {
	case nil:
		return "null"
	case int64:
		return strconv.FormatInt(tv, 10)
	default:
		return fmt.Sprintf("'%s'", tv)
	}
}

func (f *FilterInfo) filterString() string {
	switch f.Op {
	case FilterOpAnd, FilterOpOr:
		cs := make([]string, len(f.Children))
		for i, c := range f.Children {
			cs[i] = fmt.Sprintf("( %s )", c.filterString())
		}
		return strings.Join(cs, fmt.Sprintf(" %s ", f.Op))
	case FilterOpIn, FilterOpNotIn:
		strValues := make([]string, len(f.Values))
		for i, v := range f.Values {
			strValues[i] = valueString(v)
		}
		return fmt.Sprintf("%s %s [%s]", f.Field, f.Op, strings.Join(strValues, ","))
	default:
		return fmt.Sprintf("%s %s %s", f.Field, f.Op, valueString(f.Value))
	}
}

func (f *FilterInfo) String() string {

	var val strings.Builder

	val.WriteString(f.filterString())

	if len(f.Sort) > 0 {
		val.WriteString(fmt.Sprintf(" sort=%s", strings.Join(f.Sort, ",")))
		if f.Descending {
			val.WriteString(" descending")
		}
	}
	if f.Skip > 0 {
		val.WriteString(fmt.Sprintf(" skip=%d", f.Skip))
	}
	if f.Limit > 0 {
		val.WriteString(fmt.Sprintf(" limit=%d", f.Limit))
	}

	return val.String()
}

func (fb *filterBuilder) Fields() []string {
	keys := make([]string, len(fb.queryFields))
	i := 0
	for k := range fb.queryFields {
		keys[i] = k
		i++
	}
	return keys
}

type filterBuilder struct {
	ctx         context.Context
	queryFields queryFields
	sort        []string
	skip        uint64
	limit       uint64
	descending  bool
}

type baseFilter struct {
	fb       *filterBuilder
	children []Filter
	op       FilterOp
	field    string
	value    interface{}
}

func (f *baseFilter) Builder() FilterBuilder {
	return f.fb
}

func (f *baseFilter) Finalize() (fi *FilterInfo, err error) {
	var children []*FilterInfo
	var value FieldSerialization
	var values []FieldSerialization

	switch f.op {
	case FilterOpAnd, FilterOpOr:
		children = make([]*FilterInfo, len(f.children))
		for i, c := range f.children {
			if children[i], err = c.Finalize(); err != nil {
				return nil, err
			}
		}
	case FilterOpIn, FilterOpNotIn:
		fValues := f.value.([]driver.Value)
		values = make([]FieldSerialization, len(fValues))
		name := strings.ToLower(f.field)
		field, ok := f.fb.queryFields[name]
		if !ok {
			return nil, i18n.NewError(f.fb.ctx, i18n.MsgInvalidFilterField, name)
		}
		for i, fv := range fValues {
			values[i] = field.getSerialization()
			if err = values[i].Scan(fv); err != nil {
				return nil, i18n.WrapError(f.fb.ctx, err, i18n.MsgInvalidValueForFilterField, name)
			}
		}
	default:
		name := strings.ToLower(f.field)
		field, ok := f.fb.queryFields[name]
		if !ok {
			return nil, i18n.NewError(f.fb.ctx, i18n.MsgInvalidFilterField, name)
		}
		value = field.getSerialization()
		if err = value.Scan(f.value); err != nil {
			return nil, i18n.WrapError(f.fb.ctx, err, i18n.MsgInvalidValueForFilterField, name)
		}
	}
	return &FilterInfo{
		Children:   children,
		Op:         f.op,
		Field:      f.field,
		Values:     values,
		Value:      value,
		Sort:       f.fb.sort,
		Skip:       f.fb.skip,
		Limit:      f.fb.limit,
		Descending: f.fb.descending,
	}, nil
}

func (f *baseFilter) Sort(fields ...string) Filter {
	for _, field := range fields {
		if _, ok := f.fb.queryFields[field]; ok {
			f.fb.sort = append(f.fb.sort, field)
		}
	}
	return f
}

func (f *baseFilter) Skip(skip uint64) Filter {
	f.fb.skip = skip
	return f
}

func (f *baseFilter) Limit(limit uint64) Filter {
	f.fb.limit = limit
	return f
}

func (f *baseFilter) Ascending() Filter {
	f.fb.descending = false
	return f
}

func (f *baseFilter) Descending() Filter {
	f.fb.descending = true
	return f
}

type andFilter struct {
	baseFilter
}

func (fb *andFilter) Condition(children ...Filter) MultiConditionFilter {
	fb.children = append(fb.children, children...)
	return fb
}

func (fb *filterBuilder) And(and ...Filter) AndFilter {
	return &andFilter{
		baseFilter: baseFilter{
			fb:       fb,
			op:       FilterOpAnd,
			children: and,
		},
	}
}

type orFilter struct {
	baseFilter
}

func (fb *orFilter) Condition(children ...Filter) MultiConditionFilter {
	fb.children = append(fb.children, children...)
	return fb
}

func (fb *filterBuilder) Or(or ...Filter) OrFilter {
	return &orFilter{
		baseFilter: baseFilter{
			fb:       fb,
			op:       FilterOpOr,
			children: or,
		},
	}
}

func (fb *filterBuilder) Eq(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpEq, name, value)
}

func (fb *filterBuilder) Neq(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpNe, name, value)
}

func (fb *filterBuilder) In(name string, values []driver.Value) Filter {
	return fb.fieldFilter(FilterOpIn, name, values)
}

func (fb *filterBuilder) NotIn(name string, values []driver.Value) Filter {
	return fb.fieldFilter(FilterOpNotIn, name, values)
}

func (fb *filterBuilder) Lt(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpLt, name, value)
}

func (fb *filterBuilder) Gt(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpGt, name, value)
}

func (fb *filterBuilder) Gte(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpGte, name, value)
}

func (fb *filterBuilder) Lte(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpLte, name, value)
}

func (fb *filterBuilder) Contains(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpCont, name, value)
}

func (fb *filterBuilder) NotContains(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpNotCont, name, value)
}

func (fb *filterBuilder) IContains(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpICont, name, value)
}

func (fb *filterBuilder) NotIContains(name string, value driver.Value) Filter {
	return fb.fieldFilter(FilterOpNotICont, name, value)
}

func (fb *filterBuilder) fieldFilter(op FilterOp, name string, value interface{}) Filter {
	return &fieldFilter{
		baseFilter: baseFilter{
			fb:    fb,
			op:    op,
			field: name,
			value: value,
		},
	}
}

type fieldFilter struct {
	baseFilter
}
