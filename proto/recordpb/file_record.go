package recordpb

import (
	"AKFAK/utils"
	"bufio"
	"errors"
	fmt "fmt"
	"log"
	"os"
	fpath "path/filepath"
	"strconv"

	"github.com/golang/protobuf/proto"
)

// FileRecord is the used for handling/interact with the RecordBatch stored in the local file
type FileRecord struct {
	dir                string
	filename           string
	file               *os.File
	currentOffset      int64
	lastOffsetFHandler *os.File
}

const maxBaseOffsetSize = 10 // int64: 8 bytes, 2 bytes for proto message key
const maxBatchLengthSize = 6 // int32: 4 bytes, 2 bytes for proto message key
const maxPartialHeaderPeekSize = maxBaseOffsetSize + maxBatchLengthSize

///////////////////////////////////
// 		     Public Methods		     //
///////////////////////////////////

// InitialiseFileRecordFromFilepath create the FileRecord that used to handle the log read/write from/to file
func InitialiseFileRecordFromFilepath(filepath string) (*FileRecord, error) {
	// Open the file based on the filename given
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Unable to create or open the file: %v\n", err)
		return nil, err
	}
	dir, filename := fpath.Split(filepath)

	lastOffsetFHandler, err := os.OpenFile(fmt.Sprintf("%voffset", dir), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Unable to create or open the file: %v\n", err)
		return nil, err
	}

	// Create an instance of FileRecord
	fileRecord := &FileRecord{
		dir:                dir,
		filename:           filename,
		file:               file,
		currentOffset:      0,
		lastOffsetFHandler: lastOffsetFHandler,
	}

	return fileRecord, nil
}

// SetOffset allows the filerecord offset to be set 
func (fileRcd *FileRecord) SetOffset(offset int64) {
	fileRcd.currentOffset = offset
}

func (fileRcd *FileRecord) GetOffset() int64 {
	return fileRcd.currentOffset
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

// WriteToFileByBaseOffset write RecordBatch to the file based on the BaseOffset attribute and return RecordBatch with updated BaseOffset if the BaseOffset is -1 else unchange
// return nil if there is error
func (fileRcd *FileRecord) WriteToFileByBaseOffset(rcdBatch *RecordBatch) (*RecordBatch, error) {
	// dereference the input pointer as later it will change the BaseOffset of the RecordBatch
	localRcdBatch := *rcdBatch
	// Convert proto message to []byte
	if localRcdBatch.GetBaseOffset() == -1 {
		localRcdBatch.UpdateBaseOffset(fileRcd.getFileSize())
	}
	recordBatchByte, err := proto.Marshal(&localRcdBatch)
	if err != nil {
		log.Fatalf("Unable to marshall the Record message: %v\n", err)
		return nil, err
	}

	// Write or append the new RecordBatch proto message
	if _, err := fileRcd.file.WriteAt(recordBatchByte, localRcdBatch.GetBaseOffset()); err != nil {
		log.Fatalf("Unable to write to file: %v\n", err)
		return nil, err
	}
	fileRcd.file.Sync()

	// update the last log offset
	fileRcd.lastOffsetFHandler.WriteAt([]byte(fmt.Sprintf("%v", localRcdBatch.GetBaseOffset())), 0)
	fileRcd.lastOffsetFHandler.Sync()

	// return the new file size
	return &localRcdBatch, nil
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

// GetLastEndOffset get the last committed log offset
func (fileRcd *FileRecord) GetLastEndOffset() int64 {
	reader := bufio.NewScanner(fileRcd.lastOffsetFHandler)
	reader.Scan()
	offset, _ := strconv.ParseInt(reader.Text(), 10, 64)
	return offset
}

// ShiftReadOffset is used to shift the current file reading head/pointer
// ** Please use with care and make sure the offset provided is correctly
// match the record batch starting position **
func (fileRcd *FileRecord) ShiftReadOffset(offset int64) {
	fileRcd.currentOffset = offset
}

///////////////////////////////////
// 	 	    Private Methods	    	 //
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
	proto.Unmarshal(rcrBatchPartialHeaderData, recordBatch)

	// Get the BatchLength from proto message and convert to ints
	batchLength := utils.BytesToInt(recordBatch.GetBatchLength())

	return batchLength, nil
}

// hasRemaining check if there are any more bytes to read from the reader buffer
func (fileRcd *FileRecord) hasRemaining() bool {
	// Get the current file byte size
	lastLogOffset := fileRcd.GetLastEndOffset()

	// Compare the file size and the current file read position
	if lastLogOffset > fileRcd.currentOffset {
		return true
	}
	return false
}

// getFileSize get the current file size (num of bytes)
func (fileRcd *FileRecord) getFileSize() int64 {
	// Get the current file byte size
	stats, _ := fileRcd.file.Stat()
	return stats.Size()
}
