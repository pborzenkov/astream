package report

import (
	"io"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

type transaction struct {
	src      string    // source account
	dst      string    // destination account
	date     time.Time // date of the transaction
	amount   float64   // sum of the transaction
	ttype    string    // transaction type
	borrower string    // business or person that borrowed the money
	contract string    // contract number
}

// Report contains parsed report of AlfaStream investments.
type Report struct {
	transactions []*transaction // all transactions in the file, sorted by date
}

// NewFromXLSX returns an instance of a Report parsed from an XLSX file pointed
// to by r.
func NewFromXLSX(r io.ReaderAt, size int64) (*Report, error) {
	fil, err := xlsx.OpenReaderAt(r, size)
	if err != nil {
		return nil, errors.Wrap(err, "xlsx.OpenReaderAt")
	}

	sheet := fil.Sheet["Статистика"]
	if sheet == nil {
		return nil, errors.Errorf("malformed report, no sheet named 'Статистика'")
	}

	columns := make(map[string]int)
	// NOTE: what to do if it's not the first line???
	for i, c := range sheet.Rows[0].Cells {
		columns[c.Value] = i
	}
	sheet.Rows = sheet.Rows[1:]

	// first check that all cells contain valid data
	for i, r := range sheet.Rows {
		rnum := i + 2 // one deleted, one because loop starts from 0

		if _, err := r.Cells[columns["Сумма"]].Float(); err != nil {
			return nil, errors.Wrapf(err, "row(%d).Cell[\"Сумма\"].Float", rnum)
		}
		if _, err := r.Cells[columns["Дата"]].GetTime(false); err != nil {
			return nil, errors.Wrapf(err, "row(%d).Cell[\"Дата\"].GetTime", rnum)
		}
	}

	sort.Slice(sheet.Rows, func(i, j int) bool {
		idate, _ := sheet.Rows[i].Cells[columns["Дата"]].GetTime(false)
		jdate, _ := sheet.Rows[j].Cells[columns["Дата"]].GetTime(false)

		return idate.Before(jdate)
	})

	rep := &Report{
		transactions: make([]*transaction, 0, len(sheet.Rows)),
	}
	for _, r := range sheet.Rows {
		d, _ := r.Cells[columns["Дата"]].GetTime(false)
		a, _ := r.Cells[columns["Сумма"]].Float()

		rep.addTransaction(&transaction{
			src:      r.Cells[columns["Номер счета инвестора"]].Value,
			dst:      r.Cells[columns["Номер счета заемщика"]].Value,
			date:     d,
			amount:   a,
			ttype:    r.Cells[columns["Тип операции"]].Value,
			borrower: r.Cells[columns["Наименование заемщика"]].Value,
			contract: r.Cells[columns["Номер договора"]].Value,
		})
	}

	return rep, nil
}

func (r *Report) addTransaction(t *transaction) {
	r.transactions = append(r.transactions, t)
}

// Slice creates a slice of the report that contains transactions from the
// given time range.
func (r *Report) Slice(from, to time.Time) *Report {
	nr := &Report{}

	from = from.Truncate(24 * time.Hour)
	to = to.Truncate(24 * time.Hour)

	for _, t := range r.transactions {
		if (t.date.Equal(from) || t.date.After(from)) && (to.IsZero() || t.date.Before(to)) {
			nr.addTransaction(t)
		}
	}

	return nr
}
