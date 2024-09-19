package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const fileSize = 256 

func createFilesAndHash() error {
	dir := "tmp/dataset/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	// Nomes dos arquivos
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"}

	fileHashes := make([]string, 0, len(files))

	for _, file := range files {
		path := filepath.Join(dir, file)

		randomContent, err := generateRandomContent(fileSize)
		if err != nil {
			return fmt.Errorf("erro ao gerar conteúdo aleatório para %s: %v", file, err)
		}

		if err := os.WriteFile(path, randomContent, 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %v", file, err)
		}

		hash, err := calculateHash(path)
		if err != nil {
			return fmt.Errorf("erro ao calcular hash de %s: %v", file, err)
		}

		fileHashes = append(fileHashes, hash)
	}

	fmt.Println("Hashes dos arquivos:")
	for i, hash := range fileHashes {
		fmt.Printf("file%d.txt: %s\n", i+1, hash)
	}

	if err := sendHashes(fileHashes); err != nil {
		return fmt.Errorf("erro ao enviar hashes: %v", err)
	}

	return nil
}

func generateRandomContent(size int) ([]byte, error) {
	content := make([]byte, size)
	_, err := rand.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func sendHashes(hashes []string) error {
	clientPath, err := filepath.Abs("../src/client/client.go")
	if err != nil {
		return fmt.Errorf("erro ao resolver o caminho absoluto: %v", err)
	}

	args := append([]string{"run", clientPath, "send"}, hashes...)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar comando: %v", err)
	}

	return nil
}

func main() {
	if err := createFilesAndHash(); err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Println("Arquivos criados, hashes calculados e enviados com sucesso!")
	}
}
