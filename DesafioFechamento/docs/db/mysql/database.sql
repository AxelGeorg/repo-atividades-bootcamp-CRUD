-- DDL
DROP DATABASE IF EXISTS `fantasy_products`;

CREATE DATABASE `fantasy_products`;

USE `fantasy_products`;

-- Table structure for table `customers`
CREATE TABLE `customers` (
    `id` int NOT NULL AUTO_INCREMENT,
    `first_name` varchar(45) DEFAULT NULL,
    `last_name` varchar(45) DEFAULT NULL,
    `condition` tinyint(1) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

-- Table structure for table `invoices`
CREATE TABLE `invoices` (
    `id` int NOT NULL AUTO_INCREMENT,
    `datetime` datetime DEFAULT NULL,
    `customer_id` int DEFAULT NULL,
    `total` float DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_invoices_customer_id` (`customer_id`),
    CONSTRAINT `fk_invoices_customer_id` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table structure for table `products`
CREATE TABLE `products` (
    `id` int NOT NULL AUTO_INCREMENT,
    `description` varchar(100) DEFAULT NULL,
    `price` float DEFAULT NULL,
    PRIMARY KEY (`id`)
);

-- Table structure for table `sales`
CREATE TABLE `sales` (
    `id` int NOT NULL AUTO_INCREMENT,
    `quantity` int DEFAULT NULL,
    `invoice_id` int DEFAULT NULL,
    `product_id` int DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_sales_invoice_id` (`invoice_id`),
    KEY `idx_sales_product_id` (`product_id`),
    CONSTRAINT `fk_sales_invoice_id` FOREIGN KEY (`invoice_id`) REFERENCES `invoices` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_sales_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

/*
-- Add data to table customers
INSERT INTO `customers` (`first_name`,`last_name`,`condition`) VALUES
('John','Doe',1),
('Jane','Doe',1),
('John','Smith',1),
('Jane','Smith',1),
('John','Doe',0),
('Jane','Doe',0),
('John','Smith',0),
('Jane','Smith',0);

-- Add data to table invoices
INSERT INTO `invoices` (`datetime`,`customer_id`,`total`) VALUES
('2019-01-01',1,100.00),
('2019-01-01',2,200.00),
('2019-01-01',3,300.00),
('2019-01-01',4,400.00),
('2019-01-01',5,500.00),
('2019-01-01',6,600.00),
('2019-01-01',7,700.00),
('2019-01-01',8,800.00);

-- Add data to table products
INSERT INTO `products` (`description`,`price`) VALUES
('Product 1',100.00),
('Product 2',200.00),
('Product 3',300.00),
('Product 4',400.00),
('Product 5',500.00),
('Product 6',600.00),
('Product 7',700.00),
('Product 8',800.00);

-- Add data to table sales
INSERT INTO `sales` (`quantity`,`invoice_id`,`product_id`) VALUES
(1,1,1),
(2,2,2),
(3,3,3),
(4,4,4),
(5,5,5),
(6,6,6),
(7,7,7),
(8,8,8);
*/