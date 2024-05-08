package utils

import (
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"bulk/constant"
	"bulk/logger"
	"github.com/google/uuid"
)

const (
	CSVMimeType  = "text/csv"
	XLSMimeType  = "application/vnd.ms-excel"
	XLSXMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

var validStatusForBulkStatusLogs = map[int]bool{
	constant.BulkStatusFailed:            true,
	constant.BulkStatusRejected:          true,
	constant.BulkStatusTransferInitiated: true,
}

func ConvertStringToINT(feeType string) (int, error) {
	resultFeeType, errConvertFeeType := strconv.Atoi(feeType)
	if errConvertFeeType != nil {
		return 0, errConvertFeeType
	}
	return resultFeeType, nil
}

func GetTime() time.Time {
	location, _ := time.LoadLocation(constant.AsiaJakartaTimeZone)

	return time.Now().In(location)
}

func Add30Day() time.Time {
	return GetTime().AddDate(0, 0, 30)
}

func GenerateNewUniqueNumber() string {
	uuidWithHyphen := uuid.New()
	randomUUID := strings.ReplaceAll(uuidWithHyphen.String(), "-", "")

	return randomUUID
}

func GetGMT7DateTime(format string) string {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now = now.In(loc)
	return now.Format(format)
}

func PanicRecover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Errorf("stacktrace from panic: \n" + string(debug.Stack()))
		log.Errorf("panic: %v", r)
	}
}

func ValidateStatusForBulkStatusLogs(status int) bool {
	return validStatusForBulkStatusLogs[status]
}
