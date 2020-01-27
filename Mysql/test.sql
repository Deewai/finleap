Use test_db;

CREATE TABLE IF NOT EXISTS cities
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL
);
CREATE TABLE IF NOT EXISTS temperatures
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT NOT NULL,
    max INT NOT NULL,
    min INT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    FOREIGN KEY (city_id) REFERENCES cities(id)
);

CREATE TABLE IF NOT EXISTS webhooks
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT NOT NULL,
    callback_url TEXT NOT NULL,
    FOREIGN KEY (city_id) REFERENCES cities(id)
);