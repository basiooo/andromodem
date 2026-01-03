# Makefile - AndroModem Debian/Arch Packaging

PKG_NAME:=andromodem
PKG_VERSION:=3.1.1
PKG_RELEASE:=1
PKG_MAINTAINER:=BobbyUnknown
PKG_DESC:=AndroModem - Android Modem Manager
PKG_ARCH_AMD64:=amd64
PKG_ARCH_ARM64:=arm64
PKG_ARCH_ARMHF:=armhf

# Arch Linux Architecture Mappings
PKG_ARCH_ARCH_X86_64:=x86_64
PKG_ARCH_ARCH_AARCH64:=aarch64
PKG_ARCH_ARCH_ARMV7H:=armv7h

BUILD_DIR:=build
CORE_DIR:=core

# Binary mappings
BIN_AMD64:=andromodem-linux-amd64
BIN_ARM64:=andromodem-linux-arm64
BIN_ARMHF:=andromodem-linux-armv7

.PHONY: all clean build-core build-binaries build-frontend deb-amd64 deb-arm64 deb-armhf arch-x86_64 arch-aarch64 arch-armv7h build-all build-deb build-arch dev

all: build-frontend build-binaries build-deb build-arch

# Frontend build target
build-frontend:
	@echo "Building frontend..."
	@if [ ! -d "templates/andromodem_fe/node_modules" ]; then \
		echo "Installing frontend dependencies..."; \
		cd templates/andromodem_fe && npm install; \
	fi
	@echo "Building React frontend..."
	cd templates/andromodem_fe && npm run build
	@echo "Frontend built successfully!"

# Development target - run with Air for hot-reload
dev:
	@echo "Starting AndroModem in development mode with Air..."
	@command -v air > /dev/null 2>&1 || { echo "Air is not installed. Installing..."; go install github.com/air-verse/air@latest; }
	air

build-deb: deb-amd64 deb-arm64 deb-armhf

build-arch: arch-x86_64 arch-aarch64 arch-armv7h

build-all: all

build-binaries: build-frontend
	@echo "Building standalone binaries for multiple platforms..."
	mkdir -p bin
	
	@echo "Building for linux/amd64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-amd64 ./cmd/andromodem
	
	@echo "Building for linux/arm64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-arm64 ./cmd/andromodem
	
	@echo "Building for linux/386..."
	GOOS=linux GOARCH=386 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-386 ./cmd/andromodem
	
	@echo "Building for linux/arm (v7)..."
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-armv7 ./cmd/andromodem
	
	@echo "Building for linux/mips..."
	GOOS=linux GOARCH=mips GOMIPS=softfloat CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-mips ./cmd/andromodem
	
	@echo "Building for linux/mipsle..."
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-mipsle ./cmd/andromodem
	
	@echo "Building for linux/mips64..."
	GOOS=linux GOARCH=mips64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-mips64 ./cmd/andromodem
	
	@echo "Building for linux/mips64le..."
	GOOS=linux GOARCH=mips64le CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-mips64le ./cmd/andromodem
	
	@echo "Building for linux/riscv64..."
	GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-riscv64 ./cmd/andromodem
	
	@echo "Building for linux/ppc64..."
	GOOS=linux GOARCH=ppc64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-ppc64 ./cmd/andromodem
	
	@echo "Building for linux/ppc64le..."
	GOOS=linux GOARCH=ppc64le CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-ppc64le ./cmd/andromodem
	
	@echo "Building for linux/s390x..."
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-linux-s390x ./cmd/andromodem
	
	@echo "Building for android/arm64..."
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-android-arm64 ./cmd/andromodem
	
	@echo "Building for darwin/amd64 (macOS Intel)..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-darwin-amd64 ./cmd/andromodem
	
	@echo "Building for darwin/arm64 (macOS Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o bin/andromodem-darwin-arm64 ./cmd/andromodem
	
	@echo "All binaries built successfully in bin/"

build-core: build-frontend
	@echo "Building AndroModem binaries from source..."
	mkdir -p $(CORE_DIR)
	
	@echo "Building for AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o $(CORE_DIR)/$(BIN_AMD64) ./cmd/andromodem
	
	@echo "Building for ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o $(CORE_DIR)/$(BIN_ARM64) ./cmd/andromodem
	
	@echo "Building for ARMv7..."
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.Version=$(PKG_VERSION)" \
		-o $(CORE_DIR)/$(BIN_ARMHF) ./cmd/andromodem
	
	@echo "Core binaries ready in $(CORE_DIR)/"

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(CORE_DIR)
	rm -rf bin
	rm -rf templates/andromodem_fe/dist/*
	@echo "Cleaned build artifacts and frontend dist"

# Template for building Debian package
# Usage: $(call build_deb, ARCH, ANDROMODEM_BINARY)
define build_deb
	@echo "Building Debian package for $(1)..."
	$(eval PKG_DIR := $(BUILD_DIR)/$(PKG_NAME)_$(PKG_VERSION)-$(PKG_RELEASE)_$(1))
	rm -rf $(PKG_DIR)
	mkdir -p $(PKG_DIR)/usr/share/andromodem
	mkdir -p $(PKG_DIR)/lib/systemd/system
	mkdir -p $(PKG_DIR)/DEBIAN

	# Copy AndroModem Binary
	@if [ -f $(CORE_DIR)/$(2) ]; then \
		cp $(CORE_DIR)/$(2) $(PKG_DIR)/usr/share/andromodem/andromodem; \
		chmod +x $(PKG_DIR)/usr/share/andromodem/andromodem; \
	else \
		echo "Error: $(CORE_DIR)/$(2) not found!"; \
		exit 1; \
	fi
	
	# Create Control File
	echo "Package: $(PKG_NAME)" > $(PKG_DIR)/DEBIAN/control
	echo "Version: $(PKG_VERSION)-$(PKG_RELEASE)" >> $(PKG_DIR)/DEBIAN/control
	echo "Section: net" >> $(PKG_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(PKG_DIR)/DEBIAN/control
	echo "Architecture: $(1)" >> $(PKG_DIR)/DEBIAN/control
	echo "Maintainer: $(PKG_MAINTAINER)" >> $(PKG_DIR)/DEBIAN/control
	echo "Depends: adb, ca-certificates" >> $(PKG_DIR)/DEBIAN/control
	echo "Description: $(PKG_DESC)" >> $(PKG_DIR)/DEBIAN/control
	echo "  Manage Android devices as modems via ADB." >> $(PKG_DIR)/DEBIAN/control
	echo "  Provides web interface for monitoring and controlling Android modem connections." >> $(PKG_DIR)/DEBIAN/control

	# Create Systemd Service
	echo "[Unit]" > $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "Description=AndroModem Service" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "After=network.target" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "[Service]" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "ExecStart=/usr/share/andromodem/andromodem" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "Restart=always" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "RestartSec=5" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "User=root" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "[Install]" >> $(PKG_DIR)/lib/systemd/system/andromodem.service
	echo "WantedBy=multi-user.target" >> $(PKG_DIR)/lib/systemd/system/andromodem.service

	# Create postinst script
	echo "#!/bin/sh" > $(PKG_DIR)/DEBIAN/postinst
	echo "set -e" >> $(PKG_DIR)/DEBIAN/postinst
	echo "if [ \"\$$1\" = \"configure\" ]; then" >> $(PKG_DIR)/DEBIAN/postinst
	echo "    systemctl daemon-reload" >> $(PKG_DIR)/DEBIAN/postinst
	echo "    systemctl enable andromodem" >> $(PKG_DIR)/DEBIAN/postinst
	echo "    systemctl restart andromodem || true" >> $(PKG_DIR)/DEBIAN/postinst
	echo "fi" >> $(PKG_DIR)/DEBIAN/postinst
	chmod 755 $(PKG_DIR)/DEBIAN/postinst

	# Create prerm script
	echo "#!/bin/sh" > $(PKG_DIR)/DEBIAN/prerm
	echo "set -e" >> $(PKG_DIR)/DEBIAN/prerm
	echo "if [ \"\$$1\" = \"remove\" ]; then" >> $(PKG_DIR)/DEBIAN/prerm
	echo "    systemctl stop andromodem || true" >> $(PKG_DIR)/DEBIAN/prerm
	echo "    systemctl disable andromodem || true" >> $(PKG_DIR)/DEBIAN/prerm
	echo "fi" >> $(PKG_DIR)/DEBIAN/prerm
	chmod 755 $(PKG_DIR)/DEBIAN/prerm

	# Build Package
	dpkg-deb --build $(PKG_DIR)
	@echo "Package created at $(PKG_DIR).deb"
endef

deb-amd64: build-core
	$(call build_deb,$(PKG_ARCH_AMD64),$(BIN_AMD64))

deb-arm64: build-core
	$(call build_deb,$(PKG_ARCH_ARM64),$(BIN_ARM64))

deb-armhf: build-core
	$(call build_deb,$(PKG_ARCH_ARMHF),$(BIN_ARMHF))

# Template for building Arch Linux package
# Usage: $(call build_arch, ARCH, ANDROMODEM_BINARY)
define build_arch
	@echo "Building Arch package for $(1)..."
	$(eval PKG_DIR := $(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)-$(PKG_RELEASE)-$(1))
	rm -rf $(PKG_DIR)
	mkdir -p $(PKG_DIR)/usr/share/andromodem
	mkdir -p $(PKG_DIR)/usr/lib/systemd/system

	# Copy AndroModem Binary
	@if [ -f $(CORE_DIR)/$(2) ]; then \
		cp $(CORE_DIR)/$(2) $(PKG_DIR)/usr/share/andromodem/andromodem; \
		chmod +x $(PKG_DIR)/usr/share/andromodem/andromodem; \
	else \
		echo "Error: $(CORE_DIR)/$(2) not found!"; \
		exit 1; \
	fi

	# Create Systemd Service
	echo "[Unit]" > $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "Description=AndroModem Service" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "After=network.target" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "[Service]" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "ExecStart=/usr/share/andromodem/andromodem" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "Restart=always" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "RestartSec=5" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "User=root" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "[Install]" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service
	echo "WantedBy=multi-user.target" >> $(PKG_DIR)/usr/lib/systemd/system/andromodem.service

	# Create .PKGINFO
	echo "pkgname = $(PKG_NAME)" > $(PKG_DIR)/.PKGINFO
	echo "pkgver = $(PKG_VERSION)-$(PKG_RELEASE)" >> $(PKG_DIR)/.PKGINFO
	echo "pkgdesc = $(PKG_DESC)" >> $(PKG_DIR)/.PKGINFO
	echo "url = https://github.com/basiooo/andromodem" >> $(PKG_DIR)/.PKGINFO
	echo "builddate = $$(date +%s)" >> $(PKG_DIR)/.PKGINFO
	echo "packager = $(PKG_MAINTAINER)" >> $(PKG_DIR)/.PKGINFO
	echo "size = $$(du -sb $(PKG_DIR) | cut -f1)" >> $(PKG_DIR)/.PKGINFO
	echo "arch = $(1)" >> $(PKG_DIR)/.PKGINFO
	echo "depend = android-tools" >> $(PKG_DIR)/.PKGINFO
	echo "depend = ca-certificates" >> $(PKG_DIR)/.PKGINFO

	# Create .install script
	echo "post_install() {" > $(PKG_DIR)/.INSTALL
	echo "    systemctl daemon-reload" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl enable andromodem" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl start andromodem || true" >> $(PKG_DIR)/.INSTALL
	echo "}" >> $(PKG_DIR)/.INSTALL
	echo "" >> $(PKG_DIR)/.INSTALL
	echo "post_upgrade() {" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl daemon-reload" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl restart andromodem || true" >> $(PKG_DIR)/.INSTALL
	echo "}" >> $(PKG_DIR)/.INSTALL
	echo "" >> $(PKG_DIR)/.INSTALL
	echo "pre_remove() {" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl stop andromodem || true" >> $(PKG_DIR)/.INSTALL
	echo "    systemctl disable andromodem || true" >> $(PKG_DIR)/.INSTALL
	echo "}" >> $(PKG_DIR)/.INSTALL

	# Create .MTREE (file list)
	cd $(PKG_DIR) && find . -type f -o -type l | LC_ALL=C sort | sed 's/^\.\///' > .MTREE

	# Build Package (create tar.zst archive)
	cd $(PKG_DIR) && tar -cf - .PKGINFO .INSTALL .MTREE usr | zstd -19 -T0 -q -o ../$(PKG_NAME)-$(PKG_VERSION)-$(PKG_RELEASE)-$(1).pkg.tar.zst
	@echo "Package created at $(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)-$(PKG_RELEASE)-$(1).pkg.tar.zst"
endef

arch-x86_64: build-core
	$(call build_arch,$(PKG_ARCH_ARCH_X86_64),$(BIN_AMD64))

arch-aarch64: build-core
	$(call build_arch,$(PKG_ARCH_ARCH_AARCH64),$(BIN_ARM64))

arch-armv7h: build-core
	$(call build_arch,$(PKG_ARCH_ARCH_ARMV7H),$(BIN_ARMHF))
