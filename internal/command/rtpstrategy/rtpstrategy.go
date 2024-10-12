package rtpstrategy

type ProcessRTPPacketCommand struct{}

type CommandHandler struct {
}

func NewRtpCommandHandler() CommandHandler {
	return CommandHandler{}
}

func (c *CommandHandler) ProcessRTPPacket(command ProcessRTPPacketCommand) error {
	return nil
}
