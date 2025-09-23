package common

func RemoveFromSlicePosition(slice []string, position int) []string {

	return append(slice[:position], slice[position+1:]...)

}
