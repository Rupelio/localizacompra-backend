ALTER TABLE products
ADD COLUMN category_id BIGINT,
ADD CONSTRAINT fk_category
FOREIGN KEY (category_id) REFERENCES categories(id);
