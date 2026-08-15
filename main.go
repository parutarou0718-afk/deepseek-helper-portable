// DeepSeek 助手 便携版启动器（分发给别人用）
// 完全自包含：使用同目录下的 node.exe + node_modules，不依赖对方安装 Node/npm。
// 双击 -> 检测 127.0.0.1:3080 -> 未运行则后台拉起内置 dsh web（DSH_HOME 指向内置 dsh-home）
//      -> 等端口就绪 -> 打开默认浏览器。
package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	host   = "127.0.0.1"
	port   = "3080"
	webURL = "http://127.0.0.1:3080"

	createNoWindow        = 0x08000000
	createNewProcessGroup = 0x00000200
)

func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(p)
}

func portOpen() bool {
	c, err := net.DialTimeout("tcp", host+":"+port, 300*time.Millisecond)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func openURL(u string) {
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
}

// loadingURL 返回启动加载页的 file:// 地址（正确处理中文/空格路径）。
func loadingURL() string {
	p := filepath.Join(exeDir(), "loading.html")
	u := &url.URL{Scheme: "file", Path: "/" + filepath.ToSlash(p)}
	return u.String()
}

func messageBox(title, text string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString(text)
	ti, _ := syscall.UTF16PtrFromString(title)
	proc.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(ti)), 0x10)
}

// envWith 返回一份把同名变量（忽略大小写）替换成 value 的环境变量列表。
func envWith(key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(os.Environ())+1)
	for _, e := range os.Environ() {
		eq := strings.IndexByte(e, '=')
		if eq > 0 && strings.EqualFold(e[:eq], key) {
			continue
		}
		out = append(out, e)
	}
	return append(out, prefix+value)
}

func startServer() error {
	dir := exeDir()
	node := filepath.Join(dir, "node.exe")
	dshBin := filepath.Join(dir, "node_modules", "@deepseek-ai", "dsh", "lib", "bin.js")
	home := filepath.Join(dir, "dsh-home")
	work := filepath.Join(dir, "workspace")

	if err := os.MkdirAll(home, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(work, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(dir, "dsh.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command(node, dshBin, "web")
	cmd.Dir = work
	cmd.Env = envWith("DSH_HOME", home)
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: createNoWindow | createNewProcessGroup,
		HideWindow:    true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	_ = os.WriteFile(filepath.Join(dir, "dsh.pid"), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	return nil
}

func main() {
	check := false
	for _, a := range os.Args[1:] {
		if a == "--check" {
			check = true
		}
	}

	dir := exeDir()

	if check {
		node := filepath.Join(dir, "node.exe")
		dshBin := filepath.Join(dir, "node_modules", "@deepseek-ai", "dsh", "lib", "bin.js")
		lines := []string{
			"exeDir=" + dir,
			fmt.Sprintf("nodeExists=%v", fileExists(node)),
			fmt.Sprintf("dshBinExists=%v", fileExists(dshBin)),
			fmt.Sprintf("loadingHtmlExists=%v", fileExists(filepath.Join(dir, "loading.html"))),
			fmt.Sprintf("portOpen=%v", portOpen()),
		}
		report := ""
		for _, l := range lines {
			report += l + "\n"
			fmt.Println(l)
		}
		_ = os.WriteFile(filepath.Join(dir, "check.txt"), []byte(report), 0644)
		return
	}

	if portOpen() {
		openURL(webURL)
		return
	}

	// 先立刻打开加载页给用户即时反馈，再后台拉起服务；
	// 服务就绪后由加载页自动跳转到真实页面。
	openURL(loadingURL())

	if err := startServer(); err != nil {
		messageBox("启动失败", "无法启动服务：\n"+err.Error()+"\n\n请查看本目录下的 dsh.log 日志。")
		return
	}
}

func fileExists(p string) bool {
	s, err := os.Stat(p)
	return err == nil && !s.IsDir()
}
