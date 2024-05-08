CREATE USER IF NOT EXISTS 'datastream'@'%' IDENTIFIED BY '<password>';
GRANT REPLICATION SLAVE, SELECT, RELOAD, REPLICATION CLIENT, LOCK TABLES, EXECUTE ON *.* TO 'datastream'@'%';

CREATE USER IF NOT EXISTS 'bulk_disbursement_read'@'%' IDENTIFIED BY '<password>';
CREATE USER IF NOT EXISTS 'bulk_disbursement_write'@'%' IDENTIFIED BY '<password>';
GRANT SELECT ON bulk_disbursement.* TO 'bulk_disbursement_read'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE ON bulk_disbursement.* TO 'bulk_disbursement_write'@'%';
FLUSH PRIVILEGES;
COMMIT;