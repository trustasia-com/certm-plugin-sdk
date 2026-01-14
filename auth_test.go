package certm

import (
	"archive/zip"
	"crypto/ed25519"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestVerify(t *testing.T) {
	// 1. 生成密钥对
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	// 2. 准备测试文件
	tmpDir, err := os.MkdirTemp("", "plugin-test-*")
	if err != nil {
		t.Fatal(err)
	}
	// nolint:errcheck
	defer os.RemoveAll(tmpDir)

	pluginSo := filepath.Join(tmpDir, "plugin.so")
	_ = os.WriteFile(pluginSo, []byte("fake binary content"), 0644)

	pluginYaml := filepath.Join(tmpDir, "plugin.yaml")
	_ = os.WriteFile(pluginYaml, []byte("name: test-plugin"), 0644)

	// 3. 创建 Manifest
	manifest, err := CreateManifest([]string{"plugin.so", "plugin.yaml"}, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	manifestData, _ := json.Marshal(manifest)
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	_ = os.WriteFile(manifestPath, manifestData, 0644)

	// 4. 生成签名
	sig := Sign(manifestData, priv)
	sigPath := filepath.Join(tmpDir, "signature")
	_ = os.WriteFile(sigPath, sig, 0644)

	// 5. 打包 ZIP
	zipPath := filepath.Join(tmpDir, "plugin.zip")
	zipFile, _ := os.Create(zipPath)
	zw := zip.NewWriter(zipFile)

	files := []string{"plugin.so", "plugin.yaml", "manifest.json", "signature"}
	for _, name := range files {
		fw, _ := zw.Create(name)
		data, _ := os.ReadFile(filepath.Join(tmpDir, name))
		_, _ = fw.Write(data)
	}
	_ = zw.Close()
	_ = zipFile.Close()

	// 6. 执行验证
	t.Run("Valid Signature", func(t *testing.T) {
		err := Verify(zipPath, pub)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		_, otherPriv, _ := ed25519.GenerateKey(nil)
		badSig := Sign(manifestData, otherPriv)

		// 创建一个带错误签名的 zip
		badZipPath := filepath.Join(tmpDir, "bad_sig.zip")
		f, _ := os.Create(badZipPath)
		w := zip.NewWriter(f)
		for _, name := range []string{"plugin.so", "plugin.yaml", "manifest.json"} {
			fw, _ := w.Create(name)
			data, _ := os.ReadFile(filepath.Join(tmpDir, name))
			_, _ = fw.Write(data)
		}
		fw, _ := w.Create("signature")
		_, _ = fw.Write(badSig)
		_ = w.Close()
		_ = f.Close()

		err := Verify(badZipPath, pub)
		if err == nil || err.Error() != "invalid signature" {
			t.Errorf("expected invalid signature error, got %v", err)
		}
	})

	t.Run("Tampered File", func(t *testing.T) {
		// 创建一个被篡改文件的 zip
		tamperedZipPath := filepath.Join(tmpDir, "tampered.zip")
		f, _ := os.Create(tamperedZipPath)
		w := zip.NewWriter(f)
		fw, _ := w.Create("plugin.so")
		_, _ = fw.Write([]byte("tampered content")) // 修改内容

		for _, name := range []string{"plugin.yaml", "manifest.json", "signature"} {
			fw, _ := w.Create(name)
			data, _ := os.ReadFile(filepath.Join(tmpDir, name))
			_, _ = fw.Write(data)
		}
		_ = w.Close()
		_ = f.Close()

		err := Verify(tamperedZipPath, pub)
		if err == nil || !testing.Short() && err.Error() == "" {
			t.Errorf("expected integrity check failure, got %v", err)
		}
	})
}
