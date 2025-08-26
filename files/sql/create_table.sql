CREATE TABLE orders
(
    id       INT AUTO_INCREMENT PRIMARY KEY,
    user_id  INT         NOT NULL,
    quantity INT         NOT NULL,
    total DOUBLE NOT NULL,
    status   VARCHAR(50) NOT NULL,
    total_mark_up DOUBLE NOT NULL,
    total_discount DOUBLE NOT NULL
);

CREATE TABLE product_requests
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    order_id   INT NOT NULL REFERENCES orders (id),
    product_id INT NOT NULL,
    quantity   INT NOT NULL,
    mark_up DOUBLE NOT NULL,
    discount DOUBLE NOT NULL,
    final_price DOUBLE NOT NULL
);