BEGIN;
ALTER TABLE tokenpool ADD COLUMN state VARCHAR(64);
UPDATE tokenpool SET state='confirmed';
ALTER TABLE tokenpool ALTER COLUMN state SET NOT NULL;
COMMIT;
