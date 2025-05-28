SET client_min_messages TO WARNING;

ALTER TABLE reviews DROP COLUMN employment;
ALTER TABLE reviews DROP COLUMN employment_period;
ALTER TABLE reviews DROP COLUMN city;