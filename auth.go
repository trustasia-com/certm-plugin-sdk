package certm

import (
	"archive/zip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Manifest 插件清单文件
type Manifest struct {
	Files     map[string]string `json:"files"`     // 文件名 -> sha256
	Timestamp int64             `json:"timestamp"` // 签名时间戳
}

// Verify 验证插件包签名
// zipPath: ZIP压缩包路径
// publicKey: 开发者公钥
func Verify(zipPath string, publicKey ed25519.PublicKey) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	// nolint:errcheck
	defer r.Close()

	var (
		manifestFile *zip.File
		signature    []byte
	)

	// 1. 查找关键文件
	for _, f := range r.File {
		switch f.Name {
		case "manifest.json":
			manifestFile = f
		case "signature":
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("open signature file: %w", err)
			}
			signature, err = io.ReadAll(rc)
			_ = rc.Close()
			if err != nil {
				return fmt.Errorf("read signature file: %w", err)
			}
		}
	}

	if manifestFile == nil {
		return fmt.Errorf("missing manifest.json")
	}
	if len(signature) == 0 {
		return fmt.Errorf("missing signature")
	}

	// 2. 验证签名
	rc, err := manifestFile.Open()
	if err != nil {
		return fmt.Errorf("open manifest.json: %w", err)
	}
	manifestData, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil {
		return fmt.Errorf("read manifest.json: %w", err)
	}

	if !ed25519.Verify(publicKey, manifestData, signature) {
		return fmt.Errorf("invalid signature")
	}

	// 3. 验证完整性
	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("unmarshal manifest.json: %w", err)
	}

	for _, f := range r.File {
		if f.Name == "manifest.json" || f.Name == "signature" {
			continue
		}

		expectedHash, ok := manifest.Files[f.Name]
		if !ok {
			return fmt.Errorf("file %s not in manifest", f.Name)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open file %s: %w", f.Name, err)
		}

		h := sha256.New()
		if _, err := io.Copy(h, rc); err != nil {
			_ = rc.Close()
			return fmt.Errorf("hash file %s: %w", f.Name, err)
		}
		_ = rc.Close()

		actualHash := hex.EncodeToString(h.Sum(nil))
		if actualHash != expectedHash {
			return fmt.Errorf("file %s integrity check failed: expected %s, got %s", f.Name, expectedHash, actualHash)
		}
	}

	return nil
}

// Sign 对清单进行签名（用于开发者工具）
func Sign(manifestData []byte, privateKey ed25519.PrivateKey) []byte {
	return ed25519.Sign(privateKey, manifestData)
}

// CreateManifest 创建清单（用于开发者工具）
func CreateManifest(files []string, baseDir string) (*Manifest, error) {
	manifest := &Manifest{
		Files: make(map[string]string),
	}
	for _, f := range files {
		path := f
		if baseDir != "" {
			path = baseDir + "/" + f
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		h := sha256.Sum256(data)
		manifest.Files[f] = hex.EncodeToString(h[:])
	}
	return manifest, nil
}
