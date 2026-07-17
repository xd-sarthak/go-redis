package core

import "errors"

func getType(te uint8) uint8 {
	// first 4 bits
	return te & 0xF0
}

func getEncoding(te uint8) uint8 {
	// last 4 bits
	return te & 0x0F
}


func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return nil
}

func assertEncoding(te uint8, e uint8) error {
	if getEncoding(te) != e {
		return errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return nil
}