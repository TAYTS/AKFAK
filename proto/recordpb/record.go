package recordpb

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitialiseEmptyRecord return a Record pointer type with default value
func (*Record) InitialiseEmptyRecord() *Record {
	return &Record{}
}

// InitialiseRecordWithMsg return a Record pointer type with message
func (*Record) InitialiseRecordWithMsg(message string) *Record {
	msgLen, msg := convertStringToBytes(message)
	return &Record{
		ValueLen: int32(msgLen),
		Value:    msg,
	}
}

// UpdateMessage update the message inside the Record
func (rcd *Record) UpdateMessage(message string) {
	msgLen, msg := convertStringToBytes(message)

	rcd.ValueLen = int32(msgLen)
	rcd.Value = msg
}

func convertStringToBytes(message string) (int, []byte) {
	msgByte := []byte(message)
	msgLen := len(msgByte)

	return msgLen, msgByte
}
