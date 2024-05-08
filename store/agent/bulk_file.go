package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"time"

	"bulk/api/v1/repositories"
	"bulk/constant"
	errors "bulk/grpc-errors"
	"bulk/logger"
	"bulk/models/request"
	disbursementpb "bulk/pb/disbursement"
	filepb "bulk/pb/file"
	querybuilder "bulk/pkg/query-builder"
	"bulk/tracer"
	"bulk/utils"

	"go.opentelemetry.io/otel/attribute"
)

type BulkFileRepository struct {
	DB     *sql.DB
	logger logger.ILogger
}

func NewBulkFileRepository(DB *sql.DB, log logger.ILogger) *BulkFileRepository {
	return &BulkFileRepository{DB: DB, logger: log}
}

func (s *BulkFileRepository) CreateBulkFile(
	ctx context.Context,
	req *repositories.CreateBulkFileRequest,
) (*repositories.BulkFileModel, error) {
	ctx, createBulkFileSpan := tracer.StartSpan(ctx, "CreateBulkFile")
	defer createBulkFileSpan.End()

	query := `
		INSERT INTO bulk_details (merchantCode, bulkId, name, fileName, filePath, fileSize, uploadedBy, status) VALUES 
		( ?, ?, ?, ?, ?, ?, ?, ?)
	`

	resp, err := s.DB.ExecContext(
		ctx,
		query,
		req.MerchantCode,
		req.BulkId,
		req.Name,
		req.FileName,
		req.FilePath,
		req.FileSize,
		req.UserId,
		req.Status.Int64(),
	)
	if err != nil {
		s.logger.Errorf("failed to insert bulk file: %v", err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	id, err := resp.LastInsertId()
	if err != nil {
		s.logger.Errorf("failed to get last insert id: %v", err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")

	}

	return &repositories.BulkFileModel{
		Id:       id,
		FileName: req.FileName,
		FilePath: req.FilePath,
		FileSize: req.FileSize,
		UserId:   req.UserId,
	}, nil
}

func (s *BulkFileRepository) UpdateBulkFile(
	ctx context.Context,
	req *repositories.UpdateBulkFileRequest,
) error {
	ctx, updateBulkFileSpan := tracer.StartSpan(ctx, "UpdateBulkFile")
	defer updateBulkFileSpan.End()

	query := `
		UPDATE bulk_details SET filePath = ?,status=? WHERE id = ?
`
	_, err := s.DB.ExecContext(ctx, query, req.FilePath, req.Status.Int64(), req.Id)

	if err != nil {
		s.logger.Errorf("failed to update bulk file: %v", err)
		return errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")

	}

	return nil
}

func (s *BulkFileRepository) GetBulkFile(ctx context.Context, req *repositories.GetBulkFileRequest) (*repositories.BulkFileModel, error) {
	ctx, getBulkFileSpan := tracer.StartSpan(ctx, "GetBulkFile")
	defer getBulkFileSpan.End()

	var (
		query string = `SELECT id,bulkId,name,fileName,fileSize,filePath,uploadedBy,status FROM bulk_details WHERE id = ?`
		value any    = req.Id
	)

	if req.Bulkid != "" {
		query = `SELECT id,bulkId,name,fileName,fileSize,filePath,uploadedBy,status FROM bulk_details WHERE bulkId = ?`
		value = req.Bulkid
	}

	res, err := s.DB.QueryContext(ctx, query, value)
	if err != nil {
		getBulkFileSpan.RecordError(err)
		s.logger.Errorf("failed to get last insert id: %v", err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	if !res.Next() {
		s.logger.Errorf("Error not found: %v", err)
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	var result repositories.BulkFileModel
	if err := res.Scan(
		&result.Id,
		&result.BulkID,
		&result.BulkName,
		&result.FileName,
		&result.FileSize,
		&result.FilePath,
		&result.UserId,
		&result.Status,
	); err != nil {
		getBulkFileSpan.RecordError(err)
		s.logger.Errorf("failed to scan: %v", err)
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	return &result, nil
}

func (s *BulkFileRepository) GetBulkFileByBulkID(ctx context.Context, bulkID string) (*repositories.BulkFileModel, error) {
	ctx, getBulkFileByBulkIDSpan := tracer.StartSpan(ctx, "GetBulkFileByBulkID")
	defer getBulkFileByBulkIDSpan.End()

	if bulkID == "" {
		err := errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
		s.logger.Errorf("bulkDisbursmentId should not be empty, request %s", bulkID)
		getBulkFileByBulkIDSpan.RecordError(err)
		return nil, err
	}

	query := `
		SELECT id,bulkId,name,fileName,fileSize,filePath,uploadedBy,status FROM bulk_details WHERE bulkId = ?
	`

	res, err := s.DB.QueryContext(ctx, query, bulkID)
	if err != nil {
		getBulkFileByBulkIDSpan.RecordError(err)
		return nil, fmt.Errorf("failed to query GetBulkFileByBulkID err:%w, query:%s", err, query)
	}

	if !res.Next() {
		s.logger.Errorf("Error not found: %v", err)
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	var result repositories.BulkFileModel
	if err := res.Scan(
		&result.Id,
		&result.BulkID,
		&result.BulkName,
		&result.FileName,
		&result.FileSize,
		&result.FilePath,
		&result.UserId,
		&result.Status,
	); err != nil {
		getBulkFileByBulkIDSpan.RecordError(err)
		s.logger.Errorf("failed to scan GetBulkFileByBulkID: %v", err)
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	return &result, nil
}

func (s *BulkFileRepository) GetBulkData(ctx context.Context, req *repositories.GetBulkFileRequest) ([]request.BulkDisbursementDetails, error) {
	ctx, getBulkDataSpan := tracer.StartSpan(ctx, "GetBulkData")
	defer getBulkDataSpan.End()
	query := ""
	var rows *sql.Rows
	var err error
	if req.Status != 0 {
		query = "SELECT bdd.id,bdd.bulkId,bdd.accountNumber,bdd.bankCode,bdd.phoneNumber,bdd.amount,bdd.customerId,bdd.beneficiaryCorrelationId,bdd.beneficiaryId,bdd.beneficiaryStatus,bdd.disbursementReferenceNumber,bdd.disbursementStatus,bdd.status,bdd.failedReason,bdd.paymentInfo,bdd.beneficiaryName,bdd.beneficiaryBankName FROM bulk_disbursement_details bdd inner join bulk_details bd on bdd.bulkId = bd.id WHERE bd.bulkId = ? and bdd.status = ?"
		rows, err = s.DB.QueryContext(ctx, query, req.Bulkid, req.Status)
	}
	if req.Status == 0 {
		query = "SELECT bdd.id,bdd.bulkId,bdd.accountNumber,bdd.bankCode,bdd.phoneNumber,bdd.amount,bdd.customerId,bdd.beneficiaryCorrelationId,bdd.beneficiaryId,bdd.beneficiaryStatus,bdd.disbursementReferenceNumber,bdd.disbursementStatus,bdd.status,bdd.failedReason,bdd.paymentInfo,bdd.beneficiaryName,bdd.beneficiaryBankName FROM bulk_disbursement_details bdd inner join bulk_details bd on bdd.bulkId = bd.id WHERE bd.bulkId = ?"
		rows, err = s.DB.QueryContext(ctx, query, req.Bulkid)
	}
	if err != nil {
		getBulkDataSpan.RecordError(err)
		s.logger.Errorf("failed to query: %v", err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")

	}
	defer rows.Close()
	var bulkDataList []request.BulkDisbursementDetails
	for rows.Next() {
		var result request.BulkDisbursementDetails
		if err := rows.Scan(
			&result.Id,
			&result.BulkId,
			&result.AccountNumber,
			&result.BankCode,
			&result.PhoneNumber,
			&result.Amount,
			&result.CustomerId,
			&result.BeneficiaryCorrelationId,
			&result.BenficiaryId,
			&result.BeneficiaryStatus,
			&result.DisbursementReferenceNo,
			&result.DisbursementStatus,
			&result.Status,
			&result.FailedReason,
			&result.PaymentInfo,
			&result.BenficiaryName,
			&result.BeneficiaryBankName,
		); err != nil {
			getBulkDataSpan.RecordError(err)
			s.logger.Errorf("failed to scan: %v", err)
			return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")

		}
		bulkDataList = append(bulkDataList, result)

	}

	return bulkDataList, nil
}

func (s *BulkFileRepository) UpdateBulkDisbursementData(ctx context.Context, req *request.UpdateDisbursement) (int64, error) {
	ctx, updateBulkDisbursementDataSpan := tracer.StartSpan(ctx, "UpdateBulkDisbursementData")
	defer updateBulkDisbursementDataSpan.End()
	var (
		bulkDisbursementDetailsId        int64
		bulkDetailsUpdatedStatus, bulkId int
		reason, updatedBy                string
	)
	updateQuery := "UPDATE bulk_disbursement_details"
	querybuilder := querybuilder.NewQueryBuilder(updateQuery)
	status := strconv.Itoa(int(*req.Status))
	querybuilder.Set("status").Equals(status)
	if req.CustomerId != nil && len(*req.CustomerId) > 0 {
		querybuilder.AddComma("customerId").Equals(*req.CustomerId)
	}
	if req.BeneficiaryId != nil && len(*req.BeneficiaryId) > 0 {
		querybuilder.AddComma("beneficiaryId").Equals(*req.BeneficiaryId)
	}
	if req.BeneficiaryStatus != nil && *req.BeneficiaryStatus > 0 {
		status := strconv.Itoa(int(*req.BeneficiaryStatus))
		querybuilder.AddComma("beneficiaryStatus").Equals(status)
	}
	if req.DisbursementReferenceNumber != nil && len(*req.DisbursementReferenceNumber) > 0 {
		querybuilder.AddComma("disbursementReferenceNumber").Equals(*req.DisbursementReferenceNumber)
	}
	if req.DisbursementStatus != nil && *req.DisbursementStatus > 0 {
		status := strconv.Itoa(int(*req.DisbursementStatus))
		querybuilder.AddComma("disbursementStatus").Equals(status)
	}

	if req.FailedReason != nil && len(*req.FailedReason) > 0 {
		querybuilder.AddComma("failedReason").Equals(*req.FailedReason)
	}

	if req.BenficiaryName != nil && len(*req.BenficiaryName) > 0 {
		querybuilder.AddComma("beneficiaryName").Equals(*req.BenficiaryName)
	}
	if req.BeneficiaryBankName != nil && len(*req.BeneficiaryBankName) > 0 {
		querybuilder.AddComma("beneficiaryBankName").Equals(*req.BeneficiaryBankName)
	}

	dbTx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		updateBulkDisbursementDataSpan.RecordError(err)
		s.logger.Errorf("err while getting db transaction %v", err)
		return bulkDisbursementDetailsId, err
	}

	switch req.Type {
	case "INQUIRY":
		querybuilder.Where("beneficiaryCorrelationId").Equals(*req.BeneficiaryCorrelationId)
	case "INQUIRY_STATUS":
		querybuilder.Where("beneficiaryCorrelationId").Equals(*req.BeneficiaryCorrelationId)
		var status, vaildEntries, totalInquiries, verifiedInquiries int
		err = dbTx.QueryRowContext(ctx, "SELECT bulkId, status FROM bulk_disbursement_details WHERE beneficiaryCorrelationId = ?", *req.BeneficiaryCorrelationId).Scan(&bulkId, &status)
		if err != nil && err != sql.ErrNoRows {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while getting bulk data %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (get bulkData) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		// bulk_disbursement_details status is already changed to inquiry success (2)
		if status >= 2 {
			return bulkDisbursementDetailsId, nil
		}

		if _, err := dbTx.ExecContext(ctx, "UPDATE bulk_entries SET totalInquiries = totalInquiries + 1 WHERE bulkId = ?", bulkId); err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while updating the totalInquiries %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (update varifiedInquiries) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		// update verified inquiry if inquiry is finished
		if *req.Status == 2 {
			if _, err := dbTx.ExecContext(ctx, "UPDATE bulk_entries SET verifiedInquiries = verifiedInquiries + 1 WHERE bulkId = ?", bulkId); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while updating the varified inquiry %v", err)

				if err := dbTx.Rollback(); err != nil {
					updateBulkDisbursementDataSpan.RecordError(err)
					s.logger.Errorf("err while rollback transaction (update varifiedInquiries) %v", err)
					return bulkDisbursementDetailsId, err
				}
				return bulkDisbursementDetailsId, err
			}
		}

		err = dbTx.QueryRowContext(ctx, "SELECT vaildEntries, totalInquiries, verifiedInquiries FROM bulk_entries WHERE bulkId = ?", bulkId).Scan(&vaildEntries, &totalInquiries, &verifiedInquiries)
		if err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while select bulk entries, err: %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (select bulk entries) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		if vaildEntries == totalInquiries {
			bulkDetailsUpdatedStatus = constant.BulkStatusVerified
			if verifiedInquiries == 0 {
				bulkDetailsUpdatedStatus = constant.BulkStatusFailed
				reason = "All Inquiry Failed"
				updatedBy = "GENERIC_INQUIRY_CALLBACK"
			}
		}

	case "DISBURSEMENT":
		id := strconv.FormatInt(*req.Id, 10)
		querybuilder.Where("id").Equals(id)
	case "DISBURSEMENT_STATUS":
		var (
			selectQuery                                             string
			selectQueryArg                                          any
			totalDisbursements, verifiedInquiries, disbursedEntries int
		)
		if req.Id != nil {
			id := strconv.FormatInt(*req.Id, 10)
			querybuilder.Where("id").Equals(id)
			selectQuery = "SELECT amount, bulkId, status FROM bulk_disbursement_details WHERE id = ?"
			selectQueryArg = id
		} else {
			querybuilder.Where("disbursementReferenceNumber").Equals(*req.DisbursementReferenceNumber)
			selectQuery = "SELECT amount, bulkId, status FROM bulk_disbursement_details WHERE disbursementReferenceNumber = ?"
			selectQueryArg = *req.DisbursementReferenceNumber
		}

		var amount float64 = 0
		var status int
		err = dbTx.QueryRowContext(ctx, selectQuery, selectQueryArg).Scan(&amount, &bulkId, &status)
		if err != nil && err != sql.ErrNoRows {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while getting bulk data %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (get bulkData) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		// bulk_disbursement_details status is already changed to transfer success (4)
		if status >= 4 {
			return bulkDisbursementDetailsId, nil
		}

		if _, err := dbTx.ExecContext(ctx, "UPDATE bulk_entries SET totalDisbursements = totalDisbursements + 1 WHERE bulkId = ?", bulkId); err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while updating the totalDisbursements %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (update disbursedEntries) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		// update dibursedEntries if transfer is finished
		if *req.Status == 4 {
			if _, err := dbTx.ExecContext(ctx, "UPDATE bulk_entries SET disbursedEntries = disbursedEntries + 1, disbursedAmount = disbursedAmount + ? WHERE bulkId = ?", amount, bulkId); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while updating the disbursedEntries and disbursedAmount %v", err)

				if err := dbTx.Rollback(); err != nil {
					updateBulkDisbursementDataSpan.RecordError(err)
					s.logger.Errorf("err while rollback transaction (update disbursedEntries) %v", err)
					return bulkDisbursementDetailsId, err
				}
				return bulkDisbursementDetailsId, err
			}
		}

		err = dbTx.QueryRowContext(ctx, "SELECT totalDisbursements, verifiedInquiries, disbursedEntries FROM bulk_entries WHERE bulkId = ?", bulkId).Scan(&totalDisbursements, &verifiedInquiries, &disbursedEntries)
		if err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while select bulk entries for disbursement_status, err: %v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (select bulk entries for disbursement_status) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		if totalDisbursements == verifiedInquiries {
			bulkDetailsUpdatedStatus = constant.BulkStatusCompleted
			if disbursedEntries == 0 {
				bulkDetailsUpdatedStatus = constant.BulkStatusFailed
				reason = "All Disbursement Failed"
				updatedBy = "GENERIC_DISBURSEMENT_CALLBACK"
			}
		}
	}

	if bulkDetailsUpdatedStatus != 0 {
		// get bulk file because we need current bulk status for bulk_status_log table
		bulkfile, err := s.GetBulkFile(ctx, &repositories.GetBulkFileRequest{
			Id: int64(bulkId),
		})
		if err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (get bulk file) %v", err)
				return bulkDisbursementDetailsId, err
			}
		}

		if _, err := dbTx.ExecContext(ctx, "UPDATE bulk_details SET status = ? WHERE id = ?", bulkDetailsUpdatedStatus, bulkId); err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while updating bulk_details status, err: %+v", err)

			if err := dbTx.Rollback(); err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf("err while rollback transaction (update bulk_details status) %v", err)
				return bulkDisbursementDetailsId, err
			}
			return bulkDisbursementDetailsId, err
		}

		if bulkDetailsUpdatedStatus == constant.BulkStatusFailed {
			_, err = dbTx.ExecContext(ctx, constant.InsertQueryForStatusLog, bulkId, reason, updatedBy, bulkfile.Status, bulkDetailsUpdatedStatus)
			if err != nil {
				updateBulkDisbursementDataSpan.RecordError(err)
				s.logger.Errorf(constant.ErrInsertStatusLog, constant.InsertQueryForStatusLog, err)

				if err := dbTx.Rollback(); err != nil {
					updateBulkDisbursementDataSpan.RecordError(err)
					s.logger.Errorf("err while rollback transaction (inserting bulk_status_logs) %v", err)
					return bulkDisbursementDetailsId, err
				}
				return bulkDisbursementDetailsId, err
			}
		}
	}

	if err := querybuilder.Build(); err != nil {
		s.logger.Errorf("err while building the query %v", err)
	}

	s.logger.Debugf("Query string :", querybuilder.GetQuery())
	queryUpdateData := querybuilder.GetQuery()

	updateResult, err := dbTx.ExecContext(ctx, queryUpdateData)
	if err != nil {
		updateBulkDisbursementDataSpan.RecordError(err)
		s.logger.Errorf("Error while Executiing statement UPDATE bulk disbursement details: %v", err)

		if err := dbTx.Rollback(); err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while rollback transaction (update bulkdisbursement) %v", err)
			return bulkDisbursementDetailsId, err
		}
		return bulkDisbursementDetailsId, err
	}

	if err = dbTx.Commit(); err != nil {
		updateBulkDisbursementDataSpan.RecordError(err)
		s.logger.Errorf("Error while commiting UPDATE bulk disbursement details changes: %v", err)

		if err := dbTx.Rollback(); err != nil {
			updateBulkDisbursementDataSpan.RecordError(err)
			s.logger.Errorf("err while rollback transaction (update bulkdisbursement) %v", err)
			return bulkDisbursementDetailsId, err
		}
		return bulkDisbursementDetailsId, err
	}

	bulkDisbursementDetailsId, _ = updateResult.LastInsertId()
	updateBulkDisbursementDataSpan.SetAttributes(attribute.Int64(".ID", bulkDisbursementDetailsId))

	return bulkDisbursementDetailsId, nil
}

func (s *BulkFileRepository) UpdateBulkFileStatus(
	ctx context.Context,
	req *filepb.UpdateBulkStatusRequest,
) (int64, error) {
	ctx, updateBulkFileStatusSpan := tracer.StartSpan(ctx, "UpdateBulkFileStatus")
	defer updateBulkFileStatusSpan.End()

	var bulkFileid int64

	bulkData, err := s.GetBulkFile(ctx, &repositories.GetBulkFileRequest{
		Bulkid: req.BulkId,
	})
	if err != nil {
		updateBulkFileStatusSpan.RecordError(err)
		return bulkFileid, err
	}
	if req.Status == 11 {
		if bulkData.Status == 5 {
			s.logger.Debugf("Updating status to reject for this bulkId", req.BulkId)
		} else {
			return bulkFileid, errors.GetError(ctx, errors.BadRequestErrMsg, "InvalidTransaction0414")
		}
	}
	query := `
		UPDATE bulk_details SET status=? WHERE  bulkId= ?		
	`
	result, err := s.DB.ExecContext(ctx, query, req.Status, req.BulkId)
	if err != nil {
		updateBulkFileStatusSpan.RecordError(err)
		s.logger.Errorf("failed to update statu query: %v, err :%v", query, err)
		return bulkFileid, err
	}
	bulkFileid, _ = result.LastInsertId()
	updateBulkFileStatusSpan.SetAttributes(attribute.Int64(".ID", bulkFileid))
	var reason string = "N/A"
	if req.GetReason() != "" {
		reason = req.GetReason()
	}

	if utils.ValidateStatusForBulkStatusLogs(int(req.Status)) {
		// insert the data only if status is rejected and failed
		_, err = s.DB.ExecContext(ctx, constant.InsertQueryForStatusLog, bulkData.Id, reason, req.GetMerchantUserId(), bulkData.Status, req.Status)
		if err != nil {
			updateBulkFileStatusSpan.RecordError(err)
			s.logger.Errorf(constant.ErrInsertStatusLog, constant.InsertQueryForStatusLog, err)
			return bulkFileid, err
		}
	}

	return bulkFileid, nil
}

func (s *BulkFileRepository) GetDisbursements(ctx context.Context, bulkID int, req *disbursementpb.DisbursementsRequest) ([]*disbursementpb.Disbursement, int, error) {
	ctx, getDisbursementsSpan := tracer.StartSpan(ctx, "GetDisbursements")
	defer getDisbursementsSpan.End()

	count := 0

	if req == nil {
		err := errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
		s.logger.Errorf("req should not be empty, request %+v", req)
		s.logger.Errorf("Error building query %v", err)
		getDisbursementsSpan.RecordError(err)
		return nil, count, err
	}

	if req.BulkDisbursmentId == "" {
		err := errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
		s.logger.Errorf("bulkDisbursmentId should not be empty, request %+v", req)
		getDisbursementsSpan.RecordError(err)
		return nil, count, err
	}

	columns := "accountNumber,bankCode,amount,status,beneficiaryId, IF(failedReason = 'N/A', NULL, failedReason)"

	query := "SELECT %s FROM bulk_disbursement_details"

	disbursementQB := querybuilder.NewQueryBuilder(query)

	// adding the filters
	applyFiltersForDisbursementList(disbursementQB, bulkID, req)

	// building the query
	err := disbursementQB.Build()
	if err != nil {
		getDisbursementsSpan.RecordError(err)
		s.logger.Errorf("Error building query %v", err)
		return nil, count, fmt.Errorf("error building query: %v", err)
	}

	countQuery := fmt.Sprintf(disbursementQB.GetQuery(), "count(*)")

	s.logger.Debugf("GetDisbursements count query: %v", countQuery)

	if err := s.DB.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		getDisbursementsSpan.RecordError(err)
		s.logger.Errorf("Error querying disbursements count %s, err: %w", countQuery, err)
		return nil, count, err
	}

	// add sort, limit and offset as it will effect the count query so adding them after the count query
	disbursementQB.Sort("created", true)

	if req.Limit > 0 {
		disbursementQB.Limit(int(req.Limit))

		if req.Page > 0 {
			offset := (req.Page - 1) * req.Limit
			disbursementQB.Offset(int(offset))
		}
	}

	disbursementQuery := fmt.Sprintf(disbursementQB.GetQuery(), columns)

	s.logger.Debugf("GetDisbursements query: %v", disbursementQuery)

	rows, err := s.DB.QueryContext(ctx, disbursementQuery)
	if err != nil {
		getDisbursementsSpan.RecordError(err)
		s.logger.Errorf("Error querying disbursements %s, err: %w", disbursementQuery, err)
		return nil, count, err
	}

	disbursements := []*disbursementpb.Disbursement{}
	for rows.Next() {
		disbursement := new(disbursementpb.Disbursement)
		if err := rows.Scan(
			&disbursement.AccountNumber,
			&disbursement.BeneficiaryBank,
			&disbursement.Amount,
			&disbursement.Status,
			&disbursement.BeneficiaryId,
			&disbursement.StatusFailedReason,
		); err != nil {
			getDisbursementsSpan.RecordError(err)
			s.logger.Errorf("Error querying disbursements %s", disbursementQuery)
			return nil, count, fmt.Errorf("failed to scan disbursements %w", err)
		}
		disbursements = append(disbursements, disbursement)
	}

	return disbursements, count, nil
}

func (s *BulkFileRepository) GetBulkDisbursementStatusCount(ctx context.Context, bulkID int) (*disbursementpb.BulkDisbursementStatusCount, error) {
	ctx, getBulkDisbursementStatusCountSpan := tracer.StartSpan(ctx, "GetBulkDisbursementStatusCount")
	defer getBulkDisbursementStatusCountSpan.End()

	query := `SELECT IFNULL(vaildEntries,0), IFNULL(verifiedInquiries,0), IFNULL(disbursedEntries,0) FROM bulk_entries WHERE bulkID = ?`

	statusCounts := new(disbursementpb.BulkDisbursementStatusCount)

	err := s.DB.QueryRowContext(ctx, query, bulkID).Scan(&statusCounts.VaildEntries, &statusCounts.VerifiedInquiries, &statusCounts.DisbursedEntries)
	if err != nil && err != sql.ErrNoRows {
		getBulkDisbursementStatusCountSpan.RecordError(err)
		s.logger.Errorf("Error querying GetBulkDisbursementStatusCount %s, err: %w", query, err)
		return nil, err
	}

	return statusCounts, nil
}

func (s *BulkFileRepository) GetBulkList(ctx context.Context, req *filepb.BulkFilesRequest) ([]*filepb.BulkFile, int, error) {
	getBulkListCtx, getBulkListSpan := tracer.StartSpan(ctx, "GetBulkList")
	defer getBulkListSpan.End()

	count := 0

	if req.MerchantCode == "" {
		err := errors.GetError(ctx, errors.BadRequestErrMsg, "MerchantCodeIsInvalid0320")
		getBulkListSpan.RecordError(err)
		s.logger.Errorf("Error validating request %v", err)
		return nil, count, err
	}

	columns := "bulk_details.created, bulk_details.bulkId, bulk_details.name, IFNULL(bulk_entries.disbursedAmount,0), bulk_details.uploadedBy, bulk_details.status, IFNULL(bulk_entries.vaildEntries,0), IFNULL(bulk_entries.verifiedInquiries,0), IFNULL(bulk_entries.disbursedEntries,0), IFNULL(bulk_status_logs.reason, 'N/A'), IFNULL(bulk_status_logs.updatedBy, 'N/A')"

	query := "SELECT %s FROM bulk_details LEFT JOIN bulk_entries ON bulk_details.id = bulk_entries.bulkId LEFT JOIN bulk_status_logs ON bulk_status_logs.bulkId = bulk_details.id AND bulk_status_logs.updatedStatus IN (9, 11)" // fetch only failed as rejected status

	bulkListDB := querybuilder.NewQueryBuilder(query)

	// apply filters on bulkListDB
	if err := applyFiltersForBulkList(bulkListDB, req); err != nil {
		getBulkListSpan.RecordError(err)
		s.logger.Errorf("error applying filters for bulk list: %v", err)
		return nil, count, err
	}

	err := bulkListDB.Build()
	if err != nil {
		getBulkListSpan.RecordError(err)
		s.logger.Errorf("Error building bulkListDB %v", err)
		return nil, count, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	countQuery := fmt.Sprintf(bulkListDB.GetQuery(), "count(*)")

	s.logger.Debugf("GetBulkList count query: %v", countQuery)

	if err := s.DB.QueryRowContext(getBulkListCtx, countQuery).Scan(&count); err != nil {
		getBulkListSpan.RecordError(err)
		s.logger.Errorf("Error querying bulkList count %s, err %v", countQuery, err)
		return nil, count, err
	}

	bulkListDB.Sort("created", true)

	if req.Limit > 0 {
		bulkListDB.Limit(int(req.Limit))

		if req.Page > 0 {
			offset := (req.Page - 1) * req.Limit
			bulkListDB.Offset(int(offset))
		}
	}

	bulkListQuery := fmt.Sprintf(bulkListDB.GetQuery(), columns)

	s.logger.Debugf("GetBulkList query: %v", bulkListQuery)

	rows, err := s.DB.QueryContext(getBulkListCtx, bulkListQuery)
	if err != nil {
		getBulkListSpan.RecordError(err)
		s.logger.Errorf("Error querying bulkList %s, err: %v", bulkListQuery, err)
		return nil, count, err
	}

	bulkFiles := []*filepb.BulkFile{}
	for rows.Next() {
		var createAt time.Time
		bulkFile := new(filepb.BulkFile)
		bulkFile.StatusCount = new(filepb.BulkFileStatusCount)
		if err := rows.Scan(
			&createAt,
			&bulkFile.BulkDisbursementId,
			&bulkFile.BulkDisbursementName,
			&bulkFile.TotalAmount,
			&bulkFile.Uploader,
			&bulkFile.Status,
			&bulkFile.StatusCount.VaildEntries,
			&bulkFile.StatusCount.VerifiedInquiries,
			&bulkFile.StatusCount.DibursedEntries,
			&bulkFile.Reason,
			&bulkFile.UpdatedBy,
		); err != nil {
			getBulkListSpan.RecordError(err)
			s.logger.Errorf("Error scanning bulkList %s, err: %v", bulkListQuery, err)
			return nil, count, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
		}

		bulkFile.CreatedDateTime = createAt.Format(constant.BulkFileListTimeFormat)
		bulkFiles = append(bulkFiles, bulkFile)
	}

	return bulkFiles, count, nil
}

func applyFiltersForBulkList(bulkListDB *querybuilder.QueryBuilder, req *filepb.BulkFilesRequest) error {

	bulkListDB.Where("bulk_details.merchantCode").Equals(req.MerchantCode)

	if req.BulkDisbursementId != nil && *req.BulkDisbursementId != "" {
		bulkListDB.And("bulk_details.bulkId").Like("%%" + *req.BulkDisbursementId + "%%")
	}

	if req.BulkDisbursementName != nil && *req.BulkDisbursementName != "" {
		bulkListDB.And("bulk_details.name").Like("%%" + *req.BulkDisbursementName + "%%")
	}

	if req.Uploader != nil && *req.Uploader != "" {
		bulkListDB.And("bulk_details.uploadedBy").Like("%%" + *req.Uploader + "%%")
	}

	if req.Status != nil {
		bulkListDB.And("bulk_details.status").Equals(strconv.Itoa(int(*req.Status)))
	}

	if req.StartDate != nil && *req.StartDate != "" {
		sd, err := time.Parse(constant.DateFormat, *req.StartDate)
		if err != nil {
			return err
		}

		bulkListDB.And("DATE(bulk_details.created)").GreaterThanEquals(sd.UTC().Format(constant.DateFormat))
	}

	if req.EndDate != nil && *req.EndDate != "" {
		ed, err := time.Parse(constant.DateFormat, *req.EndDate)
		if err != nil {
			return err
		}

		bulkListDB.And("DATE(bulk_details.created)").LessThanEquals(ed.UTC().Format(constant.DateFormat))
	}

	return nil
}

func applyFiltersForDisbursementList(disbursementQB *querybuilder.QueryBuilder, bulkID int, req *disbursementpb.DisbursementsRequest) {

	disbursementQB.Where("bulkId").Equals(strconv.Itoa((bulkID)))

	if req.Status != nil {
		disbursementQB.And("status").Equals(strconv.Itoa(int(*req.Status)))
	}

	if req.BeneficiaryId != nil && *req.BeneficiaryId != "" {
		disbursementQB.And("beneficiaryId").Like("%%" + *req.BeneficiaryId + "%%")
	}

	if req.AccountNumber != nil && *req.AccountNumber != "" {
		disbursementQB.And("accountNumber").Like("%%" + *req.AccountNumber + "%%")
	}

	if req.BeneficiaryBank != nil && *req.BeneficiaryBank != "" {
		disbursementQB.And("bankCode").Like("%%" + *req.BeneficiaryBank + "%%")
	}
}
