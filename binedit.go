package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"os"
	"io"
	"encoding/hex"
	"bytes"
)


var (
	doRead = kingpin.Flag("read", "read binary data from stdin, output text representation to stdout").Short('r').Bool()
	doWrite = kingpin.Flag("write", "read text input from stdin, convert to binary file").Short('w').Bool()
	bytesPerLine = kingpin.Flag("bytes-per-line", "number of bytes to output per one line").Short('b').Default("16").Int()
)

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	kingpin.Parse()

	if (*doRead == *doWrite) {
		fmt.Fprintln(os.Stderr, "You must provide exactly one mode of operation (--read or --write)");
		os.Exit(1);
	}

	if (*doRead) {
		line := make([]byte, *bytesPerLine)
		c := 0

		buf := make([]byte, 1)
		for true {
			_, err := os.Stdin.Read(buf)
			if err == io.EOF {
				break
			}
			handleError(err)

			if (c > 0 && (c % *bytesPerLine) == 0) {
				printLine(line, *bytesPerLine)
				for i := range line {
					line[i] = 0
				}
			}

			if (c % *bytesPerLine) > 0 && ((c % *bytesPerLine) % 8) == 0 {
				fmt.Fprintf(os.Stdout, " ")
			}
			fmt.Fprintf(os.Stdout, "%s ", hex.EncodeToString(buf))
			line[c % *bytesPerLine] = buf[0]

			c++;
		}

		if c > 0 {
			printLine(line, c % *bytesPerLine)
		}
	}

	if (*doWrite) {
		buf := make([]byte, 1)
		buffer := new(bytes.Buffer)
		readingData := true
		for true {
			_, err := os.Stdin.Read(buf)
			if err == io.EOF {
				break
			}
			handleError(err)

			if readingData {
				if (buf[0] >= 48 && buf[0] <= 57) || (buf[0] >= 97 && buf[0] <= 102) || (buf[0] >= 65 && buf[0] <= 70) {
					buffer.WriteByte(buf[0])
					if buffer.Len() % 1024 == 0 {
						data, err := hex.DecodeString(buffer.String())
						handleError(err)
						_, err = os.Stdout.Write(data)
						handleError(err)
						buffer = new(bytes.Buffer)
					}
				} else if buf[0] == 124 {
					readingData = false
				}
			} else {
				if buf[0] == 10 {
					readingData = true
				}
			}
		}
		data, err := hex.DecodeString(buffer.String())
		handleError(err)
		_, err = os.Stdout.Write(data)
		handleError(err)
	}
}

func printLine(line []byte, c int) {
	fmt.Fprintf(os.Stdout, "| ")
	for i := 0; i < c; i++ {
		b := line[i]
		if b >= 32 && b <= 126 {
			fmt.Fprintf(os.Stdout, "%c", b)
		} else {
			fmt.Fprintf(os.Stdout, ".")
		}
	}
	fmt.Fprintf(os.Stdout, "\n")
}
