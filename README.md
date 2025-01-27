# IoT Dashboard 🚧
Go web server for running a home automation dashboard on a local network. (In development)
Intended use case is an ARM based RapsberryPi.

## Technologies
- Go
- Vue.js
- TailwindCSS

## Prerequisites

- Raspberry Pi
- Go programming language installed
- Git installed

## Setup Steps

### 1. Clone the Repository

Clone your Go application repository to your Raspberry Pi:

```sh
git clone https://github.com/leon-luna-ray/iot_dashboard.git
cd iot_dashboard/server
```

### 2. Initialize Go Modules (if needed)

If your project does not already have a `go.mod` file, initialize Go modules:

```sh
go mod init github.com/leon-luna-ray/iot_dashboard
```

### 3. Build the Go Application

Build the Go application and place the binary in the `Apps` folder:

```sh
go build -o ~/Apps/iot_dashboard ./cmd
```

### 4. Create a Systemd Service

Create a systemd service file to run the application on boot:

```sh
sudo nano /etc/systemd/system/iot_dashboard.service

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

### 5. Reload Systemd and Enable the Service

Reload systemd to recognize the new service:

```sh
sudo systemctl daemon-reload
sudo systemctl enable iot_dashboard.service
```

### 6. Start the Service

Start the service:

```sh
sudo systemctl start iot_dashboard.service
sudo systemctl status iot_dashboard.service
```

### 7. Managing the Service

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

