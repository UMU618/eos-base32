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
		var lower uint64 = 0
		if len(name) > 13 {
			fmt.Println("Length", len(name), "> 13")
			return
		} else if len(name) == 13 && name[12] != '.' {
			if name[12] >= '1' && name[12] <= '5' {
				lower = uint64(name[12] - '0')
			} else if name[12] >= 'a' && name[12] <= 'j' {
				lower = uint64(name[12] - 'a' + 6)
			} else {
				fmt.Printf("The 13th charactor is invalid: '%c'.\n", name[12])
				return
			}
		}
		for len(name) < 13 {
			name += "."
		}
		b, err := base32.NewEncoding(eosBase32).WithPadding(base32.NoPadding).DecodeString(name)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		if len(b) != 8 {
			fmt.Printf("Length %d unmatched!\n", len(b))
			return
		}
		var u uint64 = 0
		for i := 0; i < 8; i++ {
			u <<= 8
			u |= uint64(b[i])
		}
		u &^= 15
		u |= lower
		fmt.Printf("Decode(%s) = 0x%x, %d\n", name, u, u)
	} else if *intputUInt64 > 0 {
		u := *intputUInt64
		lower := u & 15
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, u&^15)
		name := base32.NewEncoding(eosBase32).WithPadding(base32.NoPadding).EncodeToString(b)
		if lower > 0 && lower <= 5 {
			name = name[:12] + string('0'+lower)
		} else if lower >= 6 {
			name = name[:12] + string('a'+lower-6)
		}
		fmt.Printf("Encode(0x%0x, %d) = %s\n", u, u, name)
	} else {
		flag.Usage()
	}
}
