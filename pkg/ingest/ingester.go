package ingest

import (
	"fmt"
	"io"
	"os"

	"github.com/ryszard/tfutils/go/tfrecord"
)

// Ingester is an interface for ingesting TFRecord file
type Ingester interface {
	// FetchUpdates fetches updates from the TFRecord file incrementally.
	// User should call this function periodically in order to get the latest updates.
	// return error means unexpected file content, user should close the Ingester.
	// Notes: FetchUpdates is not thread-safe.
	FetchUpdates() error

	Close() error

	GetRun() *Run
}

type ingester struct {
	file   *os.File
	parser *parser
	run    *Run
}

func (i *ingester) FetchUpdates() error {
	for {
		offset, err := i.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("seek fail: %w", err)
		}
		data, err := tfrecord.Read(i.file)
		if err != nil {
			_, err = i.file.Seek(offset, io.SeekStart)
			if err != nil {
				return fmt.Errorf("seek fail: %w", err)
			}
			// no more data
			break
		}
		err = i.parser.ParseRecord(data, i.run.AddScalarEvent)
		if err != nil {
			return fmt.Errorf("parse fail: %w", err)
		}
	}
	return nil
}

func (i *ingester) Close() error {
	return i.file.Close()
}

func (i *ingester) GetRun() *Run {
	return i.run
}

func NewIngester(name string, path string) (Ingester, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	run := NewRun(name)
	parser := NewParser()
	return &ingester{file: file, parser: parser, run: run}, nil
}
