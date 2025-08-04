package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {

	// 取得當前程式所在目錄
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("❌ 無法取得當前路徑：", err)
		os.Exit(1)
	}
	os.Chdir(execDir)

	// 檢查 bedrock_server_mod.exe 是否存在
	if _, err := os.Stat("bedrock_server_mod.exe"); os.IsNotExist(err) {
		fmt.Println("⚙️ 偵測到尚未初始化 LiteLoader BDS，正在進行安裝，這可能需要幾分鐘...")

		// 執行初始化程序
		cmd := exec.Command("/usr/bin/wine", "PeEditor.exe")
		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf

		if err := cmd.Run(); err != nil {
			fmt.Println("❌ 初始化失敗：", err)
			fmt.Println(outBuf.String())
			os.Exit(1)
		}

		fmt.Println("✅ 初始化完成。")
	}

	// 建立 context 控制程序結束
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動伺服器
	cmd := exec.CommandContext(ctx, "/usr/bin/wine", "bedrock_server_mod.exe")

	// 連接標準輸入輸出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("🚀 正在啟動伺服器...")

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ 無法啟動伺服器：", err)
		os.Exit(1)
	}

	// 等待程序結束（阻塞）
	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("⚠️ 伺服器異常結束：", err)
		} else {
			fmt.Println("✅ 伺服器已正常結束。")
		}
		cancel()
	}()

	// 開始讀取使用者輸入並傳送至伺服器（保持互動）
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 伺服器已關閉。")
			return
		default:
			// 讀取輸入並寫入 stdin（目前直接綁定，這裡可以省略）
			time.Sleep(100 * time.Millisecond)
		}
	}
}