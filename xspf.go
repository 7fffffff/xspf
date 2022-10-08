package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
)

type xspfWriter struct {
	wr  io.Writer
	enc *xml.Encoder
}

func newXSPFWriter(w io.Writer) *xspfWriter {
	return &xspfWriter{
		wr: w,
	}
}

func (x *xspfWriter) Begin() error {
	_, err := io.WriteString(x.wr, xml.Header)
	if err != nil {
		return err
	}
	x.enc = xml.NewEncoder(x.wr)
	x.enc.Indent("", "  ")
	err = x.enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "playlist"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: "1"},
			{Name: xml.Name{Local: "xmlns"}, Value: "http://xspf.org/ns/0/"},
			{Name: xml.Name{Local: "xmlns:vlc"}, Value: "http://www.videolan.org/vlc/playlist/ns/0/"},
		},
	})
	if err != nil {
		return err
	}
	err = x.enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "trackList"},
		Attr: []xml.Attr{},
	})
	if err != nil {
		return err
	}
	return nil
}

func (x *xspfWriter) End() error {
	err := x.enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: "trackList"},
	})
	if err != nil {
		return err
	}
	err = x.enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: "playlist"},
	})
	if err != nil {
		return err
	}
	err = x.enc.Flush()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(x.wr)
	return err
}

func (x *xspfWriter) WriteTrack(filePath string) error {
	err := x.enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "track"},
	})
	if err != nil {
		return err
	}
	err = x.enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "location"},
	})
	if err != nil {
		return err
	}
	fileURL := url.URL{
		Scheme: "file",
		Path:   filePath,
	}
	err = x.enc.EncodeToken(xml.CharData(fileURL.String()))
	if err != nil {
		return err
	}
	err = x.enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: "location"},
	})
	if err != nil {
		return err
	}
	err = x.enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: "track"},
	})
	return err
}
