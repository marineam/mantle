package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/coreos/mantle/network/journal"
)

type flakyJournal struct {
	cmd *exec.Cmd
	out io.Reader
}

func FlakyJournal(cursor string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd := exec.CommandContext(ctx, "journalctl", "--output=export", "--follow")
	if cursor == "" {
		cmd.Args = append(cmd.Args, "--boot", "--lines=all")
	} else {
		cmd.Args = append(cmd.Args, "--after-cursor", cursor)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	log.Printf("Starting: %s", cmd.Args)
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	return &flakyJournal{cmd, out}, nil
}

func (f *flakyJournal) Read(b []byte) (int, error) {
	n, err := f.out.Read(b)
	if err == io.EOF {
		if err2 := f.cmd.Wait(); err2 != nil {
			log.Printf("Stopped journalctl: %v", err2)
		}
	}
	return n, err
}

func main() {
	log.SetPrefix("       ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	w := journal.NewShortWriter(os.Stdout)
	cursor := ""

	for {
		j, err := FlakyJournal(cursor)
		if err != nil {
			log.Fatal(err)
		}
		r := journal.NewExportReader(j)
		for {
			entry, err := r.ReadEntry()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			cursor = string(entry[journal.FIELD_CURSOR])
			if err := w.WriteEntry(entry); err != nil {
				log.Fatal(err)
			}
		}
	}
}
