include $(TOPDIR)/rules.mk

PKG_NAME:=andromodem
PKG_VERSION:=3.1.1
PKG_RELEASE:=1

PKG_MAINTAINER:=BobbyUnknown <bobbyun.known88@gmail.com>
PKG_LICENSE:=MIT

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/basiooo/andromodem.git
PKG_SOURCE_VERSION:=main
PKG_HASH:=skip
PKG_SOURCE_SUBDIR:=$(PKG_NAME)-$(PKG_VERSION)
PKG_SOURCE:=$(PKG_NAME)-$(PKG_VERSION).tar.gz

PKG_BUILD_DEPENDS:=golang/host

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

GO_PKG:=github.com/basiooo/andromodem

define Package/andromodem
  SECTION:=net
  CATEGORY:=Network
  TITLE:=AndroModem - Android Modem Manager
  URL:=https://github.com/basiooo/andromodem
  DEPENDS:=+adb
endef

define Package/andromodem/description
  AndroModem - Manage Android devices as modems via ADB.
  Provides web interface for monitoring and controlling Android modem connections.
endef

define Build/Compile
	@echo "Compiling AndroModem Go binary..."
	cd $(PKG_BUILD_DIR) && \
		GOOS=linux \
		GOARCH=$(GO_ARCH) \
		$(if $(findstring arm,$(GO_ARCH)),GOARM=$(or $(GO_ARM),7)) \
		$(if $(findstring mips,$(GO_ARCH)),GOMIPS=$(or $(GO_MIPS),softfloat)) \
		CGO_ENABLED=0 \
		go build -v -trimpath -ldflags="-s -w -X main.version=$(PKG_VERSION)" \
		-o $(PKG_BUILD_DIR)/andromodem ./cmd/andromodem
endef

define Package/andromodem/install
	$(INSTALL_DIR) $(1)/usr/share/andromodem
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/andromodem $(1)/usr/share/andromodem/andromodem
endef

$(eval $(call BuildPackage,andromodem))
