# kintai-bot

## 背景
僕のインターン先では、勤務 (休憩) の開始、終了時に Slack で報告し、1日の勤務終了時などに freee 人事労務で勤怠を登録する必要がある。つまり僕は毎回の勤務で、複数回同じことを入力しており、これはブラウザを起動してサービスにログインする作業も含めると、それなりに面倒になっていた。  
Slack から freee に打刻することは freee 公式の [Slack App](https://support.freee.co.jp/hc/ja/articles/360016610812) を使用することで実現できるが、「freee人事労務」アプリケーションとのチャットウィンドウでコマンドを実行する必要がある。これでは、勤怠を社内の人に報告すると同時に勤怠管理をしている freee に打刻するという欲求が満たされない。また、 Webhook を用いて freee に打刻したら Slack に勤怠を投稿するという仕様も考えられるが、ブラウザを起動して freee にログインするのは Slack を起動して投稿することに比べてやや面倒な点、そもそも freee 人事労務に Webhook 機能が実装されていない点から、 Slack に Slash Command で勤怠を投稿すると freee に打刻されるという仕様に落ち着いた。  

## アーキテクチャ
Slack からのスラッシュコマンドを受信する Bolt Server ([is-hoku/kintai-bot-bolt](https://github.com/is-hoku/kintai-bot-bolt)) / API Server (freee API を叩く API (OAuth) Client) / DB Server から構成される。

![kintai-bot](https://user-images.githubusercontent.com/52068717/155745920-673ecd4b-512a-42b5-a57a-19a5c8f69b7c.png)

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

## freee アプリの権限

| 権限 | 参照 |更新|
|:----:|:----:|:--:|
| [人事労務] 従業員 | <li>[x] </li> | <li>[ ] </li> |
| [人事労務] 打刻 | <li>[x] </li> | <li>[x] </li> |

## 準備
はじめにやることです。  
freee API から company_id と company_name, Client ID とClient Secret を取得して環境変数に入れる。  
Slack App のトークンなどを取得して環境変数に入れる。  
db コンテナに入って以下を実行してトークンの Field を作成する。  
```
# mongo
> use tokens
> db.tokens.insert({company_id:COMPANY_ID,access_token:"",refresh_token:"",expiry:"2020-12-02T15:04:05+09:00"})
```
`/user` に以下のようなボディで POST してユーザを作成。  
```
{
  "email": "korehatestuserdesu@gmail.com",
  "freee_id": USER_ID
}
```

## Inspirations and References
- [igsr5/time-management-go](https://github.com/igsr5/time-management-go)
- [igsr5/time-management-bolt](https://github.com/igsr5/time-management-bolt)
