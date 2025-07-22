CREATE TABLE tasks (
    id INT AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    complite BOOLEAN DEFAULT false,
    date_start DATETIME DEFAULT now(),
    date_end DATETIME DEFAULT now(),
    date_insert DATETIME DEFAULT now(),
    PRIMARY KEY (id)
)