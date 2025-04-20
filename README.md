windowsのみ

GOOS=windows GOARCH=amd64 go build -o ssh.exe
でビルドする

ssh.exeをデスクトップに配置して、管理者権限で実行する。

Desktopにkey.csvを作成

プロジェクト名 ユーザー名 IP PassWord keyのpathを入れる

keyはDocumentsディレクトリに入れておく。
