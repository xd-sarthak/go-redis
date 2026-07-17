package core

import "strconv"

func deduceTypeEncoding(value string) (uint8, uint8) {
	objType := OBJ_TYPE_STRING
	if _,err := strconv.ParseInt(value, 10, 64); err == nil {
		return objType, OBJ_ENCODING_INT
	}
	if len(value) <= 44 {
		return objType, OBJ_ENCODING_EMBSTR
	}
	return objType, OBJ_ENCODING_RAW
}