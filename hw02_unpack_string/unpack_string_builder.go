package hw02unpackstring

import (
	"strings"
)

type unpackStringBuilder struct {
	*strings.Builder
	symbol rune
	repeat int
}

func newUnpackStringBuilder() unpackStringBuilder {
	return unpackStringBuilder{
		&strings.Builder{},
		0,
		0,
	}
}

func (usb *unpackStringBuilder) setSymbol(symbol rune) error {
	if _, err := usb.writeSymbol(); err != nil {
		return err
	}

	usb.symbol = symbol
	usb.repeat = 1

	return nil
}

func (usb *unpackStringBuilder) writeSymbol() (string, error) {
	if usb.repeat > 0 && usb.symbol > 0 {
		_, err := usb.WriteString(strings.Repeat(string(usb.symbol), usb.repeat))
		if err != nil {
			return "", err
		}
	}

	return usb.String(), nil
}
