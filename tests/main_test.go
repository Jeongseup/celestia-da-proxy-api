package tests

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"testing"

	"github.com/chai2010/webp"
)

func TestConvert(t *testing.T) {
	// 현재 작업 디렉토리 확인
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}
	t.Logf("Current working directory: %s", workingDir)

	testCases := []struct {
		fileName string
		// inputPath  string
		// outputPath string
	}{
		{
			fileName: "first",
			// inputPath:  "../assets/first.webp",
			// outputPath: "../assets/tmp/first.png",
		},
	}

	for _, item := range testCases {

		webpFile := fmt.Sprintf("../assets/webp/%s.webp", item.fileName)
		// pngFile := fmt.Sprintf("../assets/png/%s.png", item.fileName)
		ppmFile := fmt.Sprintf("../assets/ppm/%s.png", item.fileName)
		svgFile := fmt.Sprintf("../assets/svg/%s.svg", item.fileName)

		if err := ConvertWebPToPPM(webpFile, ppmFile); err != nil {
			t.Errorf("Error converting WebP to PNG: %s\n", err)
		}

		// PNG를 SVG로 변환
		if err := ConvertPPMToSVG(ppmFile, svgFile); err != nil {
			t.Errorf("Error converting PNG to SVG: %s\n", err)
		}

		t.Log("Conversion successful!")
	}
}

// ConvertWebPToPNG 함수는 WebP 파일을 PNG 파일로 변환합니다.
func ConvertWebPToPNG(inputFile, outputFile string) error {
	// WebP 파일을 읽기
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// WebP 파일 디코딩
	img, err := webp.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode webp file: %w", err)
	}

	// PNG 파일로 저장
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("failed to encode png file: %w", err)
	}

	return nil
}

// ConvertPNGToSVG 함수는 PNG 파일을 SVG 파일로 변환합니다.
func ConvertPNGToSVG(inputFile, outputFile string) error {
	// potrace 명령 실행
	cmd := exec.Command("potrace", inputFile, "-s", "-o", outputFile)

	// 표준 출력 및 표준 오류를 캡처합니다.
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 명령 실행
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert png to svg: %s: %w", stderr.String(), err)
	}

	fmt.Printf("Output: %s\n", out.String())
	return nil
}

// ConvertWebPToPPM 함수는 WebP 파일을 PPM 파일로 변환합니다.
func ConvertWebPToPPM(inputFile, outputFile string) error {
	// WebP 파일 열기
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// WebP 파일 디코딩
	img, err := webp.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode webp file: %w", err)
	}

	// PPM 파일로 저장
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// PPM 포맷으로 인코딩
	if err := EncodePPM(outFile, img); err != nil {
		return fmt.Errorf("failed to encode ppm file: %w", err)
	}

	return nil
}

// EncodePPM 함수는 이미지를 PPM 포맷으로 인코딩합니다.
func EncodePPM(outFile *os.File, img image.Image) error {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// PPM 헤더 작성
	_, err := fmt.Fprintf(outFile, "P3\n%d %d\n255\n", width, height)
	if err != nil {
		return fmt.Errorf("failed to write ppm header: %w", err)
	}

	// 이미지 데이터 작성
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			_, err := fmt.Fprintf(outFile, "%d %d %d ", r>>8, g>>8, b>>8)
			if err != nil {
				return fmt.Errorf("failed to write ppm data: %w", err)
			}
		}
		_, err := fmt.Fprint(outFile, "\n")
		if err != nil {
			return fmt.Errorf("failed to write ppm newline: %w", err)
		}
	}

	return nil
}

// ConvertPPMToSVG 함수는 PPM 파일을 SVG 파일로 변환합니다.
func ConvertPPMToSVG(inputFile, outputFile string) error {
	// potrace 명령 실행
	cmd := exec.Command("potrace", inputFile, "-s", "-o", outputFile)

	// 표준 출력 및 표준 오류를 캡처합니다.
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 명령 실행
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert ppm to svg: %s: %w", stderr.String(), err)
	}

	fmt.Printf("Output: %s\n", out.String())
	return nil
}
