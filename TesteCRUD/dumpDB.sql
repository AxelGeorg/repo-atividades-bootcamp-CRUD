CREATE DATABASE `my_db3`; 
USE `my_db3`;

-- Excluindo tabelas se existirem
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `warehouses`;

-- Criando a tabela `warehouses`
CREATE TABLE `warehouses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `address` varchar(150) NOT NULL,
  `telephone` varchar(150) NOT NULL,
  `capacity` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Inserindo um registro na tabela `warehouses`
INSERT INTO `warehouses` (`name`, `address`, `telephone`, `capacity`) VALUES
('Main Warehouse', '221 Baker Street', '4555666', 100);

-- Criando a tabela `products`
CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `code_value` varchar(50) DEFAULT NULL,
  `is_published` varchar(50) DEFAULT NULL,
  `expiration` date DEFAULT NULL,
  `price` decimal(5,2) DEFAULT NULL,
  `id_warehouse` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id_warehouse`) REFERENCES `warehouses`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Inserindo dados na tabela `products`
INSERT INTO `products` (`name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `id_warehouse`) VALUES
('Corn Shoots', 244, '0009-1111', '0', '2022-01-08', 23.27, 1),
('Shrimp - Baby, Cold Water', 174, '49288-0877', '0', '2022-08-04', 52.12, 1),
('Sprouts - Onion', 136, '0268-6518', '1', '2021-12-27', 91.95, 1),
('Triple Sec - Mcguinness', 107, '13537-457', '0', '2021-06-23', 72.60, 1),
('Chervil - Fresh', 81, '49430-046', '0', '2022-05-28', 34.46, 1),
('Wine - German Riesling', 212, '24385-804', '0', '2022-07-07', 5.19, 1),
('Oil - Sunflower', 169, '59779-590', '0', '2022-07-03', 7.24, 1),
('Persimmons', 238, '45802-327', '0', '2021-04-14', 60.65, 1),
('Beer - Labatt Blue', 23, '48951-1215', '1', '2022-06-23', 32.99, 1),
('Ranchero - Primerba, Paste', 59, '0268-1173', '1', '2021-04-02', 17.02, 1);