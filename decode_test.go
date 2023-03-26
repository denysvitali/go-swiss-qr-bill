package swiss_qr_code_test

import (
	"encoding/json"
	"fmt"
	sqc "github.com/denysvitali/go-swiss-qr-bill"
	. "github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"testing"
)

func getTestFile(t *testing.T, number int) string {
	return getFileByPath(t, fmt.Sprintf("./resources/samples/%d.txt", number))
}

func getInvalidTestFile(t *testing.T, number int) string {
	return getFileByPath(t, fmt.Sprintf("./resources/invalid/%d.txt", number))
}

func getFileByPath(t *testing.T, path string) string {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("unable to get test file: %v", err)
		return ""
	}
	text, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("unable to read test file: %v", err)
		return ""
	}

	return string(text)
}

func TestAll(t *testing.T) {
	mainPath := "./resources/samples"
	samplesDir, err := os.ReadDir(mainPath)
	if err != nil {
		t.Fatalf("unable to open samples dir: %v", err)
	}

	for _, v := range samplesDir {
		if !v.IsDir() {
			f, err := os.Open(path.Join(mainPath, v.Name()))
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			fileContent, err := io.ReadAll(f)
			if err != nil {
				t.Fatalf("unable to read file contents: %v", err)
			}

			qrCode, err := sqc.Decode(string(fileContent))
			if err != nil {
				t.Fatalf("unable to decode QR bill: %v", err)
			}

			Equal(t, sqc.SwissPaymentsCodeQrType, qrCode.Header.QRType)
			Equal(t, 2, qrCode.Header.Version.Major)
			Equal(t, 0, qrCode.Header.Version.Minor)
			Equal(t, sqc.Utf8, qrCode.Header.CodingType)

			s, _ := json.MarshalIndent(qrCode, "", "\t")
			fmt.Printf("%s\n", s)
		}
	}
}

func TestInvalid(t *testing.T) {
	testFile := getInvalidTestFile(t, 1)
	_, err := sqc.Decode(testFile)
	if err == nil {
		t.Fatalf("this should fail")
	}
}
