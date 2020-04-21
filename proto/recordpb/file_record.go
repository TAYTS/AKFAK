package recordpb

import (
	"AKFAK/utils"
	"errors"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
)

// FileRecord is the used for handling/interact with the RecordBatch stored in the local file
type FileRecord struct {
	filename      string
	file          *os.File
	currentOffset int64
}

const maxBaseOffsetSize = 10 // int64: 8 bytes, 2 bytes for proto message key
const maxBatchLengthSize = 6 // int32: 4 bytes, 2 bytes for proto message key
const maxPartialHeaderPeekSize = maxBaseOffsetSize + maxBatchLengthSize

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitialiseFileRecordFromFile create the FileRecord that used to handle the log read from file
func InitialiseFileRecordFromFile(filename string) (*FileRecord, error) {
	// Open the file based on the filename given
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Unable to create or open the file: %v\n", err)
		return nil, err
	}

	// Create an instance of FileRecord
	fileRecord := &FileRecord{
		filename:      filename,
		file:          file,
		currentOffset: 0,
	}

	return fileRecord, nil
}

// SetOffset allows the filerecord offset to be set 
func (fileRcd *FileRecord) SetOffset(offset int64) {
	fileRcd.currentOffset = offset
}

// ReadNextRecordBatch retrieve the next RecordBatch in the file buffer else return EOF error
func (fileRcd *FileRecord) ReadNextRecordBatch() (*RecordBatch, error) {
	// Check if the are still remaining bytes
	if !fileRcd.hasRemaining() {
		return nil, errors.New("EOF")
	}

	// Get the next RecordBatch batch length
	recordBatchLength, err := fileRcd.getNextBatchLength()
	if err != nil {
		return nil, err
	}

	// Create byte array buffer of the same size as the next RecordBatch size to store the byte data of RecordBatch retrieve from the fileReader
	recordBatchBuffer := make([]byte, recordBatchLength)

	// Read the next RecordBatch byte data into the byte buffer created earlier
	numOfReadBytes, err := fileRcd.file.ReadAt(recordBatchBuffer, fileRcd.currentOffset)
	if err != nil {
		return nil, err
	}
	fileRcd.currentOffset += int64(numOfReadBytes)

	// Create an instance of RecordBatch to store the parsed proto message
	recordBatch := &RecordBatch{}

	// Parse the RecordBatch byte data into the proto message
	unmarshalErr := proto.Unmarshal(recordBatchBuffer, recordBatch)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return recordBatch, nil
}

// WriteToFile retrieve the next RecordBatch in the file buffer else return nil
func (fileRcd *FileRecord) WriteToFile(rcdBatch *RecordBatch) error {
	// Convert proto message to []byte
	recordBatchByte, err := proto.Marshal(rcdBatch)
	if err != nil {
		log.Fatalf("Unable to marshall the Record message: %v\n", err)
		return err
	}

	// Write or append the new RecordBatch proto message
	if _, err := fileRcd.file.Write(recordBatchByte); err != nil {
		log.Fatalf("Unable to write to file: %v\n", err)
		return err
	}

	return nil
}

// CloseFile close the file
func (fileRcd *FileRecord) CloseFile() error {
	if fileRcd.file != nil {
		if err := fileRcd.file.Close(); err != nil {
			return err
		}
	}

	return nil
}

///////////////////////////////////
// 		   Private Methods		 //
///////////////////////////////////

// getBatchLength return the next RecordBatch BatchLength based on the Reader position without advancing the fileReader
func (fileRcd *FileRecord) getNextBatchLength() (int, error) {
	// Partial Header must consist of BaseOffset and BatchLength and may contains PartialLeaderEpoch
	rcrBatchPartialHeaderData := make([]byte, maxPartialHeaderPeekSize)
	_, err := fileRcd.file.ReadAt(rcrBatchPartialHeaderData, fileRcd.currentOffset)
	if err != nil {
		return 0, err
	}

	// Create RecordBatch instance which used later to store the parsed data
	recordBatch := &RecordBatch{}

	// Parse the bytes data into the proto message
	unmarshalErr := proto.Unmarshal(rcrBatchPartialHeaderData, recordBatch)
	if unmarshalErr != nil {
		return 0, unmarshalErr
	}

	// Get the BatchLength from proto message and convert to ints
	batchLength := utils.BytesToInt(recordBatch.GetBatchLength())

	return batchLength, nil
}

// hasRemaining check if there are any more bytes to read from the reader buffer
func (fileRcd *FileRecord) hasRemaining() bool {
	// Get the current file byte size
	stats, _ := fileRcd.file.Stat()
	currentFileSize := stats.Size()

	// Compare the file size and the current file read position
	if currentFileSize > fileRcd.currentOffset {
		return true
	}
	return false
}
