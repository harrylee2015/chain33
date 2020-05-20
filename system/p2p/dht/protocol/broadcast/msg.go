package broadcast

import (
	"encoding/hex"
	"github.com/33cn/chain33/types"
)

func (protocol *broadCastProtocol) recvConMsg(conMsg *types.P2PConMsg, pid, peerAddr string) error {
	msg := conMsg.GetConMsg()
	hash := hex.EncodeToString(types.Encode(msg))
	if msg == nil {
		return types.ErrInvalidParam
	}
	//嗅探包要避免重复接收
	if msg.ToPeerID != "" {
		//单点广播不过滤
		if err := protocol.postConsensus("", msg); err != nil {
			log.Error("recvConMsg", "send ConMsg to consensus Error", err.Error())
			return errSendMsgConsensus
		}
		return nil
	}
	if protocol.conMsgFilter.AddWithCheckAtomic(hash, true) {
		return nil
	}
	if err := protocol.postConsensus("", msg); err != nil {
		log.Error("recvConMsg", "send ConMsg to consensus Error", err.Error())
		return errSendMsgConsensus
	}
	//默认最大广播次数为3，超过就丢弃掉
	if route := conMsg.GetRoute(); route.TTL < defaultLtTxBroadCastTTL {
		//将节点id添加到发送过滤, 避免冗余发送
		addIgnoreSendPeerAtomic(protocol.conMsgSendFilter, hash, pid)
		route.TTL++
		protocol.broadcast(conMsg)
	}
	return nil
}
