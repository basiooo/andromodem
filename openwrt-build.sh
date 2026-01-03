#!/bin/bash
set -e

# Check if container is running, start if not
if ! docker-compose ps | grep -q "Up"; then
    echo "Starting container..."
    docker-compose up -d
    sleep 3
fi

echo "Building AndroModem packages..."

docker-compose exec openwrt-builder bash -c '
    # Download SDK if not done
    [ ! -d ./scripts ] && ./setup.sh
    
    # Create default config
    [ ! -f .config ] && make defconfig
    
    grep -q "github.com/openwrt" feeds.conf.default || \
        sed -i -E "s;git.openwrt.org/(feed|project|openwrt);github.com/openwrt;" feeds.conf.default
    
    [ ! -d ./feeds/base ] && ./scripts/feeds update base
    [ ! -d ./feeds/packages ] && ./scripts/feeds update packages
    
    # Install feeds
    ./scripts/feeds install -a
    
    # Build package
    make package/luci-app-andromodem/{clean,compile} -j 12 V=s
'

echo ""
echo "Copying packages to host..."
docker-compose exec openwrt-builder bash -c '
    mkdir -p /output
    find bin/ -name "*.ipk" -exec cp {} /output/ \;
'

echo ""
echo "Done! Packages in ./openwrt-build/"
ls -lh openwrt-build/*.ipk 2>/dev/null || echo "No packages found"
