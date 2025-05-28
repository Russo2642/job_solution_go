SET client_min_messages TO WARNING;

ALTER TABLE industries ADD COLUMN IF NOT EXISTS color VARCHAR(7);

UPDATE industries SET color = '#4285F4' WHERE id = 1; -- Разработка ПО (синий)
UPDATE industries SET color = '#34A853' WHERE id = 2; -- Веб-разработка (зеленый)
UPDATE industries SET color = '#FBBC05' WHERE id = 3; -- Мобильная разработка (желтый)
UPDATE industries SET color = '#EA4335' WHERE id = 4; -- Искусственный интеллект (красный)
UPDATE industries SET color = '#673AB7' WHERE id = 5; -- Большие данные (фиолетовый)
UPDATE industries SET color = '#E91E63' WHERE id = 6; -- Кибербезопасность (розовый)
UPDATE industries SET color = '#03A9F4' WHERE id = 7; -- Облачные вычисления (голубой)
UPDATE industries SET color = '#009688' WHERE id = 8; -- DevOps и SRE (бирюзовый)
UPDATE industries SET color = '#8BC34A' WHERE id = 9; -- Финтех (светло-зеленый)
UPDATE industries SET color = '#FF9800' WHERE id = 10; -- Блокчейн (оранжевый)
UPDATE industries SET color = '#9C27B0' WHERE id = 11; -- Игровая индустрия (пурпурный)
UPDATE industries SET color = '#00BCD4' WHERE id = 12; -- E-commerce (цвет морской волны)
UPDATE industries SET color = '#3F51B5' WHERE id = 13; -- Интернет вещей (индиго)
UPDATE industries SET color = '#607D8B' WHERE id = 14; -- Телекоммуникации (серо-синий)
UPDATE industries SET color = '#795548' WHERE id = 15; -- ERP/CRM системы (коричневый)