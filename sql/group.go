// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Author: Peter Mattis (peter@cockroachlabs.com)

package sql

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach/roachpb"
	"github.com/cockroachdb/cockroach/sql/parser"
	"github.com/cockroachdb/cockroach/util/log"
)

var aggregates = map[string]func() aggregateImpl{
	"avg":   newAvgAggregate,
	"count": newCountAggregate,
	"max":   newMaxAggregate,
	"min":   newMinAggregate,
	"sum":   newSumAggregate,
}

func (p *planner) groupBy(n *parser.Select, s *scanNode) (*groupNode, *roachpb.Error) {
	// Determine if aggregation is being performed. This check is done on the raw
	// Select expressions as simplification might have removed aggregation
	// functions (e.g. `SELECT MIN(1)` -> `SELECT 1`).
	if isAggregate := p.isAggregate(n); !isAggregate {
		return nil, nil
	}

	// Start by normalizing the GROUP BY expressions (to match what has been done to
	// the SELECT expressions in addRender) so that we can compare to them later.
	// This is done before determining if aggregation is being performed, because
	// that determination is made during validation, which will require matching
	// expressions.
	for i := range n.GroupBy {
		resolved, pErr := s.resolveQNames(n.GroupBy[i])
		if pErr != nil {
			return nil, pErr
		}

		// We could potentially skip this, since it will be checked in addRender,
		// but checking now allows early err return.
		if _, err := resolved.TypeCheck(p.evalCtx.Args); err != nil {
			return nil, roachpb.NewError(err)
		}

		norm, err := p.parser.NormalizeExpr(p.evalCtx, resolved)
		if err != nil {
			return nil, roachpb.NewError(err)
		}

		// If a col index is specified, replace it with that expression first.
		// NB: This is not a deep copy, and thus when extractAggregateFuncs runs
		// on s.render, the GroupBy expressions can contain wrapped qvalues.
		// aggregateFunc's Eval() method handles being called during grouping.
		if col, err := s.colIndex(norm); err != nil {
			return nil, roachpb.NewError(err)
		} else if col >= 0 {
			n.GroupBy[i] = s.render[col]
		} else {
			n.GroupBy[i] = norm
		}

	}

	// Normalize and check the HAVING expression too if it exists.
	if n.Having != nil {
		having, pErr := s.resolveQNames(n.Having.Expr)
		if pErr != nil {
			return nil, pErr
		}

		having, err := p.parser.NormalizeExpr(p.evalCtx, having)
		if err != nil {
			return nil, roachpb.NewError(err)
		}

		havingType, err := having.TypeCheck(p.evalCtx.Args)
		if err != nil {
			return nil, roachpb.NewError(err)
		}
		if !(havingType == parser.DummyBool || havingType == parser.DNull) {
			return nil, roachpb.NewUErrorf("argument of HAVING must be type %s, not type %s", parser.DummyBool.Type(), havingType.Type())
		}
		n.Having.Expr = having
	}

	group := &groupNode{
		planner: p,
		values:  valuesNode{columns: s.columns},
		render:  s.render,
	}

	visitor := extractAggregatesVisitor{
		n:           group,
		groupStrs:   make(map[string]struct{}, len(n.GroupBy)),
		groupedCopy: new(extractAggregatesVisitor),
	}

	for _, e := range n.GroupBy {
		visitor.groupStrs[e.String()] = struct{}{}
	}

	// A copy of the visitor that is used when a subtree appears in the GROUP BY.
	// One copy is allocated up-front, rather than potentially several on-the-fly,
	// to reduce allocations.
	*visitor.groupedCopy = visitor
	visitor.groupedCopy.groupedCopy = nil

	// Loop over the render expressions and extract any aggregate functions --
	// qvalues are also replaced (with identAggregates, which just return the last
	// value added to them for a bucket) to provide grouped-by values for each bucket.
	// After extraction, group.render will be entirely rendered from aggregateFuncs,
	// and group.funcs will contain all the functions which need to be fed values.
	for i := range group.render {
		expr, err := visitor.extract(group.render[i])
		if err != nil {
			return nil, roachpb.NewError(err)
		}
		group.render[i] = expr
	}

	if n.Having != nil {
		having, err := visitor.extract(n.Having.Expr)
		if err != nil {
			return nil, roachpb.NewError(err)
		}
		group.having = having
	}

	// Queries like `SELECT MAX(n) FROM t` expect a row of NULLs if nothing was aggregated.
	group.addNullBucketIfEmpty = len(n.GroupBy) == 0

	group.buckets = make(map[string]struct{})

	if log.V(2) {
		strs := make([]string, 0, len(group.funcs))
		for _, f := range group.funcs {
			strs = append(strs, f.String())
		}
		log.Infof("Group: %s", strings.Join(strs, ", "))
	}

	// Replace the render expressions in the scanNode with expressions that
	// compute only the arguments to the aggregate expressions.
	s.render = make([]parser.Expr, len(group.funcs))
	for i, f := range group.funcs {
		s.render[i] = f.arg
	}

	// Add the group-by expressions so they are available for bucketing.
	for _, g := range n.GroupBy {
		if err := s.addRender(parser.SelectExpr{Expr: g}); err != nil {
			return nil, err
		}
	}

	group.desiredOrdering = desiredAggregateOrdering(group.funcs)
	return group, nil
}

// A groupNode implements the planNode interface and handles the grouping logic.
// It "wraps" a planNode which is used to retrieve the ungrouped results.
type groupNode struct {
	planner *planner

	// The "wrapped" node (which returns ungrouped results).
	plan planNode

	render []parser.Expr
	having parser.Expr

	funcs []*aggregateFunc
	// The set of bucket keys.
	buckets map[string]struct{}

	addNullBucketIfEmpty bool

	values    valuesNode
	populated bool

	// During rendering, aggregateFuncs compute their result for group.currentBucket.
	currentBucket string

	desiredOrdering []int
	pErr            *roachpb.Error
}

func (n *groupNode) Columns() []column {
	return n.values.Columns()
}

func (n *groupNode) Ordering() ([]int, int) {
	// TODO(dt): aggregate buckets are returned un-ordered for now.
	return nil, 0
}

func (n *groupNode) Values() parser.DTuple {
	return n.values.Values()
}

func (n *groupNode) Next() bool {
	if !n.populated && n.pErr == nil {
		n.computeAggregates()
	}
	if n.pErr != nil {
		return false
	}
	return n.values.Next()
}

func (n *groupNode) computeAggregates() {
	var scratch []byte

	// Loop over the rows passing the values into the corresponding aggregation
	// functions.
	for n.plan.Next() {
		values := n.plan.Values()
		aggregatedValues, groupedValues := values[:len(n.funcs)], values[len(n.funcs):]

		// TODO(dt): optimization: skip buckets when underlying plan is ordered by grouped values.

		var encoded []byte
		var err error
		encoded, err = encodeDTuple(scratch, groupedValues)
		n.pErr = roachpb.NewError(err)
		if n.pErr != nil {
			return
		}

		n.buckets[string(encoded)] = struct{}{}

		// Feed the aggregateFuncs for this bucket the non-grouped values.
		for i, value := range aggregatedValues {
			if n.pErr = n.funcs[i].add(encoded, value); n.pErr != nil {
				return
			}
		}
		scratch = encoded[:0]
	}

	n.pErr = n.plan.PErr()
	if n.pErr != nil {
		return
	}

	if len(n.buckets) < 1 && n.addNullBucketIfEmpty {
		n.buckets[""] = struct{}{}
	}

	// Since this controls Eval behavior of aggregateFunc, it is not set until init is complete.
	n.populated = true

	// Render the results.
	n.values.rows = make([]parser.DTuple, 0, len(n.buckets))
	for k := range n.buckets {
		n.currentBucket = k

		if n.having != nil {
			res, err := n.having.Eval(n.planner.evalCtx)
			if err != nil {
				n.pErr = roachpb.NewError(err)
				return
			}
			if res, err := parser.GetBool(res); err != nil {
				n.pErr = roachpb.NewError(err)
				return
			} else if !res {
				continue
			}
		}

		row := make(parser.DTuple, 0, len(n.render))
		for _, r := range n.render {
			res, err := r.Eval(n.planner.evalCtx)
			if err != nil {
				n.pErr = roachpb.NewError(err)
				return
			}
			row = append(row, res)
		}

		n.values.rows = append(n.values.rows, row)
	}

}

func (n *groupNode) PErr() *roachpb.Error {
	return n.pErr
}

func (n *groupNode) ExplainPlan() (name, description string, children []planNode) {
	name = "group"
	strs := make([]string, 0, len(n.funcs))
	for _, f := range n.funcs {
		strs = append(strs, f.String())
	}
	description = strings.Join(strs, ", ")
	return name, description, []planNode{n.plan}
}

// wrap the supplied planNode with the groupNode if grouping/aggregation is required.
func (n *groupNode) wrap(plan planNode) planNode {
	if n == nil {
		return plan
	}
	n.plan = plan
	return n
}

// isNotNullFilter adds as a "col IS NOT NULL" constraint to the expression if
// the groupNode has a desired ordering on col (see
// desiredAggregateOrdering). A desired ordering will only be present if there
// is a single MIN/MAX aggregation function.
func (n *groupNode) isNotNullFilter(expr parser.Expr) parser.Expr {
	if len(n.desiredOrdering) != 1 {
		return expr
	}
	i := n.desiredOrdering[0]
	if i < 0 {
		i = -i
	}
	f := n.funcs[i-1]
	isNotNull := &parser.ComparisonExpr{
		Operator: parser.IsNot,
		Left:     f.arg,
		Right:    parser.DNull,
	}
	if expr == nil {
		return isNotNull
	}
	return &parser.AndExpr{
		Left:  expr,
		Right: isNotNull,
	}
}

// desiredAggregateOrdering computes the desired output ordering from the
// scan. It looks for an output column index containing a simple MIN/MAX
// aggregation. If zero or multiple MIN/MAX aggregations are requested then no
// ordering will be requested. A negative index indicates a MAX aggregation was
// requested for the output column.
func desiredAggregateOrdering(funcs []*aggregateFunc) []int {
	var limit int
	for i, f := range funcs {
		impl := f.create()
		switch impl.(type) {
		case *maxAggregate, *minAggregate:
			if limit != 0 || f.arg == nil {
				return nil
			}
			switch f.arg.(type) {
			case *qvalue:
				limit = i + 1
				if _, ok := impl.(*maxAggregate); ok {
					limit = -limit
				}
			default:
				return nil
			}

		default:
			return nil
		}
	}
	if limit == 0 {
		return nil
	}
	return []int{limit}
}

type extractAggregatesVisitor struct {
	n         *groupNode
	groupStrs map[string]struct{}

	// groupedCopy is nil when visitor is in an Expr subtree that appears in the GROUP BY clause.
	groupedCopy *extractAggregatesVisitor
	err         error
}

var _ parser.Visitor = &extractAggregatesVisitor{}

func (v *extractAggregatesVisitor) Visit(expr parser.Expr, pre bool) (parser.Visitor, parser.Expr) {
	if !pre || v.err != nil {
		return nil, expr
	}

	// This expression is in the GROUP BY - switch to the copy that will accept
	// qvalues for this and any subtrees.
	if _, ok := v.groupStrs[expr.String()]; ok && v.groupedCopy != nil {
		v = v.groupedCopy
	}

	switch t := expr.(type) {
	case *parser.FuncExpr:
		if len(t.Name.Indirect) > 0 {
			break
		}
		if impl, ok := aggregates[strings.ToLower(string(t.Name.Base))]; ok {
			if len(t.Exprs) != 1 {
				// Type checking has already run on these expressions thus
				// if an aggregate function of the wrong arity gets here,
				// something has gone really wrong.
				panic(fmt.Sprintf("%s has %d arguments (expected 1)", t.Name.Base, len(t.Exprs)))
			}

			f := &aggregateFunc{
				expr:    t,
				arg:     t.Exprs[0],
				create:  impl,
				group:   v.n,
				buckets: make(map[string]aggregateImpl),
			}
			if t.Type == parser.Distinct {
				f.seen = make(map[string]struct{})
			}
			v.n.funcs = append(v.n.funcs, f)
			return nil, f
		}
	case *qvalue:
		if v.groupedCopy != nil {
			v.err = fmt.Errorf("column \"%s\" must appear in the GROUP BY clause or be used in an aggregate function", t.col.Name)
			return v, expr
		}
		f := &aggregateFunc{
			expr:    t,
			arg:     t,
			create:  newIdentAggregate,
			group:   v.n,
			buckets: make(map[string]aggregateImpl),
		}
		v.n.funcs = append(v.n.funcs, f)
		return nil, f
	}
	return v, expr
}

// Extract aggregateFuncs from exprs that use aggregation and check if they are valid.
// An expression is valid if:
// - it is an aggregate expression, or
// - it appears verbatim in groupBy, or
// - it is not a qvalue, and all of its subexpressions (as defined by
// its Walk implementation) are valid
// NB: "verbatim" above is defined using a string-equality comparison
// as an approximation of a recursive tree-equality comparison.
//
// For example:
// Invalid: `SELECT k, SUM(v) FROM kv`
// - `k` is unaggregated and does not appear in the (missing) GROUP BY.
// Valid:      `SELECT k, SUM(v) FROM kv GROUP BY k`
// Also valid: `SELECT UPPER(k), SUM(v) FROM kv GROUP BY UPPER(k)`
// - `UPPER(k)` appears in GROUP BY.
// Also valid: `SELECT UPPER(k), SUM(v) FROM kv GROUP BY k`
// - `k` appears in GROUP BY, so `UPPER(k)` is OK, but...
// Invalid:    `SELECT k, SUM(v) FROM kv GROUP BY UPPER(k)`
// - `k` does not appear in GROUP BY; UPPER(k) does nothing to help here.
func (v extractAggregatesVisitor) extract(expr parser.Expr) (parser.Expr, error) {
	expr = parser.WalkExpr(&v, expr)
	return expr, v.err
}

var _ parser.Visitor = &isAggregateVisitor{}

type isAggregateVisitor struct {
	aggregated bool
}

func (v *isAggregateVisitor) Visit(expr parser.Expr, pre bool) (parser.Visitor, parser.Expr) {
	if !pre {
		return nil, expr
	}

	if t, ok := expr.(*parser.FuncExpr); ok {
		if _, ok := aggregates[strings.ToLower(string(t.Name.Base))]; ok {
			v.aggregated = true
			return nil, expr
		}
	}

	return v, expr
}

func (p *planner) isAggregate(n *parser.Select) bool {
	if n.Having != nil || len(n.GroupBy) > 0 {
		return true
	}

	for _, target := range n.Exprs {
		_ = parser.WalkExpr(&p.isAggregateVisitor, target.Expr)
		if p.isAggregateVisitor.aggregated {
			return true
		}
	}
	return false
}

var _ parser.VariableExpr = &aggregateFunc{}

type aggregateFunc struct {
	expr    parser.Expr
	arg     parser.Expr
	create  func() aggregateImpl
	group   *groupNode
	buckets map[string]aggregateImpl
	seen    map[string]struct{}
}

func (a *aggregateFunc) add(bucket []byte, d parser.Datum) *roachpb.Error {
	// NB: the compiler *should* optimize `myMap[string(myBytes)]`. See:
	// https://github.com/golang/go/commit/f5f5a8b6209f84961687d993b93ea0d397f5d5bf

	if a.seen != nil {
		encoded, pErr := encodeDatum(bucket, d)
		if pErr != nil {
			return pErr
		}
		if _, ok := a.seen[string(encoded)]; ok {
			// skip
			return nil
		}
		a.seen[string(encoded)] = struct{}{}
	}

	impl, ok := a.buckets[string(bucket)]
	if !ok {
		impl = a.create()
		a.buckets[string(bucket)] = impl
	}

	return roachpb.NewError(impl.add(d))
}

func (*aggregateFunc) Variable() {}

func (a *aggregateFunc) String() string {
	return a.expr.String()
}

func (a *aggregateFunc) Walk(v parser.Visitor) {
}

func (a *aggregateFunc) TypeCheck(args parser.MapArgs) (parser.Datum, error) {
	return a.expr.TypeCheck(args)
}

func (a *aggregateFunc) Eval(ctx parser.EvalContext) (parser.Datum, error) {
	// During init of the group buckets, grouped expressions (i.e. wrapped
	// qvalues) are Eval()'ed to determine the bucket for a row, so pass these
	// calls through to the underlying `arg` expr Eval until init is done.
	if !a.group.populated {
		return a.arg.Eval(ctx)
	}

	found, ok := a.buckets[a.group.currentBucket]
	if !ok {
		found = a.create()
	}

	datum, err := found.result()
	if err != nil {
		return nil, err
	}

	// This is almost certainly the identity. Oh well.
	return datum.Eval(ctx)
}

func encodeDatum(b []byte, d parser.Datum) ([]byte, *roachpb.Error) {
	if values, ok := d.(parser.DTuple); ok {
		dt, err := encodeDTuple(b, values)
		return dt, roachpb.NewError(err)
	}
	return encodeTableKey(b, d)
}

func encodeDTuple(b []byte, d parser.DTuple) ([]byte, error) {
	for _, val := range d {
		var pErr *roachpb.Error
		b, pErr = encodeDatum(b, val)
		if pErr != nil {
			return nil, pErr.GoError()
		}
	}
	return b, nil
}

type aggregateImpl interface {
	add(parser.Datum) error
	result() (parser.Datum, error)
}

var _ aggregateImpl = &avgAggregate{}
var _ aggregateImpl = &countAggregate{}
var _ aggregateImpl = &maxAggregate{}
var _ aggregateImpl = &minAggregate{}
var _ aggregateImpl = &sumAggregate{}
var _ aggregateImpl = &identAggregate{}

// In order to render the unaggregated (i.e. grouped) fields, during aggregation,
// the values for those fields have to be stored for each bucket.
// The `identAggregate` provides an "aggregate" function that actually
// just returns the last value passed to `add`, unchanged. For accumulating
// and rendering though it behaves like the other aggregate functions,
// allowing both those steps to avoid special-casing grouped vs aggregated fields.
type identAggregate struct {
	val parser.Datum
}

func newIdentAggregate() aggregateImpl {
	return &identAggregate{}
}

func (a *identAggregate) add(datum parser.Datum) error {
	a.val = datum
	return nil
}

func (a *identAggregate) result() (parser.Datum, error) {
	return a.val, nil
}

type avgAggregate struct {
	sumAggregate
	count int
}

func newAvgAggregate() aggregateImpl {
	return &avgAggregate{}
}

func (a *avgAggregate) add(datum parser.Datum) error {
	if datum == parser.DNull {
		return nil
	}
	if err := a.sumAggregate.add(datum); err != nil {
		return err
	}
	a.count++
	return nil
}

func (a *avgAggregate) result() (parser.Datum, error) {
	sum, err := a.sumAggregate.result()
	if err != nil {
		return parser.DNull, err
	}
	if sum == parser.DNull {
		return sum, nil
	}
	switch t := sum.(type) {
	case parser.DInt:
		return parser.DFloat(t) / parser.DFloat(a.count), nil
	case parser.DFloat:
		return t / parser.DFloat(a.count), nil
	default:
		return parser.DNull, fmt.Errorf("unexpected SUM result type: %s", t.Type())
	}
}

type countAggregate struct {
	count int
}

func newCountAggregate() aggregateImpl {
	return &countAggregate{}
}

func (a *countAggregate) add(datum parser.Datum) error {
	if datum == parser.DNull {
		return nil
	}
	switch t := datum.(type) {
	case parser.DTuple:
		for _, d := range t {
			if d != parser.DNull {
				a.count++
				break
			}
		}
	default:
		a.count++
	}
	return nil
}

func (a *countAggregate) result() (parser.Datum, error) {
	return parser.DInt(a.count), nil
}

type maxAggregate struct {
	max parser.Datum
}

func newMaxAggregate() aggregateImpl {
	return &maxAggregate{}
}

func (a *maxAggregate) add(datum parser.Datum) error {
	if datum == parser.DNull {
		return nil
	}
	if a.max == nil {
		a.max = datum
		return nil
	}
	c := a.max.Compare(datum)
	if c < 0 {
		a.max = datum
	}
	return nil
}

func (a *maxAggregate) result() (parser.Datum, error) {
	if a.max == nil {
		return parser.DNull, nil
	}
	return a.max, nil
}

type minAggregate struct {
	min parser.Datum
}

func newMinAggregate() aggregateImpl {
	return &minAggregate{}
}

func (a *minAggregate) add(datum parser.Datum) error {
	if datum == parser.DNull {
		return nil
	}
	if a.min == nil {
		a.min = datum
		return nil
	}
	c := a.min.Compare(datum)
	if c > 0 {
		a.min = datum
	}
	return nil
}

func (a *minAggregate) result() (parser.Datum, error) {
	if a.min == nil {
		return parser.DNull, nil
	}
	return a.min, nil
}

type sumAggregate struct {
	sum parser.Datum
}

func newSumAggregate() aggregateImpl {
	return &sumAggregate{}
}

func (a *sumAggregate) add(datum parser.Datum) error {
	if datum == parser.DNull {
		return nil
	}
	if a.sum == nil {
		a.sum = datum
		return nil
	}

	switch t := datum.(type) {
	case parser.DInt:
		if v, ok := a.sum.(parser.DInt); ok {
			a.sum = v + t
			return nil
		}

	case parser.DFloat:
		if v, ok := a.sum.(parser.DFloat); ok {
			a.sum = v + t
			return nil
		}
	}

	return fmt.Errorf("unexpected SUM argument type: %s", datum.Type())
}

func (a *sumAggregate) result() (parser.Datum, error) {
	if a.sum == nil {
		return parser.DNull, nil
	}
	return a.sum, nil
}
