package util

import "fmt"

func ParseInt64(arr []byte) (res int64, err error) {
	var mult int64 = 1

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] < '0' || arr[i] > '9' {
			err = fmt.Errorf("failed to parse %s to int", string(arr))
			return
		}
		res += int64(arr[i]-'0') * mult
		mult *= 10
	}
	return
}

func ParseInt64NoError(arr []byte) (res int64) {
	var mult int64 = 1

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == '-' {
			return -res
		}
		if arr[i] < '0' || arr[i] > '9' {
			return
		}
		res += int64(arr[i]-'0') * mult
		mult *= 10
	}
	return
}

func ParseInt64ArrNoError(arr []byte) (res []int64) {
	start := -1
	for i := 0; i < len(arr); i++ {
		if start == -1 && ('0' <= arr[i] && arr[i] <= '9') || arr[i] == '-' {
			start = i
		}
		if start > -1 && ('0' > arr[i] || arr[i] > '9') {
			res = append(res, ParseInt64NoError(arr[start:i]))
			start = -1
		}
	}
	if start > -1 {
		res = append(res, ParseInt64NoError(arr[start:]))
	}

	return res
}

func ParseUint8NoError(arr []byte) (res uint8) {
	var mult uint8 = 1

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] < '0' || arr[i] > '9' {
			return
		}
		res += (arr[i] - '0') * mult
		mult *= 10
	}
	return
}

func ParseUint64ArrNoError(arr []byte) (res []uint64) {
	start := -1
	for i := 0; i < len(arr); i++ {
		if start == -1 && '0' <= arr[i] && arr[i] <= '9' {
			start = i
		}
		if start > -1 && ('0' > arr[i] || arr[i] > '9') {
			res = append(res, ParseUint64NoError(arr[start:i]))
			start = -1
		}
	}
	if start > -1 {
		res = append(res, ParseUint64NoError(arr[start:]))
	}

	return res
}

func ParseUint64NoError(arr []byte) (res uint64) {
	var mult uint64 = 1

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] < '0' || arr[i] > '9' {
			return
		}
		res += uint64(arr[i]-'0') * mult
		mult *= 10
	}
	return
}

func ParseUint64IgnoreAll(arr []byte) (res uint64) {
	var mult uint64 = 1

	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] < '0' || arr[i] > '9' {
			continue
		}
		res += uint64(arr[i]-'0') * mult
		mult *= 10
	}
	return
}
