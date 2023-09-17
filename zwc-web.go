// Copyright (C) 2023 Ethan Cheng <ethanrc0528@gmail.com>
//
// This file is part of zwc-web.
//
// zwc-web is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public License as published by the Free
// Software Foundation, version 3 of the License.
//
// zwc-web is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with zwc-web. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"syscall/js"
	"strings"

	"github.com/yadayadajaychan/zwc"
)

func encodeWrapper() js.Func {
	encodeFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0].String()

		data, message, ok := strings.Cut(input, "\nEOF\n")
		if !ok {
			return "error: no EOF separator found"
		}
		if len(data) == 0 {
			return "error: no data supplied"
		}
		if len(message) == 0 {
			return "error: no message supplied"
		}

		encoding := zwc.NewEncoding(1, 4, 16)
		dst := make([]byte, encoding.EncodedMaxLen(len(data)))

		n := encoding.Encode(dst, []byte(data))

		return message[:1] + string(dst[:n]) + message[1:]
	})
	return encodeFunc
}

func decodeWrapper() js.Func {
	decodeFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0].String()

		v, e, c, err := zwc.DecodeHeader([]byte(input))
		if err != nil {
			return "error: " + err.Error()
		}

		encoding := zwc.NewEncoding(v, e, c)
		dst := make([]byte, encoding.DecodedPayloadMaxLen(len(input)))

		i := strings.IndexRune(input, zwc.V1DelimChar) + len(zwc.V1DelimCharUTF8)
		i = i + strings.IndexRune(input[i:], zwc.V1DelimChar) + len(zwc.V1DelimCharUTF8)
		n, _, err := encoding.Decode(dst, []byte(input[i:]))
		if err != nil {
			return string(dst[:n]) + "\nerror: " + err.Error()
		}

		return string(dst[:n])
	})
	return decodeFunc
}

func main() {
	js.Global().Set("encode", encodeWrapper())
	js.Global().Set("decode", decodeWrapper())
	<-make(chan bool)
}
