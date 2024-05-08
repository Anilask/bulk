package constant

const (
	HeaderContentType                = "Content-Type"
	HeaderAccept                     = "Accept"
	HeaderCorrelationID              = "X-Correlation-ID"
	MIMEApplicationJSON              = "application/json"
	CtxTransactionID                 = "transactionId"
	CtxReferenceNumber               = "referenceNumber"
	CtxMerchantCode                  = "rmerchantCode"
	CtxBulkDisbursmentID             = "bulkDisbursmentID"
	MessageOK                        = "ok"
	PageString                       = "page"
	LimitString                      = "limit"
	StartDateString                  = "startDate"
	EndDateString                    = "endDate"
	BankCode                         = "bankCode"
	FileDownloadFormat               = "02-01-2006 15:04:05"
	DateFormat                       = "2006-01-02"
	FormatTimeYYYYMMDDHHMMSS         = "20060102150405"
	FormatDateYYYYMMDD               = "20060102"
	GcsFileSuffix                    = "_000000000000.csv"
	TimeFormatStringJKT              = "02-01-2006 15:04:05"
	TimeFormatString                 = "02 Jan 2006 03:04:05 PM"
	AsiaJakartaTimeZone              = "Asia/Jakarta"
	ErrWhenHitClient                 = "Error when hit client: %vError when hit bank: %v"
	FormatPrefixTable                = "`%s.%s.%s`"
	FormatDefineQuery                = "%s %s etl_mv"
	FormatDefineQueryForMerchantList = "%s %s mtl_mv"
	LogQuery                         = " [Query]: %v"
	ErrCreateClient                  = " Error creating client: %v"
	ErrReadQuery                     = " Error reading query: %v"
	ErrReadResponse                  = " Error reading response: %v"
	ErrIterateQuery                  = " Error iterating query: %v"
	LogResponse                      = " [Response]: %v"
	ErrExecuteStatementQuery         = "Error while executing statement: %v"
	ErrPrepareStatementQuery         = "Error while preparing statement: %v"
	ErrScanningQuery                 = "Error while scanning query [%v]"
	TimestampFormat                  = "2006-01-02T15:01:05-07:00"
	LocationJakarta                  = "Asia/Jakarta"
	BulkFileListTimeFormat           = "2006-01-02T15:04:05Z"
	MerchantCodeString               = "merchantCode"
	Page                             = "page"
	Count                            = "count"
	BulkDisbursmentID                = "bulkDisbursmentId"
	BulkDisbursmentName              = "bulkDisbursmentName"
	Uploader                         = "uploader"
	Status                           = "status"
	ErrorStucture                    = "%s %s: %w"

	BulkStatusVerifing           = 4
	BulkStatusVerified           = 5
	BulkStatusTransferInitiated  = 6
	BulkStatusTransferInProgress = 7
	BulkStatusCompleted          = 8
	BulkStatusFailed             = 9
	BulkStatusRejected           = 11

	// error database
	ErrInsertStatusLog = "error inserting bulk_status_logs, query: %v, err :%v"

	// query
	InsertQueryForStatusLog = "INSERT INTO bulk_status_logs (bulkId, reason, updatedBy, currentStatus, updatedStatus) VALUES (?, ?, ?, ?, ?)"
)
const (
	// HeaderContentType   string = "Content-Type"
	// HeaderAccept        string = "Accept"
	HeaderContextTrace  string = "X-Cloud-Trace-Context"
	HeaderCorrelationId string = "X-Correlation-ID"
	MIMEApplicationJson string = "application/json"
	CtxCorrelationId    string = "correlationId"
	CtxTransactionId    string = "transactionId"
	// CtxReferenceNumber  string = "referenceNumber"
	// CtxMerchantCode     string = "merchantCode"
	CtxBankCode   string = "bankCode"
	CtxFieldsType string = "fields"
	ErrorCounts   string = "[%v] err %v"
	MessageHTTP   string = "ok"
	EnvDev        string = "development"
	EnvStage      string = "sandbox"
	CtxTaskId     string = "taskId"
	CtxMessage    string = "message"
	CtxQueuePath  string = "queuePath"
	CtxTaskName   string = "taskName"
)

const (
	// Define salt size
	SaltSize = 8
	// Defines number of times to hash
	HashCount = 6
	// Defines length of hash
	HashLength                                    = 16
	DummyAlphaNumericCorrectLength                = "pBFhaEvQ5f3lCq6vnfdRZ798rV4PSmqn"
	DummyAlphaNumericLessLength                   = "pBFhaEvQ5f3lCq6vnfdRZ798rV4PSmq"
	DummyAlphaNumericGreaterLength                = "pBFhaEvQ5f3lCq6vnfdRZ798rV4PSmqna"
	DummyAlphaNumericSpecialCharacter             = "pBFhaE@vQ5f3lCq6vnfdRZ798rV4PSmqn"
	DummyRPICode14Length                          = "RPI_ABCDE12345"
	DummyCorrectCustomerID                        = "AYOPOP-12345"
	DummyInCorrectCustomerID                      = "AYO-12345"
	DummyCustomerIDSpecialChar                    = "AYO@OP-12345"
	DummyMerchantCode                             = "AYOPOP"
	DummyInCorrectMerchantCodeWithGreaterLength   = "AYOPOPO"
	DummyMerchantCodeWithSpecialCharacters        = "AYO@YO"
	DummyPhoneNumber                              = "6285322222685"
	DummyInvalidPhoneNumberWithWrongStartingDigit = "12134530001"
	DummyInvalidPhoneNumberWithLessLength         = "30001"
	DummyPhoneNumberWithGreaterLength             = "6213453000100000000"
	DummyPhoneNumberWithSpecialChar               = "6213453000@"
	DummyCurrency                                 = "IDR"
	DummySpecialCharCurrency                      = "I@R"
	DummyInvalidCurrency                          = "IRR"
	DummyAmount                                   = "100.35"
	DummySpecialCharAmount                        = "1@0"
	DummyInvalidAmount                            = "aaa"
	DummyBankCode                                 = "002"
	DummySpecialCharBankCode                      = "@02"
	DummyInvalidBankCode                          = "0000"
	DummyOtp                                      = "999999"
	DummyInvalidOtp                               = "9"
	DummySpecialCharacterOtp                      = "99@9"
	DummyMaskedCard                               = "************1234"
	DummySpecialCharMaskedCard                    = "1@23"
	DummyInvalidMaskedCard                        = "***********1234"
	DummyInvalidNumberOfOccurrences               = 30
	DummyEmail                                    = "ayopop@mail.com"
	IndonesiaRupiah                               = "IDR"
	DummyIpAddress                                = "255.255.111.35"
	DummyInvalidIpAddressSpecialCharacters        = "ab@.abc.abc.12"
	DummyInvalidLengthIpAddress                   = "1000.40.210.253"
	DummyBeneficiaryId                            = "BE_12345abcde"
	DummyBeneficiaryIdSpecialCharacters           = "BE_12345abcd@"
	DummyBeneficiaryIdLength                      = "BE_12345abcdef"
	DummyUserId                                   = "user_FGG13GT6Fg"
	DummyUserIdSpecialChar                        = "user_FGG13GT6F@"
	DummyInvalidLengthUserId                      = "user_FGG13GT6Fg1"
	DummyAccountNumber                            = "83581952"
	DummySwiftCode                                = "SIHBIDJ1"
	DummyInvalidSwiftCode                         = "333SIHBIDJ.."
)
