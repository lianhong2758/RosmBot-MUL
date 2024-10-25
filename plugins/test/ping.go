package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() { // 主函数
	en := rosm.Register(&rosm.PluginData{
		Name: "ping",
		Help: "- /ping xxx",
	})
	en.AddRex(`^\/ping\s*(https:\/\/||http:\/\/)?(\S*)`).SetBlock(true).Handle(func(ctx *rosm.Ctx) {
		var (
			timeout      int64 = 1000 // 超时
			size         int   = 64   // 字节
			count        int   = 4    // 次数
			typ          uint8 = 8
			code         uint8 = 0
			sendCount    int
			successCount int
			failCount    int
			minTs        int64 = 1000
			maxTs        int64 = 0
			totalTs      int64
			msg          strings.Builder
		)
		type iCmp struct {
			Type        uint8
			Code        uint8
			CheckSum    uint16
			ID          uint16
			SequenceNum uint16
		}
		dstIp := ctx.Being.Rex[2]
		if ctx.Being.Rex[1] != "" {
			dstIp = dstIp[:len(dstIp)-1]
		}
		conn, err := net.DialTimeout("ip:icmp", dstIp, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		msg.WriteString(fmt.Sprintf("开始Ping %s[%s] 具有%d字节的数据:\n", dstIp, conn.RemoteAddr(), size))
		for i := range count {
			msg.WriteString(fmt.Sprintf("回复%d:", i+1))
			sendCount++
			icmp := &iCmp{
				Type:        typ,
				Code:        code,
				CheckSum:    0,
				ID:          1,
				SequenceNum: 1,
			}
			data := make([]byte, size)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, icmp)
			buffer.Write(data)
			data = buffer.Bytes()
			checkSum := checkSum(data)
			data[2] = byte(checkSum >> 8)
			data[3] = byte(checkSum)

			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
			t1 := time.Now()
			_, err := conn.Write(data)
			if err != nil {
				failCount++
				msg.WriteString("write timeout\n")
				totalTs += 1000
				maxTs = 1000 // 最长
				time.Sleep(time.Second)
				continue
			}
			buf := make([]byte, 65535)
			n, err := conn.Read(buf)
			if err != nil {
				failCount++
				msg.WriteString("read timeout\n")
				totalTs += 1000
				maxTs = 1000 // 最长
				time.Sleep(time.Second)
				continue
			}
			successCount++
			ts := time.Since(t1).Milliseconds()
			if minTs > ts {
				minTs = ts
			}
			if maxTs < ts {
				maxTs = ts
			}
			totalTs += ts
			msg.WriteString(fmt.Sprintf("字节=%d 时间=%dms TTL=%d\n", n-28, ts, buf[8]))
			time.Sleep(time.Second)
		}
		msg.WriteString("统计信息:\n")
		msg.WriteString(fmt.Sprintf("数据包: 已发送=%d 已接收=%d 丢失=%d(%d%%丢失)\n往返行程的估计时间(ms):\n最短=%dms 最长=%dms 平均=%dms",
			sendCount, successCount, failCount, int(float64(failCount)/float64(sendCount)*100), minTs, maxTs, totalTs/int64(sendCount)))
		ctx.Send(message.Text(msg.String()))
	})
}

func checkSum(data []byte) uint16 {
	length := len(data)
	index := 0
	var sum uint32 = 0
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		length -= 2
		index += 2
	}
	if length != 0 {
		sum += uint32(data[index])
	}
	hi16 := sum >> 16
	for hi16 != 0 {
		sum = hi16 + uint32(uint16(sum))
		hi16 = sum >> 16
	}
	return uint16(^sum)
}
