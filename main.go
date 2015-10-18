package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var filename string

func init() {
	flag.StringVar(&filename, "f", "", "Specify the CrdLog.log file to read")
}

// helper function to exit application with a message to Stderr.
func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// parse the line and return only the response of the transaction.
func fetchRawData(v string) string {
	slice := strings.Split(v, "2 D 0 0 6")
	return slice[len(slice)-1]
}

// convert the hex value of the response into it's string representation
func convertHexToASCII(v string) (string, error) {
	hex := bytes.NewBufferString(v)
	ascii := bytes.Buffer{}
	var d byte
	for {
		_, err := fmt.Fscanf(hex, " %X", &d)

		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		err = ascii.WriteByte(d)
		if err != nil {
			return "", err
		}
	}
	return string(ascii.Bytes()), nil
}

// decode take a line from the CrdtLog.Log file and parse the reponse
// into a human readable format
func decode(v string) ([]string, error) {

	data := fetchRawData(v)

	if len(data) <= 2 {
		return nil, errors.New("Empty line")
	}

	data, err := convertHexToASCII(data)
	if err != nil {
		return nil, errors.New("Bad format")
	}

	//remove the two hidden characters at the end of the line
	if byte(data[len(data)-1]) == 0x0 {
		data = data[:len(data)-1]
	}

	if byte(data[len(data)-1]) == 0x4E {
		data = data[:len(data)-1]
	}

	delim := 0x1C
	field := strings.Split(data, string(delim))
	return field, nil
}

// readNextLine break the file stream into line of
func readNextLine(reader *bufio.Reader) (string, error) {
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
	i := 0
	for {
		line, err := readNextLine(reader)
		if err != nil {
			break
		}
		if strings.Contains(line, ":2 D 0 0") {
			i++
			fields, err := decode(line)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(i, ": ------")
			fmt.Println(fields[0])
			for _, v := range fields[1:] {
				fmt.Println("[FS]", v)
			}
			fmt.Println(i, ": ------")
		}
	}
}
