package rock7

import (
	"errors"
)

const (
	sendURL = "https://secure.rock7mobile.com/rockblock/MT"
)

var (
	ErrLoginCred  = errors.New("Invalid login credentials")
	ErrWrongIMEI  = errors.New("No RockBLOCK with this IMEI found on your account")
	ErrNoRental   = errors.New("RockBLOCK has no line rental")
	ErrInsuffCred = errors.New("Your account has insufficient credit")
	ErrDecodHex   = errors.New("Could not decode hex data")
	ErrLongData   = errors.New("Data too long")
	ErrNoData     = errors.New("No data")
	ErrSystem     = errors.New("System Error")
	ErrDefaultSet = errors.New("Default IMEI not set")
	ErrUnknownErr = errors.New("Invalid error # returned by system")

	MappedErrNum = map[int]error{
		10: ErrLoginCred,
		11: ErrWrongIMEI,
		12: ErrNoRental,
		13: ErrInsuffCred,
		14: ErrDecodHex,
		15: ErrLongData,
		16: ErrNoData,
		99: ErrSystem,
	}
)
