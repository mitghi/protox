package server

import "testing"

const (
	cERR string = "inconsistent state, expected err==nil."
)

func TestGenerateTLSConfig(t *testing.T) {
	var s *Server = NewServer()
	_, err := s.generateTLSConfig(&TLSOptions{
		Cert: "../config/cert/server.pem",
		Key:  "../config/cert/key.pem",
	})
	if err != nil {
		t.Fatal(cERR, err)
	}
}
