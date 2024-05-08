package validator

import (
	"reflect"
	"testing"

	"bulk/constant"
	"bulk/errors"
)

func TestValidate(t *testing.T) {
	createTokenStruct := DefaultValidator{}

	typeValidator := DefaultValidator{}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		Want    bool
		WantErr *errors.CustomError
	}{
		{
			name: "Success",
			args: args{
				i: createTokenStruct,
			},
			Want:    true,
			WantErr: &errors.CustomError{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, _ := typeValidator.Validate(tt.args.i)
				if got != tt.Want {
					t.Errorf("Validate() = %v, Want = %v", got, tt.Want)
				}
			},
		)
	}
}

type common struct {
	TransactionId   string `json:"transactionId" validate:"transactionId"`
	ReferenceNumber string `json:"referenceNumber" validate:"referenceNumber"`
	CorrelationId   string `json:"correlationId" validate:"correlationId"`
	AyoconnectToken string `json:"ayoconnectToken" validate:"ayoconnectToken"`
	RpiCode         string `json:"rpiCode" validate:"rpiCode"`
}

type bifast struct {
	IpAddress                string `json:"ipAddress" validate:"ipAddress"`
	BeneficiaryId            string `json:"beneficiaryId" validate:"beneficiaryId"`
	UserId                   string `json:"userId" validate:"userId"`
	SwiftCode                string `json:"bankCode" validate:"swiftCode"`
	BeneficiaryAccountNumber string `json:"accountNumber" validate:"beneficiaryAccountNumber"`
}

type customerDetails struct {
	CustomerId   string `json:"customerId" validate:"customerId"`
	Phone        string `json:"phone" validate:"phone"`
	MerchantCode string `json:"merchantCode" validate:"merchantCode"`
	Email        string `json:"email" validate:"email"`
}

type paymentDetails struct {
	Currency   string `json:"currency" validate:"currency"`
	Amount     string `json:"amount" validate:"amount"`
	BankCode   string `json:"bankCode" validate:"bankCode"`
	Otp        string `json:"otp" validate:"otp"`
	MaskedCard string `json:"maskedCard" validate:"maskedCard"`
}
type testingValidation struct {
	RecurringInterval  int `json:"recurringInterval" validate:"recurringInterval,min=1,max=3"`
	NumberOfOccurences int `json:"numberOfOccurences" validate:"numberOfOccurences,min=1,max=24"`
	Id                 int `json:"id" validate:"-"` // will not be applicable for validation
}

type locationStruct struct {
	Latitude  string `json:"latitude" validate:"latitude,omitempty"`
	Longitude string `json:"longitude" validate:"longitude,omitempty"`
}

func TestFindTag(t *testing.T) {
	typeValidator := DefaultValidator{}

	type args struct {
		req string
	}

	test := []struct {
		name string
		args args
		Want interface{}
	}{
		{
			name: "Default Validator",
			args: args{
				req: "",
			},
			Want: typeValidator,
		},
	}

	for _, tt := range test {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := FindTag(tt.args.req); got != tt.Want {
					t.Errorf("FindTag() = %v, Want = %v", got, tt.Want)
				}
			},
		)
	}
}

func returnCommonValues(validateMap map[string]string) common {
	resp := common{
		TransactionId:   constant.DummyAlphaNumericCorrectLength,
		ReferenceNumber: constant.DummyAlphaNumericCorrectLength,
		CorrelationId:   constant.DummyAlphaNumericCorrectLength,
		AyoconnectToken: constant.DummyAlphaNumericCorrectLength,
		RpiCode:         constant.DummyRPICode14Length,
	}
	for key, value := range validateMap {
		switch key {
		case constant.CtxTransactionId:
			resp.TransactionId = value
		case constant.CtxReferenceNumber:
			resp.ReferenceNumber = value
		case "correlationId":
			resp.CorrelationId = value
		case "ayoconnectToken":
			resp.AyoconnectToken = value
		case "rpiCode":
			resp.RpiCode = value
		}
	}
	return resp
}

func buildBiFastValues(validateMap map[string]string) bifast {
	resp := bifast{
		IpAddress:                constant.DummyIpAddress,
		BeneficiaryId:            constant.DummyBeneficiaryId,
		UserId:                   constant.DummyUserId,
		SwiftCode:                constant.DummySwiftCode,
		BeneficiaryAccountNumber: constant.DummyAccountNumber,
	}
	for key, value := range validateMap {
		switch key {
		case "ipAddress":
			resp.IpAddress = value
		case "beneficiaryId":
			resp.BeneficiaryId = value
		case "userId":
			resp.UserId = value
		case "bankCode":
			resp.SwiftCode = value
		case "accountNumber":
			resp.BeneficiaryAccountNumber = value
		}
	}
	return resp
}

func buildCustomerDetails(validateMap map[string]string) customerDetails {
	resp := customerDetails{
		CustomerId:   constant.DummyCorrectCustomerID,
		Phone:        constant.DummyPhoneNumber,
		MerchantCode: constant.DummyMerchantCode,
		Email:        constant.DummyEmail,
	}
	for key, value := range validateMap {
		switch key {
		case "customerId":
			resp.CustomerId = value
		case "phone":
			resp.Phone = value
		case "merchantCode":
			resp.MerchantCode = value
		case "email":
			resp.Email = value
		}
	}
	return resp
}

func buildPaymentDetails(validateMap map[string]string) paymentDetails {
	resp := paymentDetails{
		Currency:   constant.DummyCurrency,
		Amount:     constant.DummyAmount,
		BankCode:   constant.DummyBankCode,
		Otp:        constant.DummyOtp,
		MaskedCard: constant.DummyMaskedCard,
	}
	for key, value := range validateMap {
		switch key {
		case "currency":
			resp.Currency = value
		case "amount":
			resp.Amount = value
		case "bankCode":
			resp.BankCode = value
		case "otp":
			resp.Otp = value
		case "maskedCard":
			resp.MaskedCard = value
		}
	}
	return resp
}

func buildTestingValidation(validateMap map[string]int) testingValidation {
	resp := testingValidation{
		RecurringInterval:  2,
		NumberOfOccurences: 3,
		Id:                 1,
	}
	for key, value := range validateMap {
		switch key {
		case "recurringInterval":
			resp.RecurringInterval = value
		case "numberOfOccurences":
			resp.NumberOfOccurences = value
		}
	}
	return resp
}

func buildLocationValidation(validateMap map[string]string) locationStruct {
	resp := locationStruct{
		Latitude:  "-6.175110",
		Longitude: "106.865036",
	}
	for key, value := range validateMap {
		switch key {
		case "latitude":
			resp.Latitude = value
		case "longitude":
			resp.Longitude = value
		}
	}
	return resp
}

func TestValidateRequest(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want *[]errors.ErrorItem
	}{
		{
			name: "SuccessCommonValues",
			args: args{
				s: returnCommonValues(map[string]string{}),
			},
			want: nil,
		},
		{
			name: "SuccessCustomerValues",
			args: args{
				s: buildCustomerDetails(map[string]string{}),
			},
			want: nil,
		},
		{
			name: "SuccessPaymentDetails",
			args: args{
				s: buildPaymentDetails(map[string]string{}),
			},
			want: nil,
		},
		{
			name: "SuccessTestingValidation",
			args: args{
				s: buildTestingValidation(map[string]int{}),
			},
			want: nil,
		},
		{
			name: "Fail  recurringInterval and number of occurences 0",
			args: args{
				s: buildTestingValidation(map[string]int{
					"recurringInterval":  0,
					"numberOfOccurences": 0,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0803",
					Message: "error.validator.0803",
					Details: "The 'interval' value is invalid",
				},
				{
					Code:    "0801",
					Message: "error.validator.0801",
					Details: "The number of occurrences is invalid",
				},
			},
		},
		{
			name: "SuccessLocationValidation",
			args: args{
				s: buildLocationValidation(map[string]string{}),
			},
			want: nil,
		},
		{
			name: "Fail  location validation alpha characters",
			args: args{
				s: buildLocationValidation(map[string]string{
					"latitude":  "a",
					"longitude": "c",
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0916",
					Message: "error.validator.0916",
					Details: "Request has an incorrect location information",
				},
				{
					Code:    "0916",
					Message: "error.validator.0916",
					Details: "Request has an incorrect location information",
				},
			},
		},
		{
			name: "Fail  location validation special characters",
			args: args{
				s: buildLocationValidation(map[string]string{
					"latitude":  "@",
					"longitude": "#",
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0916",
					Message: "error.validator.0916",
					Details: "Request has an incorrect location information",
				},
				{
					Code:    "0916",
					Message: "error.validator.0916",
					Details: "Request has an incorrect location information",
				},
			},
		},
		{
			name: "Special characters in currency amount bank code otp masked card",
			args: args{
				s: buildPaymentDetails(map[string]string{
					"currency":   constant.DummySpecialCharCurrency,
					"amount":     constant.DummySpecialCharAmount,
					"bankCode":   constant.DummySpecialCharBankCode,
					"otp":        constant.DummySpecialCharacterOtp,
					"maskedCard": constant.DummySpecialCharMaskedCard,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0406",
					Message: "error.validator.0406",
					Details: "Currency is not supported",
				},
				{
					Code:    "0319",
					Message: "error.validator.0319",
					Details: "The 'amount' parameter is invalid",
				},
				{
					Code:    "0016",
					Message: "error.validator.0016",
					Details: "Bank code is invalid. Please check with Ayoconnect Team.",
				},
				{
					Code:    "0011",
					Message: "error.validator.0011",
					Details: "The (OTP) passcode is invalid",
				},
				{
					Code:    "0107",
					Message: "error.validator.0107",
					Details: "Card number is invalid. Please enter a valid card number.",
				},
			},
		},
		{
			name: "Invalid amount bank code otp masked card",
			args: args{
				s: buildPaymentDetails(map[string]string{
					"currency":   constant.DummyInvalidCurrency,
					"amount":     constant.DummyInvalidAmount,
					"bankCode":   constant.DummyInvalidBankCode,
					"otp":        constant.DummyInvalidOtp,
					"maskedCard": constant.DummyInvalidMaskedCard,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0406",
					Message: "error.validator.0406",
					Details: "Currency is not supported",
				},
				{
					Code:    "0319",
					Message: "error.validator.0319",
					Details: "The 'amount' parameter is invalid",
				},
				{
					Code:    "0016",
					Message: "error.validator.0016",
					Details: "Bank code is invalid. Please check with Ayoconnect Team.",
				},
				{
					Code:    "0011",
					Message: "error.validator.0011",
					Details: "The (OTP) passcode is invalid",
				},
				{
					Code:    "0107",
					Message: "error.validator.0107",
					Details: "Card number is invalid. Please enter a valid card number.",
				},
			},
		},
		{
			name: "Fail less length of phone number",
			args: args{
				s: buildCustomerDetails(map[string]string{
					"phone": constant.DummyInvalidPhoneNumberWithLessLength,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0601",
					Message: "error.validator.0601",
					Details: "Phone number is invalid",
				},
			},
		},
		{
			name: "Fail wrong  start of phone number and empty email",
			args: args{
				s: buildCustomerDetails(map[string]string{
					"phone": constant.DummyInvalidPhoneNumberWithWrongStartingDigit,
					"email": "",
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0601",
					Message: "error.validator.0601",
					Details: "Phone number is invalid",
				},
				{
					Code:    "0309",
					Message: "error.validator.0309",
					Details: "Email-address is invalid",
				},
			},
		},
		{
			name: "Fail special characters in customer details",
			args: args{
				s: buildCustomerDetails(map[string]string{
					"customerId":   constant.DummyCustomerIDSpecialChar,
					"phone":        constant.DummyPhoneNumberWithSpecialChar,
					"merchantCode": constant.DummyMerchantCodeWithSpecialCharacters,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0305",
					Message: "error.validator.0305",
					Details: "The 'CustomerId' parameter is invalid",
				},
				{
					Code:    "0601",
					Message: "error.validator.0601",
					Details: "Phone number is invalid",
				},
				{
					Code:    "0308",
					Message: "error.validator.0308",
					Details: "Merchant code must have six alphanumeric characters",
				},
			},
		},
		{
			name: "Fail Greater length  in customer details",
			args: args{
				s: buildCustomerDetails(map[string]string{
					"customerId":   constant.DummyInCorrectCustomerID,
					"phone":        constant.DummyPhoneNumberWithGreaterLength,
					"merchantCode": constant.DummyInCorrectMerchantCodeWithGreaterLength,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0305",
					Message: "error.validator.0305",
					Details: "The 'CustomerId' parameter is invalid",
				},
				{
					Code:    "0601",
					Message: "error.validator.0601",
					Details: "Phone number is invalid",
				},
				{
					Code:    "0308",
					Message: "error.validator.0308",
					Details: "Merchant code must have six alphanumeric characters",
				},
			},
		},
		{
			name: "Fail Greater length and special characters",
			args: args{
				s: returnCommonValues(map[string]string{
					"referenceNumber": constant.DummyAlphaNumericGreaterLength,
					"correlationId":   constant.DummyAlphaNumericSpecialCharacter,
					"rpiCode":         constant.DummyAlphaNumericGreaterLength,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0311",
					Message: "error.validator.0311",
					Details: "The 'referenceNumber' parameter is invalid",
				},
				{
					Code:    "0313",
					Message: "error.validator.0313",
					Details: "The 'X-Correlation-ID' header is either missing or invalid",
				},

				{
					Code:    "0804",
					Message: "error.validator.0804",
					Details: "The recurring payment code is invalid",
				},
			},
		},
		{
			name: "Fail Lesser length and empty characters",
			args: args{
				s: returnCommonValues(map[string]string{
					"correlationId":   constant.DummyAlphaNumericLessLength,
					"transactionId":   constant.DummyAlphaNumericSpecialCharacter,
					"ayoconnectToken": "",
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0310",
					Message: "error.validator.0310",
					Details: "The 'transactionId' parameter is invalid",
				},
				{
					Code:    "0313",
					Message: "error.validator.0313",
					Details: "The 'X-Correlation-ID' header is either missing or invalid",
				},
				{
					Code:    "0315",
					Message: "error.validator.0315",
					Details: "The 'ayoconnectToken' parameter is invalid",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRequest(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateBiFastRequest(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want *[]errors.ErrorItem
	}{
		{
			name: "SuccessBiFastValues",
			args: args{
				s: buildBiFastValues(map[string]string{}),
			},
			want: nil,
		},
		{
			name: "Fail Bi fast special characters",
			args: args{
				s: buildBiFastValues(map[string]string{
					"ipAddress":     constant.DummyInvalidIpAddressSpecialCharacters,
					"beneficiaryId": constant.DummyBeneficiaryIdSpecialCharacters,
					"userId":        constant.DummyUserIdSpecialChar,
					"bankCode":      constant.DummyInvalidSwiftCode,
					"accountNumber": "",
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0902",
					Message: "error.validator.0902",
					Details: "The IP address format is invalid",
				},
				{
					Code:    "0903",
					Message: "error.validator.0903",
					Details: "The Beneficiary ID is invalid",
				},
				{
					Code:    "0304",
					Message: "error.validator.0304",
					Details: "Request has invalid public user id format",
				},
				{
					Code:    "0016",
					Message: "error.validator.0016",
					Details: "Bank code is invalid. Please check with Ayoconnect Team.",
				},
				{
					Code:    "0901",
					Message: "error.validator.0901",
					Details: "The Beneficiary account-number is invalid",
				},
			},
		},
		{
			name: "Fail Ip Address length",
			args: args{
				s: buildBiFastValues(map[string]string{
					"ipAddress":     constant.DummyInvalidLengthIpAddress,
					"beneficiaryId": constant.DummyBeneficiaryIdLength,
					"userId":        constant.DummyInvalidLengthUserId,
				}),
			},
			want: &[]errors.ErrorItem{
				{
					Code:    "0902",
					Message: "error.validator.0902",
					Details: "The IP address format is invalid",
				},
				{
					Code:    "0903",
					Message: "error.validator.0903",
					Details: "The Beneficiary ID is invalid",
				},
				{
					Code:    "0304",
					Message: "error.validator.0304",
					Details: "Request has invalid public user id format",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRequest(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
