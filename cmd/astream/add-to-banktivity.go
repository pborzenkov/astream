package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/alecthomas/kingpin"
	"github.com/pborzenkov/astream/report"
	"github.com/pkg/errors"
)

const qifTemplate = `
!Account
NAlfa•Bank Stream RUB
TInvst
^
!Type:Invst
{{- range $date := .}}
{{- range $op, $value := .Ops}}
D{{$date.Date.Format "02/01/2006"}}
{{printOp $op $value -}}
^
{{- end}}
{{- end}}
`

var qif = template.Must(template.New("qif").
	Funcs(template.FuncMap{"printOp": func(op string, val float64) string {
		switch op {
		case "инвестирование":
			return fmt.Sprintf("NBuy\nYПоток\nI1.00\nQ%0.2f\nT%0.2f\nO0.00\n", val, val)
		case "выплата ОД":
			return fmt.Sprintf("NSell\nYПоток\nI1.00\nQ%0.2f\nT%0.2f\nO.00\n", val, val)
		case "проценты":
			return fmt.Sprintf("NIntInc\nYПоток\nLInterest Income\nT%0.2f\nO0.00\n", val)
		default:
			panic(fmt.Sprintf("Unknown op %q", op))
		}
	}},
	).Parse(qifTemplate))

var addToBanktivityParams struct {
	path  string
	from  timeFlag
	to    timeFlag
	group report.GroupType
}

func addAddToBanktivityCommand(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command("add-to-banktivity", "Generate investment report aggregated by operation type and add it to Banktivity.")
	cmd.Arg(
		"report-file",
		"Report file to aggregate.",
	).Required().ExistingFileVar(&addToBanktivityParams.path)
	cmd.Flag(
		"from",
		"Start of the aggregation date range (closed).",
	).Short('f').SetValue(&addToBanktivityParams.from)
	cmd.Flag(
		"to",
		"End of the aggregation date range (open).",
	).Short('t').SetValue(&addToBanktivityParams.to)

	return cmd
}

// addToBanktivity creates an aggregate of the given report by operation for the given
// time range and adds it to Banktivity.
func addToBanktivity() error {
	r, err := reportFromFile(addToBanktivityParams.path)
	if err != nil {
		return err
	}

	if !addToBanktivityParams.from.time().IsZero() || !addToBanktivityParams.to.time().IsZero() {
		r = r.Slice(addToBanktivityParams.from.time(), addToBanktivityParams.to.time())
	}

	agg := r.AggregateByOperation(report.Daily)
	f, err := ioutil.TempFile("", "astream-*.qif")
	if err != nil {
		return errors.Wrap(err, "ioutil.TempFile")
	}
	defer os.Remove(f.Name())

	if err = qif.Execute(f, agg); err != nil {
		f.Close()
		return errors.Wrap(err, "template.Execute")
	}

	if err = f.Close(); err != nil {
		return errors.Wrap(err, "os.File.Close")
	}

	return errors.Wrap(
		exec.Command("open", "-a", "Banktivity 7", "-W", f.Name()).Run(),
		"exec.Command.Run",
	)
}
