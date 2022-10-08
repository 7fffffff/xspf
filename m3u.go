package main

import "io"

type m3uWriter struct {
	wr io.Writer
}

func (m *m3uWriter) Begin() error {
	_, err := io.WriteString(m.wr, "#EXTM3U\n")
	if err != nil {
		return err
	}
	return nil
}

func (m *m3uWriter) WriteFile(path string) error {
	_, err := io.WriteString(m.wr, path+"\n")
	if err != nil {
		return err
	}
	return nil
}
