package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/hoppscotch/proxyscotch/libfs"
	"github.com/pkg/browser"

	icon "github.com/hoppscotch/proxyscotch/icons"
	"github.com/hoppscotch/proxyscotch/inputbox"
	"github.com/hoppscotch/proxyscotch/libproxy"
	"github.com/hoppscotch/proxyscotch/notifier"
)

var (
	VersionName string
	VersionCode string
)

var (
	mStatus          *systray.MenuItem
	mCopyAccessToken *systray.MenuItem
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Proxyscotch v" + VersionName + " (" + VersionCode + ") - created by NBTX")

	/** Set up menu items. **/

	// Status
	mStatus = systray.AddMenuItem("启动中...", "")
	mStatus.Disable()
	mCopyAccessToken = systray.AddMenuItem("复制 Access Token...", "")
	mCopyAccessToken.Disable()

	systray.AddSeparator()

	// Open Hoppscotch Interface
	mOpenHoppscotch := systray.AddMenuItem("打开 Hoppscotch", "")

	systray.AddSeparator()

	// Set Proxy Authentication Token
	mSetAccessToken := systray.AddMenuItem("设置 Access Token...", "")

	systray.AddSeparator()

	// Quit Proxy
	mQuit := systray.AddMenuItem("退出 Proxyscotch", "")
	var fsPort = 33633

	/** Start proxy server. **/
	go runHoppscotchProxy(fsPort)
	// 33633 端口上监听文件服务器
	go runFileServer(fsPort)

	/** Wait for menu input. **/
	for {
		select {
		case <-mOpenHoppscotch.ClickedCh:
			_ = browser.OpenURL(fmt.Sprintf("http://127.0.0.1:%d/", fsPort))

		case <-mCopyAccessToken.ClickedCh:
			_ = clipboard.WriteAll(libproxy.GetAccessToken())
			_ = notifier.Notify("Proxyscotch", "代理 Access Token 已复制...", "代理 Access Token 已经复制到剪切板.", notifier.GetIcon())

		case <-mSetAccessToken.ClickedCh:
			newAccessToken, success := inputbox.InputBox("Proxyscotch", "请输入新的代理Access Token ...\n(Leave this blank to disable access checks.)", "")
			if success {
				libproxy.SetAccessToken(newAccessToken)

				if len(newAccessToken) == 0 {
					_ = notifier.Notify("Proxyscotch", "未开启代理访问控制...", "**任何人都可以访问你的代理服务!** 未开启代理访问控制.", notifier.GetIcon())
				} else {
					_ = notifier.Notify("Proxyscotch", "更新代理访问控制成功...", "代理 Access Token 已经成功更新.", notifier.GetIcon())
				}
			}

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onExit() {
}

func runHoppscotchProxy(fsPort int) {
	libproxy.Initialize("", "127.0.0.1:9159", fmt.Sprintf("http://127.0.0.1:%d", fsPort), "", "", onProxyStateChange, false, nil)
}

func runFileServer(fsPort int) {
	libfs.InitializeFs(fsPort)
}

func onProxyStateChange(status string, isListening bool) {
	mStatus.SetTitle(status)

	if isListening {
		mCopyAccessToken.Enable()
	}
}
