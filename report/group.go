package report

import (
	"fmt"
	"time"
)

// GroupType controls aggregation group range
type GroupType int

// Possible Group types
const (
	Daily GroupType = 1 + iota
	Weekly
	Monthly
)

var groups = [...]string{
	"daily",
	"weekly",
	"monthly",
}

func (g GroupType) String() string {
	if Daily <= g && g <= Monthly {
		return groups[g-1]
	}
	return fmt.Sprintf("%!Group(%d)", g)
}

func (g *GroupType) Set(v string) error {
	for i := range groups {
		if groups[i] == v {
			*g = GroupType(i + 1)
			return nil
		}
	}

	return fmt.Errorf("unsupported group type %q", v)
}

func (g GroupType) Same(a, b time.Time) bool {
	am, ad, bm, bd := a.Month(), a.Day(), b.Month(), b.Day()
	ay, aw := a.ISOWeek()
	by, bw := b.ISOWeek()

	switch g {
	case Daily:
		return ay == by && am == bm && ad == bd
	case Weekly:
		return ay == by && aw == bw
	case Monthly:
		return ay == by && am == bm
	default:
		return false
	}
}
