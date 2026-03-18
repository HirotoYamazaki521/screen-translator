package capture

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/kbinani/screenshot"
)

// Capturer はスクリーンキャプチャのインターフェース。
type Capturer interface {
	// Capture は指定モニターをキャプチャしてPNGバイト列を返す。
	// 前回と差分がない場合は (nil, nil) を返す。
	Capture(displayIndex int) ([]byte, error)
}

// ScreenCapturer は実際のスクリーンキャプチャを行う。
type ScreenCapturer struct {
	lastHash string
}

// New は新しい ScreenCapturer を返す。
func New() *ScreenCapturer {
	return &ScreenCapturer{}
}

// DisplayCount は接続されているモニターの数を返す。
func DisplayCount() int {
	return screenshot.NumActiveDisplays()
}

// Capture は指定モニターをキャプチャする。
// 前回と差分がない場合は (nil, nil) を返す。
func (c *ScreenCapturer) Capture(displayIndex int) ([]byte, error) {
	bounds := screenshot.GetDisplayBounds(displayIndex)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}

	hash := hashImage(img)
	if hash == c.lastHash {
		return nil, nil
	}
	c.lastHash = hash

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("png encode failed: %w", err)
	}
	return buf.Bytes(), nil
}

// hashImage は画像を 16x16 にダウンスケールしてグレースケール変換した MD5 ハッシュを返す。
// タイムスタンプ等の微細変化は意図的に無視する（APIコスト削減）。
func hashImage(img image.Image) string {
	const size = 16
	small := image.NewGray(image.Rect(0, 0, size, size))

	src := img.Bounds()
	srcW := src.Dx()
	srcH := src.Dy()

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			srcX := src.Min.X + x*srcW/size
			srcY := src.Min.Y + y*srcH/size
			c := img.At(srcX, srcY)
			gray := color.GrayModel.Convert(c).(color.Gray)
			small.SetGray(x, y, gray)
		}
	}

	return fmt.Sprintf("%x", md5.Sum(small.Pix))
}
