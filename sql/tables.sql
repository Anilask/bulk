CREATE DATABASE IF NOT EXISTS `bulk_disbursement`;

CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_details (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    merchantCode VARCHAR(6) NOT NULL,
    bulkId VARCHAR(32) NOT NULL UNIQUE,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `fileName` VARCHAR(50) NOT NULL,
    `filePath` VARCHAR(250) NOT NULL,
    `fileSize` VARCHAR(16) NOT NULL,
    uploadedBy VARCHAR(200)NOT NULL, 
    `status` TINYINT(2) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT bulk_details_PK PRIMARY KEY (`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_disbursement_details (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    bulkId BIGINT(10) UNSIGNED NOT NULL,
    accountNumber VARCHAR(20) NOT NULL,
    bankCode VARCHAR(16) NOT NULL,
    phoneNumber VARCHAR(16) NOT NULL,
    amount decimal(18,2)  NOT NULL,
    customerId VARCHAR(16) NOT NULL,
    beneficiaryCorrelationId VARCHAR(32) NOT NULL,
    beneficiaryId VARCHAR(16) nOT NULL,
    beneficiaryStatus TINYINT(1) NOT NULL,
    disbursementReferenceNumber VARCHAR(32) NOT NULL,
    disbursementStatus TINYINT(1) NOT NULL,
    failedReason VARCHAR(16) NOT NULL,
    `status` TINYINT(2) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT bulk_disbursement_details_PK PRIMARY KEY (`id`),
    CONSTRAINT bulk_disbursement_details_FK FOREIGN KEY (bulkId) REFERENCES bulk_disbursement.bulk_details(`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_entries (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    bulkId BIGINT(10) UNSIGNED NOT NULL,
    totalInquiries INT(10) NOT NULL,
    totalDisbursements INT(10) NOT NULL,
    vaildEntries INT(10) NOT NULL,
    verifiedInquiries INT(10) NOT NULL,
    disbursedEntries INT(10) NOT NULL,
    disbursedAmount decimal(18,2)  NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT bulk_entries_PK PRIMARY KEY (`id`),
    CONSTRAINT bulk_entries_FK FOREIGN KEY (bulkId) REFERENCES bulk_disbursement.bulk_details(`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;


CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_status_logs (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    bulkId BIGINT(10) UNSIGNED NOT NULL,
    reason VARCHAR(200) NOT NULL,
    updatedBy VARCHAR(200)NOT NULL, 
    currentStatus TINYINT(2) NOT NULL,
    updatedStatus TINYINT(2) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT bulk_status_logs_PK PRIMARY KEY (`id`),
    CONSTRAINT bulk_status_logs_FK FOREIGN KEY (bulkId) REFERENCES bulk_disbursement.bulk_details(`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_details_status (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    `status` TINYINT(2) UNSIGNED NOT NULL,
    `description` VARCHAR(50) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT bulk_details_status_PK PRIMARY KEY (`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS bulk_disbursement.bulk_disbursement_details_status (
    `id` BIGINT(10) UNSIGNED AUTO_INCREMENT NOT NULL,
    `status` TINYINT(2) UNSIGNED NOT NULL,
    `description` VARCHAR(50) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT bulk_disbursement_details_status_PK PRIMARY KEY (`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

INSERT INTO bulk_details_status (`status`, `description`) VALUES 
    (0, "Created"),
    (1, "File Uploading"),
    (2, "File Uploaded"),
    (3, "Inquiry Initiated"),
    (4, "Verifying"),
    (5, "Verified"),
    (6, "Transfer Initiated"),
    (7, "In Progress"),
    (8, "Completed"),
    (9, "Failed"),
    (10, "Disabled"),
    (11, "Rejected");

INSERT INTO bulk_disbursement_details_status (`status`, `description`) VALUES
    (0, "Created"),
    (1, "Inquiry Initiated"),
    (2, "Inquiry Success"),
    (3, "Transfer Initiated"),
    (4, "Transfer Success"),
    (5, "Failed");

ALTER TABLE `bulk_disbursement`.`bulk_disbursement_details` 
ADD COLUMN `paymentInfo` VARCHAR(250) NULL DEFAULT 'N/A' AFTER `amount`;
ALTER TABLE `bulk_disbursement`.`bulk_disbursement_details` 
ADD COLUMN `beneficiaryName` VARCHAR(250) NULL DEFAULT 'N/A' AFTER `beneficiaryId`;
ALTER TABLE `bulk_disbursement`.`bulk_disbursement_details` 
ADD COLUMN `beneficiaryBankName` VARCHAR(100) NULL DEFAULT 'N/A' AFTER `beneficiaryName`;