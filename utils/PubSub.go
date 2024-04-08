package utils

type DeliverMessage struct {
	Channel chan string
	Message chan []byte
}

//
//func (channel *DeliverMessage) TransferInformation() {
//	for {
//		select {
//		case
//		case msg := <-channel.Message:
//			{
//				socketConn.Mutex.Lock()
//				for Uid, cid := range socketConn.Active {
//					if ch == cid {
//						err := socketConn.Conn[Uid].WriteMessage(websocket.TextMessage, msg)
//						if err != nil {
//							_ = socketConn.Conn[Uid].Close()
//							delete(socketConn.Active, Uid)
//							delete(socketConn.Conn, Uid)
//							return
//						}
//					}
//				}
//				socketConn.Mutex.Unlock()
//			}
//		}
//	}
//}
