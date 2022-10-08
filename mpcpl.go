package main

import (
	"io"
	"strconv"
)

type mpcplWriter struct {
	wr    io.Writer
	count int
}

func (m *mpcplWriter) Begin() error {
	_, err := io.WriteString(m.wr, "MPCPLAYLIST\n")
	if err != nil {
		return err
	}
	return nil
}

func (m *mpcplWriter) WriteFile(path string) error {
	m.count++
	countStr := strconv.Itoa(m.count)
	_, err := io.WriteString(m.wr, "\n"+countStr+",type,0\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(m.wr, countStr+",filename,"+path+"\n")
	if err != nil {
		return err
	}
	return nil
}
