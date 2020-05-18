package broadcast

import (
	"github.com/33cn/chain33/types"
)

func (protocol *broadCastProtocol) recvConMsg(conMsg *types.ConsensusMsg, pid, peerAddr string) error {

	if conMsg.GetData() == nil {
		return types.ErrInvalidParam
	}

	if err := protocol.postConsensus("", conMsg); err != nil {
		log.Error("recvConMsg", "send ConMsg to consensus Error", err.Error())
		return errSendMsgConsensus
	}

	return nil
}
