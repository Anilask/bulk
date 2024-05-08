package validator

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"

	"bulk/constant"
	"bulk/errors"
)

const tagName = "validate"

type Validator interface {
	Validate(interface{}) (bool, errors.ErrorItem)
}

type DefaultValidator struct{}

func (v DefaultValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	return true, err
}

type TransactionIDValidator struct{}

func (v TransactionIDValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	l := len(val.(string))
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)
	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "Transaction ID parameter is invalid"
		// Set the error code and details as needed
		return false, err
	}
	return true, err
}

type ReferenceNumberValidator struct{}

func (v ReferenceNumberValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	l := len(val.(string))
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)
	if l == 0 || !regex.MatchString(val.(string)) {
		// Replace the following line with the correct implementation based on your codebase
		// errs := ew.NewErrorFormatter(ew.ReferenceNumberIsInvalid0311)
		// err = errors.GetErrorItem(errs)
		err.Message = "Reference number is invalid" // Example error message
		// Set the error code and details as needed
		return false, err
	}
	return true, err
}

type MerchantValidator struct{}

func (v MerchantValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regexMerchantCode := regexp.MustCompile(`^[A-Z0-9]{6}$`)
	codeUpperCase := strings.ToUpper(val.(string))
	if !regexMerchantCode.MatchString(codeUpperCase) {
		// Replace the following line with the correct implementation based on your codebase
		// errs := ew.NewErrorFormatter(ew.MerchantCodeMustHaveSixCharacters0308)
		// err = errors.GetErrorItem(errs)
		err.Message = "Merchant code must have six characters" // Example error message
		// Set the error code and details as needed
		return false, err
	}
	return true, err
}

type CustomerValidator struct{}

func (v CustomerValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regexCustomerId := regexp.MustCompile(`^[A-Z0-9]{6}-[a-zA-Z0-9]+$`)
	if !regexCustomerId.MatchString(val.(string)) {
		// Replace the following line with the correct implementation based on your codebase
		// errs := ew.NewErrorFormatter(ew.CustomerIdIsInvalid0305)
		// err = errors.GetErrorItem(errs)
		err.Message = "Customer ID is invalid" // Example error message
		// Set the error code and details as needed
		return false, err
	}
	return true, err
}

type PhoneValidator struct{}

func (v PhoneValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}

	// Min 8 digits without 628 | 08, and Max 11 digits without 628 | 08.
	// Ex: (62)812345678 | (0)812345678 -> It's the same.
	// It refers to kominfo rules.
	re := regexp.MustCompile(`(^62[0-9]{9,14}$)`)

	if !re.MatchString(val.(string)) {
		err.Message = "Phone number is invalid"
		return false, err
	}

	return true, err
}

type CurrencyValidator struct{}

func (v CurrencyValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	if !strings.EqualFold(val.(string), constant.IndonesiaRupiah) {
		err.Message = "Currency is not supported"
		return false, err
	}
	return true, err
}

type BankValidator struct{}

func (v BankValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	if len(val.(string)) != 0 {
		rBankCode := regexp.MustCompile(`^[0-9]{3}$`)
		if !rBankCode.MatchString(val.(string)) {
			err.Message = "Bank code is invalid"
			return false, err
		}
	}
	return true, err
}

type AmountValidator struct{}

func (v AmountValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^(0|[1-9]\d*)(\.\d+)?$`)
	if !regex.MatchString(val.(string)) {
		err.Message = "Amount is invalid"
		return false, err
	}

	return true, err
}

type CorrelationIdValidator struct{}

func (v CorrelationIdValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)

	l := len(val.(string))

	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "Correlation ID is missing or invalid"
		return false, err
	}
	return true, err
}

type ACorrelationIdValidator struct{}

func (v ACorrelationIdValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)

	l := len(val.(string))

	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "ACorrelation ID is missing or invalid"
		return false, err
	}
	return true, err
}

type RpiCodeValidator struct{}

func (v RpiCodeValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	l := len(val.(string))
	regex := regexp.MustCompile(`^(RPI_)[a-zA-Z0-9]{10,12}$`)
	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "RP Code is invalid"
		return false, err
	}
	return true, err
}

type AyoconnectTokenValidator struct{}

func (v AyoconnectTokenValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)
	if !regex.MatchString(val.(string)) {
		err.Message = "Invalid Ayoconnect token"
		return false, err
	}
	return true, err
}

type RecurringIntervalValidator struct {
	Min int
	Max int
}

func (v RecurringIntervalValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	if valInt, ok := val.(int); ok {
		if valInt < v.Min || valInt > v.Max {
			err.Message = "Invalid recurring interval"
			return false, err
		}
	} else {
		err.Message = "Invalid recurring interval format"
		return false, err
	}
	return true, err
}

type NumberOfOccurencesValidator struct {
	Min int
	Max int
}

func (v NumberOfOccurencesValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	if valInt, ok := val.(int); ok {
		if valInt < v.Min || valInt > v.Max {
			err.Message = "Invalid number of occurrences"
			return false, err
		}
	} else {
		err.Message = "Invalid number of occurrences format"
		return false, err
	}
	return true, err
}

type EmailIdValidator struct{}

func (v EmailIdValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(val.(string)) {
		err.Message = "Invalid email address"
		return false, err
	}

	return true, err
}

type OtpValidator struct{}

func (v OtpValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regexOtp := regexp.MustCompile(`(^[0-9]{6}|^[0-9]{4})+$`)
	if !regexOtp.MatchString(val.(string)) {
		err.Message = "Invalid OTP"
		return false, err
	}

	return true, err
}

type MaskedCardValidator struct{}

func (v MaskedCardValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	rCard := regexp.MustCompile("(^[0-9]{4}|^[0-9]{16}|^[*]{12}[0-9]{4})+$")
	if !rCard.MatchString(val.(string)) {
		err.Message = "Invalid card number format"
		return false, err
	}
	return true, err
}

type IpAddressValidator struct{}

func (v IpAddressValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	ip := net.ParseIP(val.(string))
	if ip == nil {
		err.Message = "Invalid IP address format"
		return false, err
	}
	return true, err
}

type BeneficiaryIDValidator struct{}

func (v BeneficiaryIDValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	l := len(val.(string))
	regex := regexp.MustCompile(`^BE_[a-zA-Z0-9]{10}$`)
	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "Invalid beneficiary ID format"
		return false, err
	}
	return true, err
}

type UserIdValidator struct{}

func (v UserIdValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	l := len(val.(string))
	regex := regexp.MustCompile(`^user_[a-zA-Z0-9]{10}$`)
	if l == 0 || !regex.MatchString(val.(string)) {
		err.Message = "Invalid user ID format"
		return false, err
	}
	return true, err
}

type SwiftCodeValidator struct{}

func (v SwiftCodeValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`)

	if !regex.MatchString(val.(string)) {
		err.Message = "Invalid SWIFT code format"
		return false, err
	}

	return true, err
}

type LatitudeValidator struct{}

func (v LatitudeValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`)

	if !regex.MatchString(val.(string)) {
		err.Message = "Invalid latitude format"
		return false, err
	}

	return true, err
}

type LongitudeValidator struct{}

func (v LongitudeValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}
	regex := regexp.MustCompile(`\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)

	if !regex.MatchString(val.(string)) {
		err.Message = "Invalid longitude format"
		return false, err
	}

	return true, err
}

type BeneficiaryAccountNumberValidator struct{}

func (v BeneficiaryAccountNumberValidator) Validate(val interface{}) (bool, errors.ErrorItem) {
	err := errors.ErrorItem{}

	lenAcc := len(val.(string))

	if lenAcc < 8 || lenAcc > 18 {
		err.Message = "Invalid beneficiary account number format"
		return false, err
	}

	return true, err
}

func FindTag(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {

	case constant.CtxTransactionId:
		return TransactionIDValidator{}
	case constant.CtxReferenceNumber:
		return ReferenceNumberValidator{}
	case "customerId":
		return CustomerValidator{}
	case "merchantCode":
		return MerchantValidator{}
	case "phone":
		return PhoneValidator{}
	case "rpiCode":
		return RpiCodeValidator{}
	case "currency":
		return CurrencyValidator{}
	case "amount":
		return AmountValidator{}
	case "bankCode":
		return BankValidator{}
	case "correlationId":
		return CorrelationIdValidator{}
	case "acorrelationId":
		return ACorrelationIdValidator{}
	case "ayoconnectToken":
		return AyoconnectTokenValidator{}
	case "recurringInterval":
		validator := RecurringIntervalValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "numberOfOccurences":
		validator := NumberOfOccurencesValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		return EmailIdValidator{}
	case "otp":
		return OtpValidator{}
	case "maskedCard":
		return MaskedCardValidator{}
	case "ipAddress":
		return IpAddressValidator{}
	case "beneficiaryId":
		return BeneficiaryIDValidator{}
	case "userId":
		return UserIdValidator{}
	case "swiftCode":
		return SwiftCodeValidator{}
	case "beneficiaryAccountNumber":
		return BeneficiaryAccountNumberValidator{}
	case "latitude":
		return LongitudeValidator{}
	case "longitude":
		return LongitudeValidator{}
	default:
		return DefaultValidator{}
	}
}

func ValidateRequest(s interface{}) *[]errors.ErrorItem {
	errorItems := []errors.ErrorItem{}
	v := reflect.ValueOf(s)

	for i := 0; i < reflect.Indirect(v).NumField(); i++ {
		tag := reflect.Indirect(v).Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		validator := FindTag(tag)
		err := errors.ErrorItem{}
		valid, err := validator.Validate(reflect.Indirect(v).Field(i).Interface())

		if !valid && (errors.ErrorItem{}) != err {
			errorItems = append(errorItems, err)
		}
	}

	if len(errorItems) > 0 {
		return &errorItems
	}

	return nil
}
