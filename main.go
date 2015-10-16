package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var filename string

func init() {
	flag.StringVar(&filename, "f", "", "Specify the CrdLog.log file to read")
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func decode(v string) (string, error) {
	slice := strings.Split(v, "2 D 0 0 6")
	if len(slice) > 1 {
		if len(slice[1]) == 0 {
			return "", errors.New("Empty line")
		}

		var d byte
		enc := bytes.NewBufferString(slice[1])
		dec := bytes.Buffer{}
		for {
			_, err := fmt.Fscanf(enc, " %X", &d)
			_ = dec.WriteByte(d)
			if err != nil {
				break
			}

		}
		return string(dec.Bytes()), nil
	}
	return "", errors.New("Bad format")
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	lLen := len(line)

	if lLen == 0 {
		return "", err
	}

	if line[lLen-2] == '\r' {
		return line[:lLen-2], err
	}

	return line[:lLen-1], err
}

func main() {
	flag.Parse()

	if len(filename) == 0 {
		fmt.Fprintln(os.Stderr, "Missing CrdLog file")
		os.Exit(-1)
	}

	stream, err := os.Open(filename)
	exitOnError(err)
	defer stream.Close()

	reader := bufio.NewReader(stream)
	for {
		line, err := readLine(reader)
		if err != nil {
			break
		}
		if strings.Contains(line, ":2 D 0 0") {
			str, err := decode(line)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			del := str[3]
			field := strings.Split(str, string(del))

			fmt.Println("------")
			fmt.Println(field[0])
			for _, v := range field[1:] {
				fmt.Println("[FS]", v)
			}
			fmt.Println("------")
		}
	}
}
