CREATE DATABASE IF NOT EXISTS inventory;

USE inventory;

CREATE TABLE
    IF NOT EXISTS inventory (
        id INT AUTO_INCREMENT PRIMARY KEY,
        ip VARCHAR(15),
        hostname VARCHAR(255),
        services TEXT,
        os_version VARCHAR(255),
        open_ports TEXT,
        last_reported_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        next_report_date TIMESTAMP,
        health VARCHAR(50),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );