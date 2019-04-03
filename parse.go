// +build go1.11

package passwd

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

const (
	separatorChar = '$'
	separatorRune = rune('$')
	separatorStr  = "$"
)

var (
	rangeTableSeparator = rangetable.New(separatorRune)
)

func token(c rune) bool {
	return unicode.Is(rangeTableSeparator, c)
}

func parseFromHashToParams(hashed []byte) (interface{}, error) {

	fields := strings.FieldsFunc(string(hashed), token)
	//fmt.Printf("FIELDS: %q\n", fields)

	switch fields[0] {
	case idBcrypt:
		//fmt.Printf("bcrypt compare!\n")
		bp := BcryptParams{}
		//return bp.Compare(hashed, password)
		return bp, nil
	case idScrypt:
		//fmt.Printf("scrypt compare!\n")
		sp, err := newScryptParamsFromFields(fields[1:]) // mismatch.
		if err != nil {
			// XXX wrapp the error
			return nil, err
		}
		return *sp, nil
	case idArgon2i:
		fallthrough
	case idArgon2id:
		//fmt.Printf("argon2id compare!\n")
		ap, err := newArgon2ParamsFromFields(fields[1:]) // mismatch.
		if err != nil {
			// XXX wrapp the error
			return nil, err
		}
		//return ap.Compare(hashed, password)
		return *ap, nil
	}
	return nil, fmt.Errorf("invalid")
}

func parseFromHashToSalt(hashed []byte) ([]byte, error) {
	fields := strings.FieldsFunc(string(hashed), token)
	if len(fields) == 0 {
		return nil, fmt.Errorf("invalid format")
	}
	fmt.Printf("prout fields: %q\n", fields)
	switch fields[0] {
	case idBcrypt:
		return nil, nil
	case idScrypt:
		fallthrough
	case idArgon2i:
		fallthrough
	case idArgon2id:
		salt, err := base64Decode([]byte(fields[1])) // process the salt
		if err != nil {
			return nil, err
		}
		return salt, nil
	}
	return nil, fmt.Errorf("invalid format")

}
