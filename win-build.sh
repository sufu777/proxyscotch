# 导入版本信息
. ./version.properties

[ -d "out/" ] || mkdir "out/"

[ -f "rsrc.syso" ] && rm rsrc.syso
go get github.com/akavel/rsrc
go install github.com/akavel/rsrc

rsrc -arch="amd64" -manifest="$OUTPUT_DIR/proxyscotch.manifest" -ico="icons/icon.ico" -o rsrc.syso
CGO_ENABLED=1 GOOS="windows" GOARCH="amd64" go build -ldflags "-X main.VersionName=$VERSION_NAME -X main.VersionCode=$VERSION_CODE -H=windowsgui" -o "$OUTPUT_DIR/proxyscotch-amd64.exe"
rm rsrc.syso

mkdir "$OUTPUT_DIR/icons"
cp icons/icon.png "$OUTPUT_DIR/icons/icon.png"

mkdir "$OUTPUT_DIR/data"

rm "$OUTPUT_DIR/proxyscotch.manifest"

mv "$OUTPUT_DIR/proxyscotch-amd64.exe" "$OUTPUT_DIR/Proxyscotch-Desktop-Windows-amd64-v${VERSION_NAME}.exe"