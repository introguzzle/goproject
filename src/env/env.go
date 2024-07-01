package env

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Key struct {
	Value string
	Valid bool
}

func Get(key string) Key {
	for _, s := range format() {
		value := strings.Split(s, "=")

		if value != nil && key == value[0] {
			return Key{
				Value: trim(value[1]),
				Valid: true,
			}
		}
	}

	return Key{
		Value: "",
		Valid: false,
	}
}

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(s, "\n"), "\""))
}

func format() []string {
	return strings.Split(string(read()), "\n")
}

func read() []byte {
	f, err := os.OpenFile(".env", os.O_RDONLY, 0666)
	defer func(f *os.File) {
		if f != nil {
			_ = f.Close()
		}
	}(f)

	if err != nil {
		fErr("Error opening file", err)
		return nil
	}

	data := make([]byte, 0)
	buffer := make([]byte, 64)

	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			fErr("Error reading file", err)
			return nil
		}

		if n == 0 {
			break
		}

		data = append(data, buffer[:n]...)
	}

	return data
}

func fErr(s string, err error) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s", s, err))
}
