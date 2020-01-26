package dingo

import (
	"bytes"
	"errors"
	"unicode"
)

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digits = []byte("0123456789")

// FormatDefName is the function used to turn the definition name
// into something that can be used in the generated container.
func FormatDefName(name string) string {
	formatted := bytes.NewBuffer(nil)
	start := true

	for _, c := range name {
		if !bytes.ContainsRune(chars, c) {
			start = true
			continue
		}

		if !start {
			formatted.WriteRune(c)
			continue
		}

		formatted.WriteRune(unicode.ToUpper(c))
		start = false
	}

	return formatted.String()
}

// DefNameIsAllowed returns an error if the definition name is not allowed.
func DefNameIsAllowed(name string) error {
	names := []string{"", "C", "ErrorCallback", "Container", "NewContainer"}

	formatted := FormatDefName(name)

	for _, n := range names {
		if n == formatted {
			return errors.New("DefName '" + name + "' is not allowed (reserved key word)")
		}
	}

	if bytes.ContainsRune(digits, rune(formatted[0])) {
		return errors.New("DefName '" + name + "' is not allowed (first char is a digit)")
	}

	return nil
}

// FormatPkgName formats a package name by keeping only the letters.
func FormatPkgName(name string) string {
	formatted := bytes.NewBuffer(nil)

	for _, c := range []byte(name) {
		if bytes.Contains(letters, []byte{c}) {
			formatted.WriteByte(c)
		}
	}

	return formatted.String()
}
