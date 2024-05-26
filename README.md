# Go Infra Inventory

Go Infra Inventory is a Go application designed to collect and report infrastructure details from various clients (VMs). It uses JWT authentication to securely communicate with the server and stores configuration details in environment variables.

## Features

- Collects details like IP, hostname, services, OS version, and open ports.
- Uses JWT for secure authentication.
- Automatically renews tokens using a refresh token.
- Reports inventory details to the server every 15 minutes.
- Stores tokens securely in a temporary file.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose (for setting up the MySQL server)
- A running MySQL server

### Setting Up the Server

1. **Create the `.env` file for the server:**

   ```env
   MYSQL_USER=user
   MYSQL_PASSWORD=password
   MYSQL_DATABASE=inventory
   MYSQL_PORT=3306
   SERVER_PORT=8080
   INVENTORY_TABLE=inventory
   JWT_SECRET=your-256-bit-secret
   JWT_REFRESH_SECRET=your-refresh-256-bit-secret
   ```

2. **Create the `docker-compose.yml` file:**

   ```yaml
   version: "3.8"

   services:
     db:
       image: mysql:8.0
       container_name: inventory_db
       environment:
         MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
         MYSQL_DATABASE: ${MYSQL_DATABASE}
         MYSQL_USER: ${MYSQL_USER}
         MYSQL_PASSWORD: ${MYSQL_PASSWORD}
       ports:
         - "${MYSQL_PORT}:3306"
       volumes:
         - db_data:/var/lib/mysql
         - ./init.sql:/docker-entrypoint-initdb.d/init.sql
       networks:
         - inventory_network

   volumes:
     db_data:

   networks:
     inventory_network:
   ```

3. **Create the `init.sql` file to initialize the database:**

   ```sql
   CREATE DATABASE IF NOT EXISTS inventory;

   USE inventory;

   CREATE TABLE IF NOT EXISTS inventory (
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
   ```

4. **Start the MySQL service:**

   ```sh
   docker-compose up -d
   ```

5. **Run the server:**

   ```sh
   go run cmd/server/main.go
   ```

### Setting Up the Client

Server needs to generate a Inital JWT for the client.

```sh
tools/GenerateJWT.go
```

1. **Create the `.env` file for the client:**

   ```env
   SERVER_URL=http://your-server-ip:8080
   INITIAL_ACCESS_TOKEN=your-initial-jwt-token
   ```

2. **Run the client:**

   ```sh
   go run cmd/client/main.go
   ```

### Compiling for Different Operating Systems

#### Compiling the Server

For Linux:

```sh
GOOS=linux GOARCH=amd64 go build -o server-linux server/cmd/server/main.go
```

For macOS:

```sh
GOOS=darwin GOARCH=amd64 go build -o server-mac server/cmd/server/main.go
```

For Windows:

```sh
GOOS=windows GOARCH=amd64 go build -o server-windows.exe server/cmd/server/main.go
```

#### Compiling the Client

For Linux:

```sh
GOOS=linux GOARCH=amd64 go build -o client-linux client/cmd/client/main.go
```

For macOS:

```sh
GOOS=darwin GOARCH=amd64 go build -o client-mac client/cmd/client/main.go
```

For Windows:

```sh
GOOS=windows GOARCH=amd64 go build -o client-windows.exe client/cmd/client/main.go
```

## Usage

- Server: Runs the API that collects and stores inventory data.
- Client: Collects data from the VM and reports it to the server.

## Inventory Data Collected

- IP Address
- Hostname
- Services (e.g., SQL Server, IIS, Nginx, Docker)
- OS Version
- Open Ports
- Last Reported Date
- Next Report Date
- Health Status

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### License

This project is licensed under the MIT License.
