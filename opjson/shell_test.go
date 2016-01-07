package opjson

import (
	"encoding/json"
	"log"
	"testing"
)

func TestSFirst(t *testing.T) {
	log.Println("TEST SFIRST")
}

func TestShell(t *testing.T) {
	sh, ok := MakeShell(nil)
	log.Println("MADE NILL SHELL:", sh, ok)
	log.Println(string(sh.Data))
	raw, err := json.Marshal(sh)
	if err != nil {
		log.Println("ERROR MARSHALLING SHELL:", err)
		return
	}
	log.Println("MARSHALLED:", string(raw))
}
