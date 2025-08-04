package main

import (
"bufio"
"bytes"
"fmt"
"os"
"os/exec"
"path/filepath"
)

func main() {
inR, inW, _ := os.Pipe()
outR, outW, _ := os.Pipe()
dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

done := make(chan struct{})  

process, _ := os.StartProcess("/bin/sh", nil, &os.ProcAttr{  
	Files: []*os.File{inR, outW, outW},  
	Dir:   dir,  
})  

// check bedrock_server_mod.exe if not exists  
if _, err := os.Stat("bedrock_server_mod.exe"); os.IsNotExist(err) {  
	fmt.Println("正在初始化 LeviLamina，這可能需要幾分鐘...")  

	cmd := exec.Command("/usr/bin/wine", "LLPeEditor.exe")  
	var out bytes.Buffer  
	cmd.Stdout = &out  
	err := cmd.Run()  
	fmt.Println(out.String(), err)  
}  

go func() {  
	// read console  
	reader := bufio.NewReader(os.Stdin)  
	writer := bufio.NewWriter(inW)  
	
	writer.WriteString("export WINEPREFIX=/home/container/.wine\n")
	writer.WriteString("wineboot --init\n")
	writer.Flush()
	writer.WriteString("winetricks -q vcrun2019\n")
	writer.Flush()
	writer.WriteString("wine bedrock_server_mod.exe\n")  
	writer.Flush()  

	for {  
		text, _ := reader.ReadString('\n')  
		inW.Write([]byte(text))  
		// writer.WriteString(text)  
		writer.Flush()  

	}  
}()  

go func() {  
	scanner := bufio.NewScanner(outR)  
	text := ""  
	for scanner.Scan() {  
		text = scanner.Text()  

		fmt.Println(text)  

		// if find := strings.Contains(text, "type \"help\" or \"?\""); find {  
		// 	fmt.Println("服务器启动成功，现在可以进入服务器了！")  
		// 	break  
		// }  

		// if find := strings.Contains(text, "Quit correctly"); find {  
		// 	// fmt.Println("服务器已正常关闭。")  
		// 	break  
		// }  

	}  
	process.Signal(os.Kill)  
	done <- struct{}{}  
	fmt.Println("伺服器已關閉。")  
}()  

process.Wait()  

// buffer := new(bytes.Buffer)  
// buffer.ReadFrom(outR)  
// fmt.Println(buffer.String())  

os.Exit(0)

}