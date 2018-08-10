package main

import (
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pborzenkov/astream/report"
	"github.com/pkg/errors"
)

func reportFromFile(path string) (*report.Report, error) {
	fil, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}
	defer fil.Close()

	st, err := fil.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "fil.Stat")
	}

	r, err := report.NewFromXLSX(fil, st.Size())
	if err != nil {
		return nil, errors.Wrap(err, "report.NewFromXLSX")
	}

	return r, nil
}

type timeFlag time.Time

func (tf *timeFlag) Set(value string) error {
	t, err := dateparse.ParseStrict(value)
	if err != nil {
		return errors.Wrap(err, "dateparse.ParseStrict")
	}

	*(*time.Time)(tf) = t
	return nil
}

func (tf *timeFlag) String() string {
	return tf.String()
}

func (tf *timeFlag) time() time.Time {
	return *(*time.Time)(tf)
}
