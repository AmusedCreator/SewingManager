-- Бэкап базы данных SewingDB
-- Время: 20241015_122839

-- Данные таблицы Nomenclature
INSERT INTO Nomenclature VALUES ('1', 'Shirt', '15.50');
INSERT INTO Nomenclature VALUES ('2', 'Pants', '25.00');
INSERT INTO Nomenclature VALUES ('3', 'Dress', '40.75');
INSERT INTO Nomenclature VALUES ('4', 'Jacket', '60.00');
INSERT INTO Nomenclature VALUES ('5', 'Skirt', '20.30');

-- Данные таблицы Tasks
INSERT INTO Tasks VALUES ('1', '1001', '01-09-2024', '10-09-2024', 'Client A', '1', '10', '10', '10 shirts for Client A');
INSERT INTO Tasks VALUES ('2', '1002', '02-09-2024', '15-09-2024', 'Client B', '3', '5', '5', '5 dresses for Client B');
INSERT INTO Tasks VALUES ('3', '1003', '03-09-2024', '20-09-2024', 'Client C', '2', '7', '7', '7 pants for Client C');
INSERT INTO Tasks VALUES ('4', '1004', '04-09-2024', '25-09-2024', 'Client D', '4', '3', '3', '3 jackets for Client D');
INSERT INTO Tasks VALUES ('5', '1005', '05-09-2024', '30-09-2024', 'Client E', '5', '8', '8', '8 skirts for Client E');

-- Данные таблицы Workers
INSERT INTO Workers VALUES ('1', 'John', 'Doe', 'Experienced tailor specializing in shirts');
INSERT INTO Workers VALUES ('2', 'Jane', 'Smith', 'Expert in dressmaking');
INSERT INTO Workers VALUES ('3', 'Emily', 'Brown', 'Works with pants and jackets');
INSERT INTO Workers VALUES ('4', 'Michael', 'Johnson', 'Skirt and jacket specialist');
INSERT INTO Workers VALUES ('5', 'Sara', 'Miller', 'Apprentice, learning all types of sewing');

-- Данные таблицы Task_Workers
INSERT INTO Task_Workers VALUES ('1', '1', '1', '5', '01-09-2024', '77.50');
INSERT INTO Task_Workers VALUES ('2', '1', '1', '5', '02-09-2024', '77.50');
INSERT INTO Task_Workers VALUES ('3', '2', '2', '3', '03-09-2024', '122.25');
INSERT INTO Task_Workers VALUES ('4', '2', '2', '2', '04-09-2024', '81.50');
INSERT INTO Task_Workers VALUES ('5', '3', '3', '4', '05-09-2024', '100.00');
INSERT INTO Task_Workers VALUES ('6', '3', '3', '3', '06-09-2024', '75.00');
INSERT INTO Task_Workers VALUES ('7', '4', '4', '3', '07-09-2024', '180.00');
INSERT INTO Task_Workers VALUES ('8', '5', '5', '8', '08-09-2024', '162.40');

