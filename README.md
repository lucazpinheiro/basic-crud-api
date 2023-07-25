CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  created_at timestamp default NOW(),
  name VARCHAR(250) NOT NULL,
  isAvailable boolean NOT NULL,
  price float NOT NULL,
  description VARCHAR(250)
);

SELECT * FROM products

INSERT INTO products (name, isAvailable, price, description) VALUES ('a', true, 34, 'aaaa');
INSERT INTO products (name, isAvailable, price, description) VALUES ('b', true, 34, 'bbbb');
INSERT INTO products (name, isAvailable, price, description) VALUES ('c', true, 54.6, 'cccc');
INSERT INTO products (name, isAvailable, price ) VALUES ('d', true, 41);
INSERT INTO products (name, isAvailable, price, description) VALUES ('e', true, 901, '');

------------
curl request to create a product
```
curl -d '{"name": "iphone branco", "price": 49, "available": true}' -H "Content-Type: application/json" -X POST http://localhost:3000/products
```