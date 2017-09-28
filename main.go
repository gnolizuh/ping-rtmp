package main

import (
	"fmt"
	"os"
	"time"
	"github.com/zhangpeihao/goflv"
	"github.com/zhangpeihao/gortmp"
	"encoding/binary"
	"github.com/urfave/cli"
)

func NewPublish(url string) error {
	h := new(OutBoundHandler)
	h.role = 0
	h.name = "ping-rtmp"

	var err error
	h.c, err = gortmp.Dial(url, h, 100)
	if err != nil {
		return err
	}

	err = h.c.Connect()
	if err != nil {
		return err
	}

	return nil
}

func NewPlay(url string) error {
	h := new(OutBoundHandler)
	h.role = 1
	h.maxtimes = 100
	h.report = 0.0
	h.nreport = 0
	h.name = "ping-rtmp"

	var err error
	h.c, err = gortmp.Dial(url, h, 100)
	if err != nil {
		return err
	}

	err = h.c.Connect()
	if err != nil {
		return err
	}

	return nil
}

type OutBoundHandler struct {
	role     int
	status   uint
	c        gortmp.OutboundConn
	maxtimes uint
	report   float64
	nreport  int64
	name     string
}

func (h *OutBoundHandler) OnStatus(c gortmp.OutboundConn) {
	if c == nil {
		return
	}

	var err error
	h.status, err = c.Status()
	if err != nil {
		fmt.Println(err)
	}
}

func (h OutBoundHandler) OnClosed(conn gortmp.Conn) {}
func (h OutBoundHandler) OnReceived(conn gortmp.Conn, message *gortmp.Message) {
	if h.role == 1 {
		switch message.Type {
		case gortmp.VIDEO_TYPE:
			now1 := int64(binary.BigEndian.Uint64(message.Buf.Bytes()))
			now := time.Now().UnixNano()
			diff := float64(now - now1) / 1000000 // convert ns to ms
			h.report = h.report + diff
			h.nreport ++
			fmt.Printf("%d bytes from: sid=%d csid=%d time=%.3f ms\n",
				message.Buf.Len(), message.StreamID, message.ChunkStreamID, diff)
		}
	}
}

func (h OutBoundHandler) OnReceivedRtmpCommand(conn gortmp.Conn, command *gortmp.Command) {}

func (h *OutBoundHandler) OnStreamCreated(conn gortmp.OutboundConn, stream gortmp.OutboundStream) {
	sh := new(StreamHandle)
	sh.h = h
	stream.Attach(sh)

	if h.role == 0 {
		err := stream.Publish(h.name, "live")
		if err != nil {
			fmt.Printf("Publish error: %s", err.Error())
			return
		}
	} else {
		err := stream.Play(h.name, nil, nil, nil)
		if err != nil {
			fmt.Printf("Play error: %s", err.Error())
			return
		}
	}
}

type StreamHandle struct {
	h *OutBoundHandler
}

func (h StreamHandle) OnPlayStart(stream gortmp.OutboundStream) {}
func (h StreamHandle) OnReceived(conn gortmp.Conn, message *gortmp.Message) {}

func (h StreamHandle) OnPublishStart(stream gortmp.OutboundStream) {
	// loop send ping-pong packets.
	go func() {
		pts := uint32(0)
		epoch := time.Now().UnixNano()

		for h.h.status == gortmp.OUTBOUND_CONN_STATUS_CREATE_STREAM_OK {
			now := time.Now().UnixNano()
			data := make([]byte, 32)
			binary.BigEndian.PutUint64(data, uint64(now))

			header := flv.TagHeader{
				TagType:   flv.VIDEO_TAG,
				DataSize:  8,
				Timestamp: pts}

			delta := (uint32)((now - epoch) / 1000000)

			if err := stream.PublishData(header.TagType, data, delta); err != nil {
				fmt.Println("PublishData() error:", err)
				break
			}

			time.Sleep(time.Millisecond * time.Duration(1000))
		}
	} ()
}

func main() {
	ping := cli.NewApp()
	ping.Name = "ping-rtmp"
	ping.Usage = "RTMP layer PingPong test."
	ping.Version = "1.0"
	ping.Flags = []cli.Flag {
		cli.StringFlag{
			Name:  "push",
			Value: "rtmp://127.0.0.1/live/",
			Usage: "specify a push url for publisher",
		},
		cli.StringFlag{
			Name:  "pull",
			Value: "rtmp://127.0.0.1/live/",
			Usage: "specify a pull url for player",
		},
	}

	ping.Action = func(c *cli.Context) error {
		config := Config{}
		config.Push = c.String("push")
		config.Pull = c.String("pull")

		if err := NewPublish(config.Push); err != nil {
			fmt.Printf("publish(%s) request rejected due to [%s].",
				config.Push, err.Error())
			return err
		}

		if err := NewPlay(config.Pull); err != nil {
			fmt.Printf("play(%s) request rejected due to [%s].",
				config.Pull, err.Error())
			return err
		}

		fmt.Printf("PING %s <-> %s: 32 data bytes\n",
			config.Push, config.Pull)

		finish := make(chan int)
		for {
			select {
			case <-finish:
				break
			}
		}

		return nil
	}

	ping.Run(os.Args)
}
