package engine

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"strings"
)

// ChecksumResult holds computed cryptographic hashes of a downloaded file.
type ChecksumResult struct {
	SHA256   string `json:"sha256"`
	MD5      string `json:"md5"`
	SHA1     string `json:"sha1,omitempty"`
	Verified *bool  `json:"verified,omitempty"`
}

// ComputeFileChecksums calculates SHA-256 and MD5 in a single streaming pass over the file
// using a 128KB pooled buffer.
func ComputeFileChecksums(filePath string, expectedHash string) (*ChecksumResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h256 := sha256.New()
	hMD5 := md5.New()

	var writers []io.Writer
	writers = append(writers, h256, hMD5)

	var hSHA1 hash.Hash
	cleanExpected := strings.ToLower(strings.TrimSpace(expectedHash))
	if len(cleanExpected) == 40 {
		hSHA1 = sha1.New()
		writers = append(writers, hSHA1)
	}

	multi := io.MultiWriter(writers...)
	bufPtr := BufferPool.Get().(*[]byte)
	defer BufferPool.Put(bufPtr)

	if _, err := io.CopyBuffer(multi, f, *bufPtr); err != nil {
		return nil, err
	}

	res := &ChecksumResult{
		SHA256: hex.EncodeToString(h256.Sum(nil)),
		MD5:    hex.EncodeToString(hMD5.Sum(nil)),
	}
	if hSHA1 != nil {
		res.SHA1 = hex.EncodeToString(hSHA1.Sum(nil))
	}

	if cleanExpected != "" {
		matched := strings.EqualFold(cleanExpected, res.SHA256) ||
			strings.EqualFold(cleanExpected, res.MD5) ||
			(res.SHA1 != "" && strings.EqualFold(cleanExpected, res.SHA1))
		res.Verified = &matched
	}

	return res, nil
}
