package main

import (
	"context"
	"os"
	"sync"
	"time"

	"screen-translator/internal/capture"
	"screen-translator/internal/config"
	"screen-translator/internal/translate"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App はWailsアプリケーションのメイン構造体。
type App struct {
	ctx        context.Context
	cfg        config.Config
	capturer   *capture.ScreenCapturer
	translator *translate.ClaudeTranslator

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
}

// NewApp は新しい App を返す。
func NewApp() *App {
	return &App{}
}

// startup はアプリ起動時にWailsから呼ばれる。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load()
	if err != nil {
		runtime.LogWarningf(ctx, "config load failed: %v", err)
		cfg = config.Default()
	}
	a.cfg = cfg
	a.capturer = capture.New()
	a.translator = translate.New(os.Getenv("ANTHROPIC_API_KEY"))
}

// --- フロントエンドから呼ばれるメソッド ---

// StartTranslation はキャプチャ・翻訳ループを開始する。
func (a *App) StartTranslation() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return
	}
	a.running = true
	a.stopCh = make(chan struct{})

	go a.loop()
}

// StopTranslation はキャプチャ・翻訳ループを停止する。
func (a *App) StopTranslation() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return
	}
	close(a.stopCh)
	a.running = false
}

// GetConfig は現在の設定を返す。
func (a *App) GetConfig() config.Config {
	return a.cfg
}

// SaveConfig は設定を保存し、ランタイムに反映する。
func (a *App) SaveConfig(cfg config.Config) error {
	if err := config.Save(cfg); err != nil {
		return err
	}
	a.mu.Lock()
	a.cfg = cfg
	a.mu.Unlock()
	return nil
}

// DisplayCount は接続モニター数を返す。
func (a *App) DisplayCount() int {
	return capture.DisplayCount()
}

// --- 内部ループ ---

func (a *App) loop() {
	a.runOnce()

	ticker := time.NewTicker(time.Duration(a.cfg.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ticker.C:
			a.runOnce()
		}
	}
}

func (a *App) runOnce() {
	a.mu.Lock()
	displayIndex := a.cfg.DisplayIndex
	a.mu.Unlock()

	png, err := a.capturer.Capture(displayIndex)
	if err != nil {
		runtime.LogErrorf(a.ctx, "capture failed: %v", err)
		runtime.EventsEmit(a.ctx, "translation:error", err.Error())
		return
	}
	if png == nil {
		// 差分なし
		return
	}

	result, err := a.translator.Translate(a.ctx, png)
	if err != nil {
		runtime.LogErrorf(a.ctx, "translate failed: %v", err)
		runtime.EventsEmit(a.ctx, "translation:error", err.Error())
		return
	}

	runtime.EventsEmit(a.ctx, "translation:result", result)
}
