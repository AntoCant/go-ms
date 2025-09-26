CREATE TABLE products (
    id    TEXT PRIMARY KEY,
    name  TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL CHECK (price > 0),
    stock INT NOT NULL CHECK (stock >= 0)
);
