package rtpstrategy

import (
	"log"

	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pion/rtp"
)

type ProcessRTPPacketCommand struct {
	Packet *rtp.Packet
	Format format.Format
}

type CommandHandler struct {
}

func NewRtpCommandHandler() CommandHandler {
	return CommandHandler{}
}

func (c *CommandHandler) ProcessRTPPacket(cmd ProcessRTPPacketCommand) (*rtp.Packet, error) {
	log.Println("Processing RTP packet", cmd.Packet)
	return cmd.Packet, nil
}
