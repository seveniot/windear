package codec

import (
	"io"
	"github.com/SevenIOT/windear/ex"
)

func DecodeRemainLen(reader io.Reader) (len uint32, err error){
	var (
		multiplier uint32 = 0
		b = make([]byte, 1)
	)

	for {
		var l int

		l,err = reader.Read(b)

		if l==0{
			continue
		}

		if err!=nil{
			return 0,err
		}

		//encodedByte AND 127
		len = len+uint32(b[0]&0x7F)<<multiplier
		multiplier += 7

		//multiplier>128*128*128
		if multiplier>21{
			//err = errors.New(fmt.Sprintf("illegal remain length multiplier:%v",multiplier))
			return 0,ex.IncorrectRemainLengthMultiplier
		}

		if b[0]&0x80 ==0{
			return
		}
	}
}

func EncodeRemainLength(len uint32)(buf []byte){
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
	for{
		//encodedByte = X MOD 128
		//MOD 128 equals AND 127
		var b = uint8(len&0x7F)
		len = len >> 7

		if len>0{
			b = b|0x80
		}

		buf = append(buf,b)

		if len<=0{
			return
		}
	}
}

