package rock7

import (
	"errors"
)

const (
	sendURL = "https://secure.rock7mobile.com/rockblock/MT"
)

var (
	ErrLoginCred  = errors.New("invalid login credentials")
	ErrWrongIMEI  = errors.New("no RockBLOCK with this IMEI found on your account")
	ErrNoRental   = errors.New("rockBLOCK has no line rental")
	ErrInsuffCred = errors.New("your account has insufficient credit")
	ErrDecodHex   = errors.New("could not decode hex data")
	ErrLongData   = errors.New("data too long")
	ErrNoData     = errors.New("no data")
	ErrSystem     = errors.New("system Error")
	ErrDefaultSet = errors.New("default IMEI not set")
	ErrUnknownErr = errors.New("invalid error # returned by system")

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
