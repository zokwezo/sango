package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

const (
	LoPitch = iota
	MidPitch
	HiPitch
)

func main() {
	ch := make(chan byte)
	go func(ch chan<- byte) {
		// Disable input buffering
		err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		if err != nil {
			panic(err)
		}
		// Do not display entered characters on the screen
		err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		if err != nil {
			panic(err)
		}
		var b []byte = make([]byte, 1)
		for {
			n, err := os.Stdin.Read(b)
			if n == 0 && errors.Is(err, io.EOF) {
				close(ch)
				return
			}
			if err != nil {
				panic(err)
			}
			for _, c := range b {
				if c == 4 {
					// Re-display entered characters on the screen
					err := exec.Command("stty", "-F", "/dev/tty", "echo").Run()
					if err != nil {
						panic(err)
					}
					close(ch)
					return
				}
				ch <- c
			}
		}
	}(ch)

	pitch := LoPitch
	for c := range ch {
		switch c {
		case 'j', 'J':
			pitch++
		default:
			b := []byte{c}
			if c, found := asciiToUtf8[pitch%3][c]; found {
				b = c
			}
			pitch = LoPitch
			if _, err := os.Stdout.Write(b); err != nil {
				panic(err)
			}
		}
	}
}

var asciiToUtf8 = map[int]map[byte][]byte{
	LoPitch: {
		'A': ([]byte)("A"),
		'E': ([]byte)("E"),
		'X': ([]byte)("Ɛ"),
		'I': ([]byte)("I"),
		'O': ([]byte)("O"),
		'C': ([]byte)("Ɔ"),
		'U': ([]byte)("U"),
		'a': ([]byte)("a"),
		'e': ([]byte)("e"),
		'x': ([]byte)("ɛ"),
		'i': ([]byte)("i"),
		'o': ([]byte)("o"),
		'c': ([]byte)("ɔ"),
		'u': ([]byte)("u"),
	},
	MidPitch: {
		'A': ([]byte)("Ä"),
		'E': ([]byte)("Ë"),
		'X': ([]byte)("Ɛ̈"),
		'I': ([]byte)("Ï"),
		'O': ([]byte)("Ö"),
		'C': ([]byte)("Ɔ̈"),
		'U': ([]byte)("Ü"),
		'a': ([]byte)("ä"),
		'e': ([]byte)("ë"),
		'x': ([]byte)("ɛ̈"),
		'i': ([]byte)("ï"),
		'o': ([]byte)("ö"),
		'c': ([]byte)("ɔ̈"),
		'u': ([]byte)("ü"),
	},
	HiPitch: {
		'A': ([]byte)("Â"),
		'E': ([]byte)("Ê"),
		'X': ([]byte)("Ɛ̂"),
		'I': ([]byte)("Î"),
		'O': ([]byte)("Ô"),
		'C': ([]byte)("Ɔ̂"),
		'U': ([]byte)("Û"),
		'a': ([]byte)("â"),
		'e': ([]byte)("ê"),
		'x': ([]byte)("ɛ̂"),
		'i': ([]byte)("î"),
		'o': ([]byte)("ô"),
		'c': ([]byte)("ɔ̂"),
		'u': ([]byte)("û"),
	},
}
