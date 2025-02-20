# IoT Dashboard 🚧
Go web server for running a home automation dashboard on a local network (in development).
Intended use case is an ARM based RapsberryPi.

## Technologies
- Go
- Vue.js
- TailwindCSS

## Prerequisites

- Raspberry Pi
- Go installed
- Node.js installed
- Pnpm installed
- Git installed

## Setup

### 1. Clone the Repository

Clone the repository to your Raspberry Pi and navigate to the project directory:

```sh
git clone https://github.com/leon-luna-ray/iot_dashboard.git
cd ./backend
```

### 2. Initialize Go Modules (if needed)

If project does not already have a `go.mod` file, initialize Go modules:

```sh
go mod init github.com/leon-luna-ray/iot_dashboard
```

### 3. Build the frontend
If pnpm is not already installed, you can install it using npm:

```sh
npm install -g pnpm
```

Navigate to the frontend directory and build the Vue.js application using pnpm:

```sh
cd frontend
pnpm install
pnpm run build
```

### 4. Add enviorment variables (if needed)
Add the needed enivorment variables for connecting to external APIs. Follow .env.example.

If you're connected to the RaspberryPi via SSH you can create the file inside the project directory with nano.

```sh
nano .env
```

### 5. Build the Go Application

Build the Go application and place the binary in the desired directory:

```sh
go build -o ./bin
```

### 6. Create a Systemd Service

Create a systemd service file to run the application on boot:

```sh
sudo nano /etc/systemd/system/iot_dashboard.service
```

Add service updating your relative path:

```sh
[Unit]
Description=IoT Dashboard Service
After=network.target

[Service]
ExecStart=/home/user/Apps/iot_dashboard
Restart=always
User=user
Group=user
Environment=PATH=/usr/bin:/usr/local/bin
Environment=GO_ENV=production
WorkingDirectory=/home/user/Apps

[Install]
WantedBy=multi-user.target
```

### 7. Reload Systemd and Enable the Service

Reload systemd to recognize the new service:

```sh
sudo systemctl daemon-reload
sudo systemctl enable iot_dashboard.service
```

### 8. Start the Service

Start the service:

```sh
sudo systemctl start iot_dashboard.service
sudo systemctl status iot_dashboard.service
```


## Updates

To update the application with the latest changes from the main branch, follow these steps:

### 1. Pull the Latest Changes

Navigate to the project directory and pull the latest changes from the main branch:

```sh
cd ./iot_dashboard
git pull origin main
```
### 2. Build the frontend frontend

Navigate to the frontend directory and rebuild the Vue.js application using pnpm:

```sh
cd frontend
pnpm install
pnpm run build
```

### 3. Build the Go Application

Rebuild the Go application and place the binary in the desired directory:

```sh
cd ../backend
go mod tidy
go build -o ./bin
```

### 4. Restart the Service

Restart the systemd service to apply the updates:

```sh
sudo systemctl restart iot_dashboard.service
```

### 5. Check the Service Status

Ensure the service is running correctly:

```sh
sudo systemctl status iot_dashboard.service
```

## Management Commands

Use the following commands to manage the service:

- **Check the status of the service**:
  ```sh
  sudo systemctl status iot_dashboard.service
  ```
- **Start the service**:
  ```sh
  sudo systemctl start iot_dashboard.service
  ```

- **Stop the service**:
  ```sh
  sudo systemctl stop iot_dashboard.service
  ```
- **Restart the service**:
  ```sh
  sudo systemctl restart iot_dashboard.service
  ```

- **Enable the service to start on boot**:
  ```sh
  sudo systemctl enable iot_dashboard.service
  ```

- **Disable the service from starting on boot**:
    ```sh
    sudo systemctl disable iot_dashboard.service
    ```
