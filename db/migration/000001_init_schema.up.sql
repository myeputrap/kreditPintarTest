-- Create Consumer table
CREATE TABLE `consumers` (
    `id`                INT             NOT NULL AUTO_INCREMENT,
    `nik`               VARCHAR(20)     NOT NULL,
    `name`              VARCHAR(100)    NOT NULL,
    `place_of_birth`    VARCHAR(50)     NOT NULL,
    `birth_date`        DATE            NOT NULL,
    `salary`            DECIMAL(15, 2)  NOT NULL,
    `email`             VARCHAR(20)     NOT NULL,
    `phone_number`      VARCHAR(20)     NOT NULL,
    `dtm_crt`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dtm_upd`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY (`nik`) USING BTREE,
    KEY (`name`) USING BTREE,
    KEY (`birth_date`),
    KEY (`phone_number`) USING BTREE,
    KEY (`dtm_crt`),
    KEY (`dtm_upd`)
);

-- Create Transaction table
CREATE TABLE `transactions` (
    `id`                INT             NOT NULL AUTO_INCREMENT,
    `consumer_id`       INT             NOT NULL,
    `contract_number`   VARCHAR(50)     NOT NULL,
    `otr`               DECIMAL(15, 2)  NOT NULL,
    `admin_fee`         DECIMAL(15, 2)  NOT NULL,
    `installment_count` INT             NOT NULL,
    `interest_amount`   DECIMAL(15, 2)  NOT NULL,
    `purchase_amount`   DECIMAL(15, 2)  NOT NULL,
    `asset_name`        VARCHAR(100)    NOT NULL,
    `status`            VARCHAR(10)     NOT NULL,
    `dtm_crt`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dtm_upd`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY (`consumer_id`),
    KEY (`contract_number`) USING BTREE,
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`)
);

-- Create Credit Card table
CREATE TABLE `credit_cards` (
    `id`                INT             NOT NULL AUTO_INCREMENT,
    `consumer_id`       INT             NOT NULL,
    `card_number`       VARCHAR(16)     NOT NULL,
    `expiration_date`   DATE            NOT NULL,
    `cvv`               INT             NOT NULL,
    `credit_limit`      DECIMAL(15, 2)  NOT NULL,
    `current_balance`   DECIMAL(15, 2)  NOT NULL,
    `dtm_crt`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dtm_upd`           TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY (`consumer_id`),
    KEY (`card_number`) USING BTREE,
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`)
);

CREATE TABLE `billing` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `consumer_id` INT NOT NULL,
    `transaction_id` INT NOT NULL,
    `bill_amount` DECIMAL(15, 2) NOT NULL,
    `due_date` DATE NOT NULL,
    `status` VARCHAR(20) NOT NULL,
    `dtm_crt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dtm_upd` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY (`consumer_id`),
    KEY (`transaction_id`),
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`),
    FOREIGN KEY (`transaction_id`) REFERENCES `transactions`(`id`)
);

-- Create Admin table
CREATE TABLE `admin` (
    `user_name`         VARCHAR(50)     NOT NULL,
    `password`          VARCHAR(100)    NOT NULL,
    PRIMARY KEY (`user_name`)
);

INSERT INTO `admin` (`user_name`, `password`)
VALUES ('admin', '$2a$12$6B4sCTeUIZ9ctISGwb.6tODSgRUGLwQWQGPo0QqH6mdfkIyk55x8K');