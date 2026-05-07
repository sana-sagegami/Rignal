# CLAUDE.md — frontend (Rignal)

このファイルは Claude Code がプロジェクトを理解するためのコンテキストです。
作業前に必ず参照してください。

---

## プロジェクト概要

**Rignal** は Oura Ring 4 の生体データを活用したコンディション予測 × スケジュール最適化 iOS アプリです。

- 起動時に Google カレンダーから翌日の最初の予定時刻を取得し、バックエンドへ渡す
- バックエンドから「今日のコンディションスコア・集中ピーク時間帯・推奨就寝時刻」を取得してダッシュボードに表示する
- 推奨就寝時刻をローカル通知でリマインドする
- ホーム画面ウィジェットでコンディションスコアと集中ピーク時間帯を常時表示する
- App Store 非公開・個人利用のため APNs は使用しない

---

## 技術スタック

| 役割 | 技術 |
|---|---|
| UI フレームワーク | SwiftUI |
| プロジェクト管理 | Tuist |
| Google 認証 | Google Sign-In SDK (`GoogleSignIn`) |
| Google Calendar 取得 | Google API Client Library (`GoogleAPIClientForREST/Calendar`) |
| 通知 | UserNotifications framework（ローカル通知） |
| ウィジェット | WidgetKit |
| データ共有（App ↔ Widget） | App Groups + UserDefaults |
| ネットワーク | URLSession（サードパーティ不使用） |
| 状態管理 | `@Observable` (Swift 5.9+) |

---

## アーキテクチャ

```
起動
  ↓
Google Sign-In（未認証時のみ）
  ↓
Google Calendar API
  → 翌日の最初の予定時刻を取得（例: 09:00）
  ↓
バックエンド GET /summary?next_event=09:00
  ↓
SummaryResponse（condition_score, focus_peak_start/end, recommend_bedtime, sleep_debt_minutes）
  ↓
┌─────────────────────────────────────┐
│  DashboardView（ダッシュボード画面）  │
└─────────────────────────────────────┘
  ↓
UserDefaults（App Group）へ保存
  ↓
┌──────────────────┐   ローカル通知スケジュール
│  RignalWidget      │   （推奨就寝時刻の 30 分前）
└──────────────────┘
```

---

## ディレクトリ構成

```
frontend/
├── Project.swift                        # Tuist プロジェクト定義（Target 追加時はここを編集）
├── Tuist/
│   └── Package.swift                    # 外部依存パッケージ（Google SDK など）
│
├── RignalApp/
│   ├── Sources/
│   │   ├── RignalApp.swift          # エントリーポイント
│   │   │
│   │   ├── Views/
│   │   │   └── DashboardView.swift      # ➕ メイン画面
│   │   │
│   │   ├── ViewModels/
│   │   │   └── DashboardViewModel.swift # ➕ @Observable。API 呼び出し・状態管理
│   │   │
│   │   ├── Models/
│   │   │   └── SummaryResponse.swift    # ➕ バックエンドレスポンスの Codable 型
│   │   │
│   │   ├── Services/
│   │   │   ├── BackendService.swift     # ➕ Go バックエンドへの URLSession ラッパー
│   │   │   ├── CalendarService.swift    # ➕ Google Calendar API 呼び出し
│   │   │   └── NotificationService.swift# ➕ ローカル通知スケジュール管理
│   │   │
│   │   └── Auth/
│   │       └── GoogleAuthManager.swift  # ➕ Google Sign-In の認証状態管理
│   │
│   ├── Resources/
│   │   └── Assets.xcassets/
│   │
│   └── Tests/
│       └── RignalAppTests.swift
│
└── RignalWidget/                          # ➕ WidgetKit Extension
    ├── Sources/
    │   ├── RignalWidgetBundle.swift
    │   ├── RignalWidget.swift             # ➕ Widget 定義（TimelineProvider）
    │   └── WidgetEntryView.swift        # ➕ Widget UI
    └── Resources/
```

> `➕` は追加ファイルです。

---

## 主要画面：DashboardView

表示する情報:

| 項目 | 内容 |
|---|---|
| コンディションスコア | 0–100 の数値 + 段階ラベル（Good / Fair / Low） |
| 集中ピーク時間帯 | 例「09:30 〜 12:30」 |
| 推奨就寝時刻 | 例「23:00」（翌日予定がない場合は非表示） |
| 睡眠負債 | 例「45 分の借金あり」 |
| 翌日の最初の予定 | Google カレンダーから取得した予定名 + 時刻 |

画面遷移は 1 画面のみ（ダッシュボードのみ）。設定（Oura トークン・Google アカウント）はOS の設定アプリへ誘導する形でも可。

---

## Google カレンダー連携（CalendarService.swift）

1. Google Sign-In で認証済みの `GIDGoogleUser` を取得
2. `GTLRCalendarService` を使って `calendarId = "primary"` の翌日イベントを検索
3. 時刻でソートして最初のイベントの `start.dateTime` を取得
4. `HH:MM` 形式に変換して `BackendService` へ渡す

スコープ: `https://www.googleapis.com/auth/calendar.readonly`

---

## バックエンド通信（BackendService.swift）

```swift
// リクエスト例
GET {BACKEND_BASE_URL}/summary?next_event=09:00

// レスポンス型
struct SummaryResponse: Codable {
    let date: String
    let conditionScore: Int
    let focusPeakStart: Date?
    let focusPeakEnd: Date?
    let recommendBedtime: Date?
    let sleepDebtMinutes: Int
}
```

- `BACKEND_BASE_URL` は `Info.plist` の `BackendBaseURL` キーから取得する
- 認証ヘッダーは現時点では不要（ローカルネットワーク内での利用を想定）

---

## ローカル通知（NotificationService.swift）

- サマリー取得後、`recommend_bedtime` の **30 分前** に就寝リマインド通知をスケジュール
- 既存の通知があれば上書きする（identifier: `"rignal.bedtime.reminder"`）
- 許可リクエストはアプリ初回起動時に行う

---

## ウィジェット（RignalWidget）

- App Group (`group.dev.tuist.RignalApp`) を使って本体アプリと `UserDefaults` を共有
- 本体アプリがサマリーを取得したら App Group の UserDefaults に保存
- `TimelineProvider` が UserDefaults から読み取って表示
- 更新タイミング: 本体アプリ起動時 + 毎朝 7:00（`.atEnd` ポリシー）
- サイズ: `.systemSmall` のみ対応

---

## App Groups 設定

- Identifier: `group.dev.tuist.RignalApp`
- 本体アプリと RignalWidget Extension の両方で有効化が必要
- `Project.swift` の entitlements に追記する

---

## 環境設定

`Info.plist` に以下のキーを追加:

```
BackendBaseURL = http://192.168.x.x:8081   （開発時: ローカル IP）
GIDClientID    = {GoogleのOAuthクライアントID}
```

Google Cloud Console で OAuth クライアントを作成し、`GIDClientID` を取得すること。
URL Scheme に `com.googleusercontent.apps.{クライアントID}` を登録すること。

---

## コーディング規約

- `@Observable` を使う（`ObservableObject` は使わない）
- ビジネスロジックは `Services/` に置く。View は表示のみ
- 非同期処理は `async/await` で統一（Combine は使わない）
- エラーは `throws` で上位に伝播させ、ViewModel でキャッチして表示
- `UserDefaults` への直接アクセスは `NotificationService` と Widget のみ許可

---

## 実装優先順位（Phase）

### Phase 1（表示まで動かす）
1. `BackendService.swift` — `/summary` を叩いてレスポンスを返す
2. `DashboardViewModel.swift` — BackendService を呼んで状態を保持
3. `DashboardView.swift` — スコア・ピーク・就寝時刻を表示

### Phase 2（Google Calendar 連携）
4. `GoogleAuthManager.swift` — Google Sign-In の認証フロー
5. `CalendarService.swift` — 翌日の最初の予定時刻を取得して ViewModel へ渡す

### Phase 3（通知・ウィジェット）
6. `NotificationService.swift` — 就寝リマインド通知のスケジュール
7. `RignalWidget` — App Group 経由でデータを読み取り表示

---

## 作業時の注意

- Tuist でターゲットを追加する場合は必ず `Project.swift` を編集してから `tuist generate` を実行する
- Google Sign-In の URL Scheme 登録を忘れると認証コールバックが返ってこない
- Widget Extension は独立したターゲットのため、App Group を両方に設定しないとデータが共有されない
- ローカル通知の許可が得られていない場合は通知をスケジュールせずにサイレントに失敗させる
