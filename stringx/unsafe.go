package stringx

import "unsafe"

// UnsafeStringToBytes converts a string to a byte slice.
//
// This is a shallow copy, meaning the returned byte slice
// reuses the underlying array of the string. You must not
// modify the returned byte slice under any circumstances.
//
// # Safety
//
//   - The returned slice aliases the string's read-only storage.
//     Writing through it is undefined behavior and, on Go tip,
//     typically corrupts interned/constant strings or crashes.
//   - Do NOT pass the returned slice to any API that may retain
//     or mutate it (e.g., append into it, io.Reader.Read, bufio,
//     crypto hashers with scratch buffer). Prefer []byte(s) when
//     the callee's retention/mutation contract is unknown.
//   - Keep the source string alive for as long as the slice is
//     in use. The slice header does not pin the string header,
//     so a dead string may have its storage reclaimed while the
//     slice is still referenced.
//   - Using the slice as a map key has implementation-defined
//     hashing behavior and is not supported; use the string
//     directly as the map key instead.
//   - Only use this when profiling shows the copy of []byte(s) is
//     a measurable bottleneck on a hot path. In all other cases
//     the safe conversion []byte(s) is the correct choice.
func UnsafeStringToBytes(s string) []byte {
	// unsafe.StringData is unspecified for the empty string, so we provide a strict interpretation
	if len(s) == 0 {
		return nil
	}
	// Copied from go 1.20.1 os.File.WriteString
	// https://github.com/golang/go/blob/202a1a57064127c3f19d96df57b9f9586145e21c/src/os/file.go#L246
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// func UnsafeStringToBytes(s string) []byte {
// 	if len(s) == 0 {
// 		return nil
// 	}
// 	x := (*[2]uintptr)(unsafe.Pointer(&s))
// 	h := [3]uintptr{x[0], x[1], x[1]}
// 	return *(*[]byte)(unsafe.Pointer(&h))
// }

// func UnsafeStringToBytes(s string) (b []byte) {
// 	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
// 	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
// 	bh.Data = sh.Data
// 	bh.Len = sh.Len
// 	bh.Cap = sh.Len
// 	return b
// }

// UnsafeBytesToString converts a byte slice to a string.
//
// This is a shallow copy, meaning the returned string reuses the
// underlying array of the byte slice. It is your responsibility
// to keep the input byte slice alive until you no longer access
// the string.
//
// # Safety
//
//   - Strings in Go are immutable. Once this function returns, any
//     further mutation of b (direct writes, append that reuses the
//     backing array, io.Reader.Read into b, encoding/encrypted
//     in-place operations, etc.) silently violates the immutability
//     invariant. Callers MUST treat b as frozen for the lifetime of
//     every returned string that still refers to it.
//   - Do NOT use the returned string as a map key or as a key to
//     any hash-based structure unless b is guaranteed never to
//     change. Once a key's underlying bytes are mutated, lookups
//     become undefined (the new hash differs from the stored hash).
//   - Do NOT let the returned string escape to code you do not
//     control (callback, goroutines with unclear lifetime, cached
//     values, returned values) unless you can guarantee b outlives
//     every such reference. When in doubt, copy via string(b).
//   - If b was obtained from a pool (sync.Pool, bytes.Buffer reuse,
//     bufio.Reader.ReadSlice) returning b to the pool while the
//     string is still observable is a use-after-free.
//   - Only use this when profiling shows the copy of string(b) is a
//     measurable bottleneck on a hot path. In all other cases the
//     safe conversion string(b) is the correct choice.
func UnsafeBytesToString(b []byte) string {
	// unsafe.SliceData relies on cap whereas we want to rely on len
	if len(b) == 0 {
		return ""
	}
	// Copied from go 1.20.1 strings.Builder.String
	// https://github.com/golang/go/blob/202a1a57064127c3f19d96df57b9f9586145e21c/src/strings/builder.go#L48
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// func UnsafeBytesToString(b []byte) string {
// 	if len(b) == 0 {
// 		return ""
// 	}
// 	return *(*string)(unsafe.Pointer(&b))
// }
