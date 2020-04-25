package cluster

// UniqueInsertInt32 insert int32 into the slice uniquely
func uniqueInsertInt32(inpVal int32, inputSlice []int32) []int32 {
	for _, val := range inputSlice {
		if val == inpVal {
			return inputSlice
		}
	}
	inputSlice = append(inputSlice, inpVal)
	return inputSlice
}
