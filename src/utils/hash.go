package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createFilesAndHash() error {
	// Criar o diretório tmp/dataset/ se não existir
	dir := "tmp/dataset/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	// Nomes dos arquivos
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"}

	// Map para armazenar os hashes
	fileHashes := make(map[string]string)

	// Criar arquivos e calcular hashes
	for _, file := range files {
		path := filepath.Join(dir, file)
		content := fmt.Sprintf("Conteúdo de %s", file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %v", file, err)
		}

		// Calcular o hash SHA-256 do arquivo
		hash, err := calculateHash(path)
		if err != nil {
			return fmt.Errorf("erro ao calcular hash de %s: %v", file, err)
		}

		// Armazenar o hash no map
		fileHashes[file] = hash
	}

	// Imprimir os hashes no final
	fmt.Println("Hashes dos arquivos:")
	for file, hash := range fileHashes {
		fmt.Printf("%s: %s\n", file, hash)
	}

	return nil
}

func calculateHash(filePath string) (string, error) {
	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Inicializar o hash SHA-256
	hasher := sha256.New()

	// Copiar o conteúdo do arquivo para o hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Retornar o hash em formato hexadecimal
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func main() {
	// Executar a função para criar arquivos e calcular hashes
	if err := createFilesAndHash(); err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Println("Arquivos criados e hashes calculados com sucesso!")
	}
}
