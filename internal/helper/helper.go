package helper

import (
	"fmt"
)

var (
	unityList = []string{"", "k", "m", "g"}
)

/*
Percentage returns a percentage
*/
func Percentage(value uint64, limit uint64) float64 {
	if limit <= 0 {
		return 0
	}
	//fmt.Printf("%d/%f\n", value, float64(value))
	return (float64(value) / float64(limit)) * float64(100)
}

/*
FormatMemory formats memory values
*/
func FormatMemory(value uint64) string {
	return unityChooser(float64(value), 0)
}

func unityChooser(value float64, unity int) string {
	if (value > float64(1024)) && (unity <= len(unityList)) {
		return unityChooser(value/1024.0, unity+1)
	}
	return fmt.Sprintf("%01.2f%s", value, unityList[unity])
}
