package scrapper

import "io"

type Detail struct {
	Name         string 
	Mobile       string
	CompanyName  string
	Category     string
	VisitCardUrl string
	Email        string
}

type Response interface {
	io.Reader
}

type HTMLResponse struct {
	Data []byte
}

// Read method to implement io.Reader interface
func (r *HTMLResponse) Read(p []byte) (n int, err error) {
	copy(p, r.Data)
	if len(r.Data) > len(p) {
		r.Data = r.Data[len(p):]
		return len(p), nil
	}
	n = len(r.Data)
	r.Data = nil
	return n, io.EOF
}

var DetailsChannel = make(chan Detail)
var HTMLChannel = make(chan []byte)

