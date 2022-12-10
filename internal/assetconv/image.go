package assetconv

import (
	"io"
	"os/exec"
	"strings"
)

// CanToWebP checks if a specified format can be converted to WebP.
func CanToWebP(format string) bool {
	// PNG, JPEG, and TIFF common exts
	supported := []string{"png", "jpg", "jpeg", "tiff", "tif"}
	for _, v := range supported {
		if strings.ToLower(v) == format {
			return true
		}
	}
	return false
}

// ToWebP converts supported images to WebP format by calling `cwebp`.
// Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
func ToWebP(image []byte) ([]byte, error) {
	cmd := exec.Command("cwebp", "-quiet", "-o", "-", "--", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(image))
	}()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	ret, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
