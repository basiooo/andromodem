# 🚀 AndroModem

<div align="center">

![Codecov](https://codecov.io/gh/basiooo/andromodem/branch/main/graph/badge.svg)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/basiooo/andromodem)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/basiooo/andromodem/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/basiooo/andromodem)
![Go Version](https://img.shields.io/github/go-mod/go-version/basiooo/andromodem)
![License](https://img.shields.io/github/license/basiooo/andromodem)

**Android device management system with monitoring and automation capabilities**

</div>

---

## 📖 Overview
AndroModem is a comprehensive Android device management solution. Built with a Go backend and React frontend, it provides real-time monitoring, automated failover mechanisms, and an intuitive web-based control panel for Android devices connected via ADB (Android Debug Bridge).

> **🛜 Designed for OpenWrt Users:**
AndroModem is primarily designed for OpenWrt users who use Android phones as USB modems. It enables ADB-based monitoring and control directly from the router.



## ✨ Features

### 🔧 **Device Management**
- **Real-time Device Detection**: Automatically discover and monitor connected Android devices
- **Comprehensive Device Info**: View detailed specifications, battery status, memory, storage, and system information
- **Power Control**: Remote power operations (reboot, power off, recovery mode, bootloader)
- **Root Detection**: Automatic detection of device root status
- **Multi-device Support**: Manage multiple Android devices simultaneously

### 🌐 **Network Management**
- **Network Information**: View IP routes, APN settings, and SIM card details
- **Mobile Data Control**: Toggle mobile data on/off
- **Airplane Mode Control**: Toggle airplane mode with legacy Android support
- **Signal Strength Monitoring**: Signal strength and network type detection

### 📊 **Advanced Monitoring System**
- **Monitoring Methods**: HTTP, HTTPS, ICMP ping monitoring
- **Automated Failover**: Automatic airplane mode toggle on connection failures
- **Configurable Thresholds**: Custom failure limits and check intervals
- **Real-time Logging**: Live monitoring logs

### 🌐 **Modern Web Interface**
- **Responsive Design**: Clean, modern UI built with React 19 and Tailwind CSS
- **Real-time Device Updates**: Live device status via Server-Sent Events (SSE)
- **Device Selector**: Intuitive switching between multiple devices
- **Theme Support**: Multiple UI themes with DaisyUI components


## 🚀 Installation

### Prerequisites
- **Go 1.22.2+** - [Download Go](https://golang.org/dl/)
- **Node.js 18+** and **npm** - [Download Node.js](https://nodejs.org/)
- **Android Debug Bridge (ADB)** - [Install ADB](https://developer.android.com/studio/command-line/adb)
- **Android devices** with USB debugging enabled

### 📦 Quick Start

#### Option 1: Download Pre-built Binary
```bash
# Download the latest release for your platform
wget https://github.com/basiooo/andromodem/releases/latest/download/andromodem_linux_amd64
chmod +x andromodem_linux_amd64
./andromodem_linux_amd64
```

#### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/basiooo/andromodem.git
cd andromodem

# go frontend directory
cd templates/andromodem_fe
# install frontend dependencies
npm install
# build frontend
npm run build

# back to root
cd ../../

# Build the application
go build -o andromodem cmd/andromodem/main.go

# Run the application
./andromodem
```

#### Option 3: OpenWrt Installation
```bash
# One-line installation on OpenWrt
wget -O - https://raw.githubusercontent.com/basiooo/andromodem/main/andromodem_openwrt.sh | sh -s install
```

---

## 🎮 Usage

### Starting the Application
```bash
# Start AndroModem
./andromodem

# Check version
./andromodem --version
```

### Web Interface
Once started, access the web interface at:
- **Local**: http://localhost:49153
- **Network**: http://YOUR_IP:49153

### OpenWrt Integration
After installation on OpenWrt, access via:
- **LuCI Interface**: Status → AndroModem
- **Direct Access**: http://router-ip:49153

---

## 🔧 Development

### Backend Development
```bash
# Install dependencies
go mod download

# Run with hot reload
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...
```

### Frontend Development
```bash
# Navigate to frontend directory
cd templates/andromodem_fe

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

---

## 🤝 Contributing

We welcome contributions!

### Development Setup
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [goadb](https://github.com/basiooo/goadb) - ADB client library
- [Chi](https://github.com/go-chi/chi) - HTTP router
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [React](https://reactjs.org/) - Frontend framework
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [DaisyUI](https://daisyui.com/) - UI components

---

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/basiooo/andromodem/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/basiooo/andromodem/discussions)
- 📧 **Contact**: [Bagas Julianto](https://github.com/basiooo)

---

<div align="center">

**Made with ❤️ by [Bagas Julianto](https://github.com/basiooo)**

⭐ Star this repository if you find it helpful!

</div>
        
