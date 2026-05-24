ALTER TABLE boxes RENAME COLUMN content TO description;
ALTER TABLE boxes ALTER COLUMN description TYPE TEXT USING '';
ALTER TABLE boxes ALTER COLUMN description SET DEFAULT '';
