package encode

func RemainLength(len uint32) (buf []byte) {
	//do
	//encodedByte = X MOD 128
	//X = X DIV 128
	//// if there are more data to encode, set the top bit of this byte
	//if ( X > 0 )
	//encodedByte = encodedByte OR 128
	//endif
	//'output' encodedByte
	//while ( X > 0 )
	//
	//Where MOD is the modulo operator (% in C), DIV is integer division (/ in C), and OR is bit-wise or (| in C).
	for {
		//encodedByte = X MOD 128
		//MOD 128 equals AND 127
		var b = uint8(len & 0x7F)
		len = len >> 7

		if len > 0 {
			b = b | 0x80
		}

		buf = append(buf, b)

		if len <= 0 {
			return
		}
	}
}
