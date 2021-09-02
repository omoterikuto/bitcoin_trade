USE brc_trade;
CREATE TABLE IF NOT EXISTS signal_events (time DATETIME PRIMARY KEY NOT NULL, product_code STRING, side STRING, price FLOAT, size FLOAT);