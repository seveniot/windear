package tools

import (
	"net"
	"bytes"
	"crypto/sha1"
	"io"
	"time"
	"math/rand"
	"strconv"
	"encoding/hex"
	"github.com/sirupsen/logrus"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/4/12
 * 
*/

func GenerateNodeId() (addr string){
	h := sha1.New()
	io.WriteString(h,getMacAddr())
	io.WriteString(h,time.Now().Format("2006010215040599"))
	io.WriteString(h,strconv.FormatFloat(rand.ExpFloat64(),'d',-1,64))

	return hex.EncodeToString(h.Sum(nil))
}

func getMacAddr()(addr string){
	interfaces, err := net.Interfaces()

	if err!=nil{
		logrus.Fatal(err)
	}

	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				addr = i.HardwareAddr.String()
				break
			}
		}
	}

	return
}
