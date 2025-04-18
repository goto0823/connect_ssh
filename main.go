package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"encoding/csv"
)


type ServerInfo struct {
	HostName string
	User      string
	IP        string
	Path	 string
}


func main() {
	excelPath := filepath.Join(os.Getenv("USERPROFILE"), "Desktop", "key.csv")
	// excelPath := filepath.Join(os.Getenv("USERPROFILE"), "key.csv")

	servers, err := readServerData(excelPath)
	if err != nil {
		fmt.Println("Error reading server data:", err)
		return
	}

	if len(servers) == 0 {
		fmt.Println("サーバー情報が見つかりません")
		return
	}

	fmt.Println("=== サーバーリスト ===")
	for i, server := range servers {
		fmt.Printf("[%d] %s (%s@%s)\n", i+1, server.HostName, server.User, server.IP, server.Path)
	}

	selectedIndex := getSelectedServerIndex(len(servers))
	if selectedIndex < 0 {
		return
	}
	
	selectedServer := servers[selectedIndex]
	
	// SSH接続を実行
	connectSSH(selectedServer)
}


func readServerData(filePath string) ([]ServerInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("CSVファイルを開けません: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVファイルの読み込みに失敗しました: %w", err)
	}

	var servers []ServerInfo
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 4 {
			continue 
		}
		
		
		server := ServerInfo{
			HostName: strings.TrimSpace(row[0]),
			User:     strings.TrimSpace(row[1]),
			IP:       strings.TrimSpace(row[2]),
			Path:     strings.TrimSpace(row[3]),
		}
		
		if server.HostName != "" && server.User != "" && server.IP != "" && server.Path != "" {
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func getSelectedServerIndex(maxIndex int) int {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("接続するサーバーの番号を入力してください (qで終了): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "q" || input == "quit" || input == "exit" {
			return -1
		}
		
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > maxIndex {
			fmt.Printf("1から%dまでの番号を入力してください\n", maxIndex)
			continue
		}
		
		return index - 1
	}
}

// SSHで接続する
func connectSSH(server ServerInfo) {
	// 鍵のパスを確認
	keyPath := server.Path
	if !filepath.IsAbs(keyPath) {
		// 相対パスの場合、ユーザーのホームディレクトリからの相対パスとして扱う
		keyPath = filepath.Join(os.Getenv("USERPROFILE"), keyPath)
	}
	
	// コマンドを構築
	cmdArgs := []string{
		"-i", keyPath,
		fmt.Sprintf("%s@%s", server.User, server.IP),
	}
	
	fmt.Println(cmdArgs)
	fmt.Printf("接続中: %s (%s@%s)...\n\n", server.HostName, server.User, server.IP)
	
	// SSHコマンドを実行
	cmd := exec.Command("C:\\Windows\\System32\\OpenSSH\\ssh.exe", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("SSH接続エラー: %v\n", err)
		fmt.Println("何かキーを押すと終了します...")
		fmt.Scanln() 
	}
}