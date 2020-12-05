# AtenaOCIOracle
# 宛名の管理プログラム

- 選択、送受信履歴をとりあえず保存できるようにした物
- SQLDeveloper や TablePlusで編集したり、宛名職人からの書き出しを想定しているため編集機能すらない

- 前もってsample.sql.txt を読み込ませて ATENA テーブルと、ATENAUSER テーブルのデータ構造を作ってください
  - `sqlplus admin/hogehoge@hogedb_tp`
  - `@sample.sql.txt`

- ビルド
  - `go build -o AtenaOCIOracle`

- 実行
  - `./AtenaOCIOracle &` でサーバが起動
  - 環境変数 OCISTRING に "admin/hogehoge@hogedb_tp" などを指定すること
  - 環境変数 PORT が指定されていないと 3002を使います
  - webブラウザで http://host:3002/ にアクセス

- 機能
  - login 画面で USER: test PASS: test でログインすると、アドレス帳一覧が表示
  - パスワードは変更可能だが現在はマルチユーザーに対応していない
  - 左側の選択チェックボックスで選択かどうかを選べます
  - 選択した合計数が上に表示
  - 今年・去年の送受・喪中がプルダウンメニューで選べます
  - 上の「選択・送受を確定」ボタンで DBに書き込みます
  - 「選択」「名前」をクリックすると、選択優先・名前優先でソートされます
  - ソート情報、ログイン情報は cookie に保存される
  - CRUDはまだつけてない
  - NYCARDHISTORY には変換して使っていない場合 "00000000000000002005"を入れていないと boundary errorが出ますので入れておけ
  - サンプルデータ作成：https://tm-webtools.com/Tools/TestData

- 今後
  - 年度rotateボタンの作成
  - 簡単なCRUD
  - CSVの書き出し機能
