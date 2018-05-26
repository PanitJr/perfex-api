CREATE SCHEMA `perfex`;

USE `perfex`;

CREATE TABLE `clinic_products` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `related_products` TEXT NULL,
  `price` DECIMAL(18,2) NOT NULL,
  `status` VARCHAR(30) NOT NULL,
  `info` VARCHAR(500) NULL,
PRIMARY KEY (`id`));

CREATE TABLE `products` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `category` VARCHAR(255) NULL,
  `brand` VARCHAR(255) NULL,
  `price` DECIMAL(18,2) NOT NULL,
  `amount` DECIMAL(18,2) NOT NULL,
  `info` VARCHAR(500) NULL,
  `unit` VARCHAR(50) NULL,
  `last_update` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`));

CREATE TABLE `clinic_records` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `customer_id` INT NOT NULL,
  `clinic_product_id` INT NOT NULL,
  `employee_id` INT NULL,
  `doctor_id` INT NULL,
  `price` DECIMAL(18,2) NOT NULL,
  `discount` DECIMAL(18,2) NULL DEFAULT 0.00,
  `paid` DECIMAL(18,2) NULL DEFAULT 0.00,
  `left` DECIMAL(18,2) NULL DEFAULT 0.00,
  `create_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  `last_update` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`));

CREATE TABLE `doctors` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `hour_rate` DECIMAL(18,2) NOT NULL,
  `dt_rate` DECIMAL(18,2) NOT NULL,
  `phone`VARCHAR(255) NULL,
  `email`VARCHAR(255) NULL,
  `info` TEXT NULL,
PRIMARY KEY (`id`));

CREATE TABLE `customers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `gender` VARCHAR(30) NOT NULL,
  `type` VARCHAR(30) NOT NULL DEFAULT 'NORMAL',
  `national_id` VARCHAR(30) NOT NULL UNIQUE,
  `phone`VARCHAR(255) NULL,
  `email`VARCHAR(255) NULL,
  `date_of_birth` DATETIME NOT NULL,
  `info` TEXT NULL,
  `join_date` DATETIME NOT NULL,
  `last_update` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`));

CREATE TABLE `employees` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `hour_rate` INT NOT NULL,
  `phone`VARCHAR(255) NULL,
  `email`VARCHAR(255) NULL,
  `info` TEXT NULL,
PRIMARY KEY (`id`));

CREATE TABLE `doctor_timesheets` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `doctor_id` INT NOT NULL,
  `check_in` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `check_out` DATETIME NULL,
  `calculated_dt` DECIMAL(18,2) NULL DEFAULT 0,
  `calculated_hour` DECIMAL(18,2) NULL DEFAULT 0,
  `total` DECIMAL(18,2) NULL  DEFAULT 0,
  `aditional` DECIMAL(18,2) NULL  DEFAULT 0,
  `status` VARCHAR(30) NOT NULL  DEFAULT 'CHECKIN',
  `info` TEXT NULL,
PRIMARY KEY (`id`));

CREATE TABLE `employee_timesheets` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `employee_id` INT NOT NULL,
  `check_in` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `check_out` DATETIME NULL,
  `calculated_hour` DECIMAL(18,2) NULL DEFAULT 0,
  `total` DECIMAL(18,2) NULL  DEFAULT 0,
  `aditional` DECIMAL(18,2) NULL  DEFAULT 0,
  `status` VARCHAR(30) NOT NULL  DEFAULT 'CHECKIN',
  `info` TEXT NULL,
PRIMARY KEY (`id`));