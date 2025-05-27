package stringx

import "unsafe"

// StringToBytes converts a string to a byte slice.
//
// This is a shallow copy, means that the returned byte slice reuse
// the underlying array in string, so you can't change the returned
// byte slice in any situations.
func StringToBytes(s string) []byte {
	// unsafe.StringData is unspecified for the empty string, so we provide a strict interpretation
	if len(s) == 0 {
		return nil
	}
	// Copied from go 1.20.1 os.File.WriteString
	// https://github.com/golang/go/blob/202a1a57064127c3f19d96df57b9f9586145e21c/src/os/file.go#L246
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// func StringToBytes(s string) []byte {
// 	if len(s) == 0 {
// 		return nil
// 	}
// 	x := (*[2]uintptr)(unsafe.Pointer(&s))
// 	h := [3]uintptr{x[0], x[1], x[1]}
// 	return *(*[]byte)(unsafe.Pointer(&h))
// }

// func StringToBytes(s string) (b []byte) {
// 	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
// 	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
// 	bh.Data = sh.Data
// 	bh.Len = sh.Len
// 	bh.Cap = sh.Len
// 	return b
// }

// BytesToString converts a byte slice to a string.
//
// This is a shallow copy, means that the returned string reuse the
// underlying array in byte slice, it's your responsibility to keep
// the input byte slice survive until you don't access the string anymore.
func BytesToString(b []byte) string {
	// unsafe.SliceData relies on cap whereas we want to rely on len
	if len(b) == 0 {
		return ""
	}
	// Copied from go 1.20.1 strings.Builder.String
	// https://github.com/golang/go/blob/202a1a57064127c3f19d96df57b9f9586145e21c/src/strings/builder.go#L48
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// func BytesToString(b []byte) string {
// 	if len(b) == 0 {
// 		return ""
// 	}
// 	return *(*string)(unsafe.Pointer(&b))
// }
