package utils

import "strconv"

//RUintToString converts optional uint(*uint) to optional string(*string)
func RUintToString(u *uint) *string {
	var result *string
	if u != nil {
		tmp := strconv.Itoa(int(*u))
		result = &tmp
	}
	return result
}

//RStringToUint converts optional string(*string) to optional uint(*uint)
func RStringToUint(s *string) *uint {
	var result *uint
	if s != nil {
		i, err := strconv.Atoi(*s)
		if err != nil {
			panic(err)
		}
		ui := uint(i)
		result = &ui
	}
	return result
}
