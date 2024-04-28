#!/bin/sh

opkg update
echo "Installing Adb, Wget and Curl"
opkg install adb
opkg install wget
opkg install curl


if [ ! -d "/opt" ]; then
    echo "Creating directory /opt..."
    mkdir -p /opt
fi

echo "Kill existing andromodem proccess."
killall /opt/andromodem

arch=$(uname -m)

if [[ $arch == *"arm"* ]]; then
    download_url="https://andromodem.bagasjulianto.my.id/download/latest/andromodem_arm"
elif [[ $arch == *"x86_64"* || $arch == *"amd64"* ]]; then
    download_url="https://andromodem.bagasjulianto.my.id/download/latest/andromodem_amd64"
elif [[ $arch == *"x86"* ]]; then
    download_url="https://andromodem.bagasjulianto.my.id/download/latest/andromodem_x86"
elif [[ $arch == *"aarch64"* || $arch == *"arm64"* ]]; then
    download_url="https://andromodem.bagasjulianto.my.id/download/latest/andromodem_arm64"
else
    echo "Unsupported architecture: $arch"
    exit 1
fi

echo "Downloading Andromodem for $arch device..."
if curl -o /opt/andromodem "$download_url"; then
    chmod +x /opt/andromodem
else
    echo "Failed download using curl. Trying with wget..."
    if wget -O /opt/andromodem "$download_url"; then
        chmod +x /opt/andromodem
    else
        echo "Both curl and wget failed. Please check your internet connection and try again."
        exit 1
    fi
fi

echo "Add Andromodem in luci menu"
mkdir -p /usr/lib/lua/luci/view

cat << EOF > /usr/lib/lua/luci/view/andromodem.htm
<%+header%>
<div class="cbi-map"><br>
<iframe id="andromodem" style="width: 100%; min-height: 750px; border: none;"></iframe>
</div>
<script type="text/javascript">
document.getElementById("andromodem").src = \`http://\${window.location.hostname}:49153\`;
</script>
<%+footer%>
EOF

mkdir -p /usr/lib/lua/luci/controller

cat << EOF > /usr/lib/lua/luci/controller/andromodem.lua
module("luci.controller.andromodem", package.seeall)
function index()
entry({"admin","status","andromodem"}, template("andromodem"), _("Andromodem"), 4).leaf=true
end
EOF

cat <<EOF > /etc/init.d/andromodem
#!/bin/sh /etc/rc.common
# Copyright (C) 2024 Bagas Julianto

USE_PROCD=1
START=99

stop_service() {
    killall andromodem
}

start_service() {
    procd_open_instance
    procd_set_param command "/opt/andromodem"
    procd_close_instance
}
EOF


# Remove fucking startup code
if grep -qxF '/opt/andromodem' /etc/rc.local; then
    sed -i '/\/opt\/andromodem/d' /etc/rc.local
fi

echo "Starting Andromodem"
chmod 775 /etc/init.d/andromodem
/etc/init.d/andromodem start
