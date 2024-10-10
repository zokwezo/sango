package transcode

import (
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	testCases = map[string][]SSE{
		// English has the wrong lengths (falsely including symbols)
		// French doesn't work at all.
		"Hello":   {0x4548, 0x4465, 0x436c, 0x426c, 0x416f},
		"Bɛ̂-bïn": {0xbed3, 0x94dd},
		"bebi":    {0xa8db, 0x88d5},
	}
)

func TestEncode(t *testing.T) {
	for input, expect := range testCases {
		log.Printf("input = %q\n", input)
		in := strings.NewReader(input)
		actual := encode(in)
		nActual := len(actual)
		nExpect := len(expect)
		if nActual != nExpect {
			log.Printf("input       = %q\n", input)
			log.Printf("len(actual) = %v\n", nActual)
			log.Printf("len(expect) = %v\n", nExpect)
		}
		for k := range max(nActual, nExpect) {
			if k < nActual {
				if k < nExpect {
					if actual[k] != expect[k] {
						log.Printf("input      = %q\n", input)
						log.Printf("actual[%v] = %04x = %016b\n", k, actual[k], actual[k])
						log.Printf("expect[%v] = %04x = %016b\n", k, expect[k], expect[k])
					}
				} else {
					log.Printf("input      = %q\n", input)
					log.Printf("actual[%v] = %04x = %016b\n", k, actual[k], actual[k])
					log.Printf("expect[%v] not defined\n", k)
				}
			} else if k < nExpect {
				log.Printf("input      = %q\n", input)
				log.Printf("actual[%v] not defined\n", k)
				log.Printf("expect[%v] = %04x = %016b\n", k, expect[k], expect[k])
			}
		}
	}
}
