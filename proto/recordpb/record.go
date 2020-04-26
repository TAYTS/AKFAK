package recordpb

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitialiseEmptyRecord return a Record pointer type with default value
func InitialiseEmptyRecord() *Record {
	return &Record{}
}

// InitialiseRecordWithMsg return a Record pointer type with message
func InitialiseRecordWithMsg(message string) *Record {
	// msgLen, msg := convertStringToBytes(message)
	return &Record{
		Value: message,
	}
}

// UpdateMessage update the message inside the Record
func (rcd *Record) UpdateMessage(message string) {
	rcd.Value = message
}
