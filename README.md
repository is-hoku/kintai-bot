# kintai-bot

## 背景
僕のインターン先では、勤務 (休憩) の開始、終了時に Slack で報告し、1日の勤務終了時などに freee 人事労務で打刻する必要がある。つまり僕は毎回の勤務で、複数回同じことを入力しており、これはブラウザを起動してサービスにログインする作業も含めると、それなりに面倒になっていた。
  
一回の入力で全ての勤怠入力をするために、この Slack App を作成した。

## アーキテクチャ
Slack からのスラッシュコマンドを受信する Bolt Server (https://github.com/is-hoku/kintai-bot-bolt) 、 freee API を叩く API Client と、ユーザの Slack に登録してある Email と freee ID、また Company ID と Token を紐付けるための DB, API Server から構成される。

![kintai-bot-figma](https://user-images.githubusercontent.com/52068717/151113691-66bc4745-b2cb-47b4-a9ef-3cc4c333dac4.png)

## 機能
Slack App に実装されている Slash Command は以下。
- `/auth`
- `/clock_in`
- `/clock_out`
- `/break_begin`
- `/break_out`

![2022-01-26_12-12](https://user-images.githubusercontent.com/52068717/151122877-51f42ae3-cfec-41ef-89f0-a1bfa8f5f160.png)

### `/auth`
freee API を叩いて認可URLを返す。ブラウザで認可すると freee API から生成されたアクセストークン、リフレッシュトークンを DB に保存する。  
実行ははじめの一回のみで良い。認可は freee の管理者権限を持つアカウントで行う必要がある。

### `/clock_in`
freee API の `/api/v1/employees/{emp_id}/time_clocks` を叩く。  
`type` は `clock_in`、 `datetime` は現在時刻。 `basetime` は考慮しない。つまり日をまたいでの勤務は想定していない。  
また、 freee API で不正な勤怠登録とされる操作 (勤務をはじめる前に退勤するなど) はエラーを返す。  
API を叩いて打刻した後、 Slash Command を実行したアカウントとしてメッセージ (「稼動します!」など) をチャンネルに投稿する。  

実行後の Slack の例  
![2022-01-26_12-15](https://user-images.githubusercontent.com/52068717/151122774-bc0af799-c335-4a96-ae20-d28cb3110497.png)  
実行後の freee の例  
![2022-01-26_12-15_1](https://user-images.githubusercontent.com/52068717/151127919-995aab64-2f1f-4ec8-a95a-d9e930fce2c0.png)


### `/clock_out`
退勤する。

### `/break_begin`
休憩を開始する。

### `/break_end`
休憩を終わる。

## 準備
以下の事項は実装されていないため、この Slack App を使用する前に手動でセットする必要がある。  
- DB へユーザの emp_id (freee_id) とメールアドレス (email) を登録
- company_id と company_name の取得と環境変数へのセット
- Slack のシークレットトークンなどの環境変数へのセット
- freee App の Client ID と Client Secret の取得と環境変数へのセット
