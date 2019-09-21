package utils

// CheckErr - panics error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
