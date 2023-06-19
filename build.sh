#!/bin/sh
echo "gofumpting"
gofumpt -w -l .

if [ -t 0 ]; then
	echo "building for your platform"
	wails build --debug --v 1 --webview2 embed -o "mf"
    echo "upxing"
	upx ./build/bin/mf -5
	echo "Executable can be found as /build/bin/mf"
	exit 0
fi

while IFS= read -r PLATFORM; do
	$EXTRAFLAGS=""
	if echo "$PLATFORM" | grep -iq "darwin" && [ "$GOOS" != "darwin" ]; then
		echo "Cross-compilation to darwin NOT supported"
	else
	echo "wails building for platform: $PLATFORM"
	PLATFORM_NAME=$(echo "$PLATFORM" | tr '/' '_')
	
	if echo "$PLATFORM" | grep -iq "windows"; then
		PLATFORM_NAME=$(echo "$PLATFORM_NAME.exe")
		# $EXTRAFLAGS="-nsis"
	fi

    wails build --debug --v 2 --webview2 embed --upx --platform "$PLATFORM" -o "mf_$PLATFORM_NAME" $EXTRAFLAGS

done