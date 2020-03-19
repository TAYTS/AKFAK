package recordpb

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitialiseRecord return a RecordBatch pointer type with default value
func (*Record) InitialiseRecord(baseTimestamp int32, baseOffset int32) *Record {
	return &Record{}
}

// UpdateMessage update the message inside the Record
func (rcd *Record) UpdateMessage(message string) {
	msgByte := []byte(message)
	msgLen := len(msgByte)

	rcd.ValueLen = int32(msgLen)
	rcd.Value = msgByte
}
