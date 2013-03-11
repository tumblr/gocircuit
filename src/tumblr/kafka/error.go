package kafka

import "errors"

var (
	ErrIO           = errors.New("io or network")
	ErrClosed       = errors.New("already closed")
	ErrArg          = errors.New("invalid argument")
	ErrWire         = errors.New("invalid or dangerous data in wire format")
	ErrNotSupported = errors.New("not supported")
	ErrCompression  = errors.New("unknown compression")
	ErrChecksum     = errors.New("checksum")
	ErrNoBrokers    = errors.New("no brokers")
)

type KafkaError error

var (
	KafkaErrUnknown          = KafkaError(errors.New("kafka: unknown error"))
	KafkaErrNoError          = KafkaError(nil)
	KafkaErrOffsetOutOfRange = KafkaError(errors.New("kafka: offset out of range"))
	KafkaErrInvalidMessage   = KafkaError(errors.New("kafka: invalid message"))
	KafkaErrWrongPartition   = KafkaError(errors.New("kafka: wrong partition"))
	KafkaErrInvalidFetchSize = KafkaError(errors.New("kafka: invalid fetch size"))
)

func KafkaErrorCode(err KafkaError) ErrorCode {
	switch err {
	case KafkaErrUnknown:
		return ErrorCodeUnknown
	case KafkaErrNoError:
		return ErrorCodeNoError
	case KafkaErrOffsetOutOfRange:
		return ErrorCodeOffsetOutOfRange
	case KafkaErrInvalidMessage:
		return ErrorCodeInvalidMessage
	case KafkaErrWrongPartition:
		return ErrorCodeWrongPartition
	case KafkaErrInvalidFetchSize:
		return ErrorCodeInvalidFetchSize
	}
	panic("unknown kafka error")
}

func KafkaCodeError(code ErrorCode) KafkaError {
	switch code {
	case ErrorCodeUnknown:
		return KafkaErrUnknown
	case ErrorCodeNoError:
		return KafkaErrNoError
	case ErrorCodeOffsetOutOfRange:
		return KafkaErrOffsetOutOfRange
	case ErrorCodeInvalidMessage:
		return KafkaErrInvalidMessage
	case ErrorCodeWrongPartition:
		return KafkaErrWrongPartition
	case ErrorCodeInvalidFetchSize:
		return KafkaErrInvalidFetchSize
	}
	panic("unknown kafka error code")
}

// ErrorCode represents a Kafka response error code
type ErrorCode int16

const (
	ErrorCodeUnknown ErrorCode = iota - 1
	ErrorCodeNoError
	ErrorCodeOffsetOutOfRange
	ErrorCodeInvalidMessage
	ErrorCodeWrongPartition
	ErrorCodeInvalidFetchSize
)

func isValidErrorCode(e ErrorCode) bool {
	return e >= ErrorCodeUnknown && e <= ErrorCodeInvalidFetchSize
}

func (x ErrorCode) String() string {
	switch x {
	case ErrorCodeUnknown:
		return "unknown"
	case ErrorCodeNoError:
		return "ok"
	case ErrorCodeOffsetOutOfRange:
		return "offset out of range"
	case ErrorCodeInvalidMessage:
		return "invalid message"
	case ErrorCodeWrongPartition:
		return "wrong partition"
	case ErrorCodeInvalidFetchSize:
		return "invalid fetch size"
	}
	return "Error code not implemented"
}
