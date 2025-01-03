package editor

import (
	"bytes"
	"os"
	"testing"
)

func TestDecompressXML(t *testing.T) {
	data, err := os.ReadFile("/Users/parapeng/Downloads/app-release/AndroidManifest.xml")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	m1 := ModifyInfo[string]{DefaultManifest.VersionName, "6.6.6"}
	m2 := ModifyInfo[uint32]{DefaultManifest.VersionCode, 666}
	m3 := ModifyInfo[string]{DefaultManifest.Label, "TestDemo"}
	m4 := ModifyInfo[string]{DefaultManifest.Package, "com.test.webvievwtr"}
	result, err := ModifyAll(data, m1, m2, m3, m4)
	if err != nil {
		t.Fatalf("Failed to modify manifest: %v", err)
	}
	os.WriteFile("/Users/parapeng/Downloads/app-release/AndroidManifest.new.xml", result, 0644)
}

func TestAdjustStringLength(t *testing.T) {
	testCases := []struct {
		name     string
		old      string
		new      string
		expected string
	}{
		{
			name:     "new longer than old",
			old:      "abc",
			new:      "abcdef",
			expected: "abc",
		},
		{
			name:     "new shorter than old",
			old:      "abcdef",
			new:      "abc",
			expected: "abc   ",
		},
		{
			name:     "same length",
			old:      "abc",
			new:      "def",
			expected: "def",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldBytes := changeString([]byte(tc.old))
			newBytes := changeString([]byte(tc.new))
			result := adjustStringLength(oldBytes, newBytes)
			expected := changeString([]byte(tc.expected))
			if !bytes.Equal(result, expected) {
				t.Errorf("Expected %x, got %x", expected, result)
			}
		})
	}
}
