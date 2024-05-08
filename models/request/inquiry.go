package request

type (
	Inquiry struct {
		TransactionId   string     `json:"transactionId"`
		ReferenceNumber string     `json:"referenceNumber"`
		CorrelationId   string     `json:"correlationId"`
		MerchantCode    string     `json:"merchantCode"`
		BulkId          int64      `json:"bulkId"`
		Data            []BulkData `json:"data"`
	}
	BulkData struct {
		BeneficiaryAccountNumber string  `json:"beneficiaryAccountNumber"`
		BeneficiaryBankCode      string  `json:"beneficiaryBankCode"`
		MobileNumber             string  `json:"mobileNumber"`
		Amount                   float64 `json:"amount"`
		PaymentInfo              string  `json:"paymentInfo"`
		InquiryCorrelationId     string  `json:"inquiryCorrelationId"`
	}
	Dibursement struct {
		TransactionId   string                    `json:"transactionId"`
		ReferenceNumber string                    `json:"referenceNumber"`
		CorrelationId   string                    `json:"correlationId"`
		MerchantCode    string                    `json:"merchantCode"`
		MerchantUserId  string                    `json:"merchantCodeId"`
		BulkId          string                    `json:"bulkdId"`
		Data            []BulkDisbursementDetails `json:"data"`
	}
	UpdateDisbursement struct {
		TransactionId               string  `json:"transactionId"`
		ReferenceNumber             string  `json:"referenceNumber"`
		CustomerId                  *string `json:"customerId,omitempty"`
		BeneficiaryCorrelationId    *string `json:"beneficiaryCorrelationId,omitempty"`
		BeneficiaryId               *string `json:"beneficiaryId,omitempty"`
		BeneficiaryStatus           *int32  `json:"beneficiaryStatus,omitempty"`
		DisbursementReferenceNumber *string `json:"disbursementReferenceNumber,omitempty"`
		DisbursementStatus          *int32  `json:"disbursementStatus,omitempty"`
		Status                      *int32  `json:"status,omitempty"`
		FailedReason                *string `json:"failedReason,omitempty"`
		BulkId                      *int64  `json:"bulkId,omitempty"`
		Id                          *int64  `json:"id,omitempty"`
		BenficiaryName              *string `json:"beneficiaryName"`
		BeneficiaryBankName         *string `json:"beneficiaryBankName"`
		Type                        string  `json:"type"`
	}
	BulkDisbursementDetails struct {
		Id                       int64   `json:"id"`
		BulkId                   string  `json:"bulkId"`
		AccountNumber            string  `json:"accountNumber"`
		BankCode                 string  `json:"bankCode"`
		PhoneNumber              string  `json:"phoneNumber"`
		Amount                   float64 `json:"amount"`
		CustomerId               string  `json:"customerId"`
		BeneficiaryCorrelationId string  `json:"beneficiaryCorrelationId"`
		BenficiaryId             string  `json:"beneficiaryId"`
		BeneficiaryStatus        int32   `json:"beneficiaryStatus"`
		DisbursementReferenceNo  string  `json:"disbursementReferenceNumber"`
		DisbursementStatus       int32   `json:"disbursementStatus"`
		Status                   int32   `json:"status"`
		FailedReason             string  `json:"failedReason"`
		PaymentInfo              string  `json:"paymentInfo"`
		BenficiaryName           string  `json:"beneficiaryName"`
		BeneficiaryBankName      string  `json:"beneficiaryBankName"`
	}
)
