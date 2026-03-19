import { useState, useEffect, useCallback } from "react";
import { EventsOn } from "../wailsjs/runtime/runtime";
import {
  StartTranslation,
  StopTranslation,
  GetConfig,
  SaveConfig,
  DisplayCount,
} from "../wailsjs/go/main/App";
import { config as configModels } from "../wailsjs/go/models";

type Status = "idle" | "running" | "error";

function App() {
  const [status, setStatus]           = useState<Status>("idle");
  const [translation, setTranslation] = useState("");
  const [updatedAt, setUpdatedAt]     = useState<Date | null>(null);
  const [errorMsg, setErrorMsg]       = useState("");
  const [showSettings, setShowSettings] = useState(false);
  const [config, setConfig]           = useState<configModels.Config | null>(null);
  const [displayCount, setDisplayCount] = useState(1);
  const [toast, setToast]             = useState("");

  // 初期化
  useEffect(() => {
    GetConfig().then(setConfig);
    DisplayCount().then(setDisplayCount);

    EventsOn("translation:result", (result: string) => {
      setTranslation(result);
      setUpdatedAt(new Date());
      setStatus("running");
    });

    EventsOn("translation:error", (msg: string) => {
      setErrorMsg(msg);
      setStatus("error");
      showToast("エラーが発生しました");
    });
  }, []);

  const showToast = useCallback((msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(""), 5000);
  }, []);

  const handleStart = () => {
    setStatus("running");
    setErrorMsg("");
    StartTranslation();
  };

  const handleStop = () => {
    setStatus("idle");
    StopTranslation();
  };

  const handleSaveConfig = async () => {
    if (!config) return;
    await SaveConfig(config);
    setShowSettings(false);
    showToast("設定を保存しました");
  };

  const statusColor = {
    idle:    "text-gray-400",
    running: "text-success",
    error:   "text-danger",
  }[status];

  const statusLabel = {
    idle:    "停止中",
    running: "翻訳中",
    error:   "エラー",
  }[status];

  return (
    <div className="flex flex-col h-screen bg-bg text-text">
      {/* ステータスバー */}
      <header className="flex items-center justify-between px-4 py-2 bg-card border-b border-white/10">
        <div className="flex items-center gap-2">
          <span className={`text-sm font-medium ${statusColor}`}>● {statusLabel}</span>
        </div>
        <div className="flex items-center gap-2">
          {status === "running" ? (
            <button
              onClick={handleStop}
              className="px-3 py-1 text-sm rounded bg-danger/20 text-danger hover:bg-danger/30 transition-colors"
            >
              停止
            </button>
          ) : (
            <button
              onClick={handleStart}
              className="px-3 py-1 text-sm rounded bg-accent/20 text-accent hover:bg-accent/30 transition-colors"
            >
              開始
            </button>
          )}
          <button
            onClick={() => setShowSettings(true)}
            className="px-3 py-1 text-sm rounded bg-white/10 hover:bg-white/20 transition-colors"
          >
            設定
          </button>
        </div>
      </header>

      {/* 翻訳テキストエリア */}
      <main className="flex-1 overflow-y-auto p-6">
        {translation ? (
          <p
            className="whitespace-pre-wrap leading-relaxed"
            style={{ fontSize: config?.font_size ?? 18 }}
          >
            {translation}
          </p>
        ) : (
          <p className="text-white/30 text-center mt-20">
            「開始」を押すと翻訳が始まります
          </p>
        )}
      </main>

      {/* 最終更新時刻 */}
      <footer className="px-4 py-2 text-xs text-white/30 border-t border-white/10">
        {updatedAt
          ? `最終更新: ${updatedAt.toLocaleTimeString("ja-JP")}`
          : "未更新"}
      </footer>

      {/* 設定モーダル */}
      {showSettings && config && (
        <div className="absolute inset-0 bg-black/60 flex items-center justify-center z-10">
          <div className="bg-card rounded-xl p-6 w-96 space-y-5 shadow-xl">
            <h2 className="text-lg font-bold">設定</h2>

            {/* モニター選択 */}
            <label className="block space-y-1">
              <span className="text-sm text-white/60">キャプチャするモニター</span>
              <select
                className="w-full bg-bg border border-white/20 rounded px-3 py-2 text-sm"
                value={config.display_index}
                onChange={(e) =>
                  setConfig(configModels.Config.createFrom({ ...config, display_index: Number(e.target.value) }))
                }
              >
                {Array.from({ length: displayCount }, (_, i) => (
                  <option key={i} value={i}>
                    モニター {i + 1}
                  </option>
                ))}
              </select>
            </label>

            {/* インターバル */}
            <label className="block space-y-1">
              <span className="text-sm text-white/60">
                キャプチャ間隔: {config.interval_seconds}秒
              </span>
              <input
                type="range"
                min={1}
                max={30}
                value={config.interval_seconds}
                onChange={(e) =>
                  setConfig(configModels.Config.createFrom({ ...config, interval_seconds: Number(e.target.value) }))
                }
                className="w-full accent-accent"
              />
            </label>

            {/* フォントサイズ */}
            <label className="block space-y-1">
              <span className="text-sm text-white/60">
                フォントサイズ: {config.font_size}px
              </span>
              <input
                type="range"
                min={12}
                max={32}
                value={config.font_size}
                onChange={(e) =>
                  setConfig(configModels.Config.createFrom({ ...config, font_size: Number(e.target.value) }))
                }
                className="w-full accent-accent"
              />
            </label>

            {/* APIキー */}
            <label className="block space-y-1">
              <span className="text-sm text-white/60">Anthropic APIキー</span>
              <input
                type="password"
                placeholder="sk-ant-..."
                value={config.api_key ?? ""}
                onChange={(e) =>
                  setConfig(configModels.Config.createFrom({ ...config, api_key: e.target.value }))
                }
                className="w-full bg-bg border border-white/20 rounded px-3 py-2 text-sm font-mono"
              />
              <p className="text-xs text-white/30">
                console.anthropic.com で取得できます
              </p>
            </label>

            <div className="flex justify-end gap-2 pt-2">
              <button
                onClick={() => setShowSettings(false)}
                className="px-4 py-2 text-sm rounded bg-white/10 hover:bg-white/20 transition-colors"
              >
                キャンセル
              </button>
              <button
                onClick={handleSaveConfig}
                className="px-4 py-2 text-sm rounded bg-accent hover:bg-accent/80 transition-colors"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}

      {/* トースト通知 */}
      {toast && (
        <div className="absolute bottom-6 left-1/2 -translate-x-1/2 bg-card border border-white/20 px-4 py-2 rounded-lg text-sm shadow-lg">
          {toast}
        </div>
      )}
    </div>
  );
}

export default App;
