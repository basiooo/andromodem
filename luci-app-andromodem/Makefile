include $(TOPDIR)/rules.mk

PKG_NAME:=luci-app-andromodem
PKG_VERSION:=3.1.1
PKG_RELEASE:=1

PKG_MAINTAINER:=BobbyUnknown <bobbyun.known88@gmail.com>

LUCI_TITLE:=LuCI support for Andromodem
LUCI_DEPENDS:=+andromodem
LUCI_PKGARCH:=all

PKG_LICENSE:=MIT

include $(INCLUDE_DIR)/package.mk

define Package/$(PKG_NAME)
  CATEGORY:=LuCI
  SUBMENU:=3. Applications
  TITLE:=$(LUCI_TITLE)
  PKGARCH:=$(LUCI_PKGARCH)
  DEPENDS:=$(LUCI_DEPENDS)
  MAINTAINER:=$(PKG_MAINTAINER)
endef

define Package/$(PKG_NAME)/description
    LuCI web interface for Andromodem. Provides Android modem management through ADB.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
endef

define Build/Compile
endef

define Package/$(PKG_NAME)/install
	# 1. Install Init Script
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./root/etc/init.d/andromodem $(1)/etc/init.d/andromodem

	# 2. Install LuCI Controller
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/controller
	$(INSTALL_DATA) ./luasrc/controller/andromodem.lua $(1)/usr/lib/lua/luci/controller/andromodem.lua

	# 3. Install LuCI View
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/view
	$(INSTALL_DATA) ./luasrc/view/andromodem.htm $(1)/usr/lib/lua/luci/view/andromodem.htm
endef

define Package/$(PKG_NAME)/postinst
#!/bin/sh
[ -n "$$IPKG_INSTROOT" ] || {
	/etc/init.d/andromodem enable
	/etc/init.d/andromodem start
}
endef

define Package/$(PKG_NAME)/prerm
#!/bin/sh
[ -n "$$IPKG_INSTROOT" ] || {
	/etc/init.d/andromodem stop
	/etc/init.d/andromodem disable
}
endef

$(eval $(call BuildPackage,$(PKG_NAME)))
