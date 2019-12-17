/*
 * @Description: eos-base32
 * @Copyright: UMU618
 * @Author: UMU618 <umu618@hotmail.com>
 * @Date: 2019-12-17 17:42:13
 * @LastEditTime: 2019-12-17 17:42:13
 * @LastEditors: UMU618
 */
package main // import "github.com/UMU618/eos-base32"

import (
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
)

// encoding/base32 is NOT suitable for decode/encode EOS name
const eosBase32 = ".12345abcdefghijklmnopqrstuvwxyz"

func main() {
	inputEosName := flag.String("decode", "", "Base32 encoded EOS name.")
	intputUInt64 := flag.Uint64("encode", 0, "UInt64.")
	flag.Parse()
	if len(*inputEosName) > 0 {
		name := *inputEosName
		if len(name) > 13 {
			fmt.Println("Length", len(name), "> 13")
			return
		} else if len(name) == 13 && name[12] != '.' {
			fmt.Printf("The 13th charactor is '%c', result may wrong.\n", name[12])
			return
		}
		for len(name) < 13 {
			name += "."
		}
		enc := base32.NewEncoding(eosBase32).WithPadding(base32.NoPadding)
		b := make([]byte, enc.DecodedLen(len(name)+1))
		n, err := enc.Decode(b, []byte(name))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		if n != 8 {
			fmt.Println("Length unmatched!", n, len(b))
			return
		}
		var u uint64 = 0
		for i := 0; i < n; i++ {
			u <<= 8
			u |= uint64(b[i])
		}
		fmt.Printf("Decode(%s) = 0x%x, %d\n", name, u, u)
	} else if *intputUInt64 > 0 {
		enc := base32.NewEncoding(eosBase32).WithPadding(base32.NoPadding)
		u := *intputUInt64
		if (u & 15) != 0 {
			fmt.Printf("The lowest 4 bits is 0x%x, result may wrong.\n", u&15)
			return
		}
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, u)
		name := enc.EncodeToString(b)
		fmt.Printf("Encode(0x%0x, %d) = %s\n", u, u, name)
	} else {
		flag.Usage()
	}
}
