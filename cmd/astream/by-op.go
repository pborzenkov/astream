package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/alecthomas/kingpin"
	"github.com/pborzenkov/astream/report"
	"github.com/pkg/errors"
)

var byOpParams struct {
	path  string
	from  timeFlag
	to    timeFlag
	group report.GroupType
}

const (
	hsep = "================"
	sep  = "----------------"
)

func addByOpCommand(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command("by-op", "Generate investment report aggregated by operation type.")
	cmd.Arg(
		"report-file",
		"Report file to aggregate.",
	).Required().ExistingFileVar(&byOpParams.path)
	cmd.Flag(
		"from",
		"Start of the aggregation date range (closed).",
	).Short('f').SetValue(&byOpParams.from)
	cmd.Flag(
		"to",
		"End of the aggregation date range (open).",
	).Short('t').SetValue(&byOpParams.to)
	cmd.Flag(
		"group",
		"Group operations by the given range.",
	).Short('g').Default("daily").SetValue(&byOpParams.group)

	return cmd
}

// ByOp creates an aggregate of the given report by operation for the given
// time range and outputs it to stdout.
func byOp() error {
	r, err := reportFromFile(byOpParams.path)
	if err != nil {
		return err
	}

	if !byOpParams.from.time().IsZero() || !byOpParams.to.time().IsZero() {
		r = r.Slice(byOpParams.from.time(), byOpParams.to.time())
	}

	agg := r.AggregateByOperation(byOpParams.group)

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintf(tw, "Date\tOperation\tAmount\n%s\t%s\t%s\n", hsep, hsep, hsep)
	for _, a := range agg {
		for _, op := range a.Ops.Sorted() {
			fmt.Fprintf(tw, "%s\t%s\t%.2f\n", a.Date.Format("2006/01/02"), op.Name, op.Amount)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", sep, sep, sep)
	}

	return errors.Wrap(tw.Flush(), "tabwriter.Flush")
}
