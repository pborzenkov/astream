package report

import (
	"sort"
	"time"
)

type Operations map[string]float64

type Operation struct {
	Name   string
	Amount float64
}

func (o Operations) Add(op string, val float64) {
	o[op] += val
}

func (o Operations) Sorted() []*Operation {
	ops := make([]*Operation, 0, len(o))
	for o, v := range o {
		ops = append(ops, &Operation{
			Name:   o,
			Amount: v,
		})
	}

	sort.Slice(ops, func(i, j int) bool {
		return ops[i].Name < ops[j].Name
	})

	return ops
}

type AggregationByOp struct {
	Date time.Time
	Type GroupType
	Ops  Operations
}

type AggregateByOperationOption func(*aggregateByOperationConfig)

type aggregateByOperationConfig struct {
	separateFine bool
}

func WithSeparateFine(v bool) AggregateByOperationOption {
	return func(cfg *aggregateByOperationConfig) {
		cfg.separateFine = v
	}
}

func (r *Report) AggregateByOperation(g GroupType, opts ...AggregateByOperationOption) []*AggregationByOp {
	var c time.Time
	var a []*AggregationByOp

	var cfg aggregateByOperationConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	for _, t := range r.transactions {
		if !g.Same(c, t.date) {
			a = append(a, &AggregationByOp{
				Date: t.date,
				Type: g,
				Ops:  make(map[string]float64),
			})
			c = t.date
		}
		if !cfg.separateFine && t.ttype == "пени" {
			t.ttype = "проценты"
		}

		a[len(a)-1].Ops.Add(t.ttype, t.amount)
	}

	return a
}
