# btc_trade

## 概要
ビットコインの自動トレーディングアプリ。8種類のインディケーター(SMA,EMA,BBands,Ichimoku,Volume,Rsi,MACD,HV)を実装。ローソク足は1分、1時間、1日単位に対応。
バックテストモードではそれぞれのインディケーターのパラメータを設定でき、イベント表示させることで実際には取引を行うことはせずにトレードのテストが行える。
トレードモードでは保持コイン数に対しての使用可能率や取引を停止させる割合も設定でき、設定期間内で利益を出せるインディケーターの上位のものを使用してトレードする。

## 使用技術
### 言語
- Golang 1.15
- JaveScript
- HTML
- CSS
### デプロイ環境
- GCP(GAE、Mysql 8.0)

### 外部API、サービス
- bitflyerAPI
- GoogleChart

## 工夫点
- docker,docker composeによるローカル開発環境の構築
- MarshallJSON実装
- DateFrameの実装
- バックテストの実装
- ホットリロードの導入
- goroutineによるリアルタイムでのビットコイン情報の取得とトレード・分析の並列処理
- testdateディレクトリ内のjsonデータによるユニットテスト
- RoundTripの実装によるHTTP Clientのテスト
