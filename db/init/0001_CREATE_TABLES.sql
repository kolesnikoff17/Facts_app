CREATE TABLE Facts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE Links (
    id SERIAL PRIMARY KEY,
    fact_id INT NOT NULL,
    link VARCHAR(255) NOT NULL,
    FOREIGN KEY (fact_id) REFERENCES Facts (id)
);