#!/bin/sh
# Bagas Juliato @basiooo <bagasjulianto.my.id>

if ! command -v adb >/dev/null 2>&1 ; then
    echo "ADB belum terinstall. Memulai proses instalasi..."

    # Install ADB
    opkg update
    opkg install android-tools-adb

    # Cek apakah instalasi berhasil
    if [ $? -eq 0 ]; then
        echo "Instalasi ADB berhasil."
    else
        echo "Gagal menginstal ADB. Pastikan Perangkat OPENWRT Anda terhubung ke internet."
        exit 1
    fi
else
    echo "ADB sudah terinstall."
fi

exit 0