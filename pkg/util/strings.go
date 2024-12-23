package util

import "unsafe"

func BytesEqualString(b []byte, s string) bool {
	return ToUnsafeString(b) == s
}

func ToUnsafeString(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}
