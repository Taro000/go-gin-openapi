# go-gin-webapi

アカウント登録機能が付いたTODOリストを作る。

## 技術スタック

- Go
- Gin
- MySQL
- Firebase Authentication
- OpenAPI

## 機能一覧

- 会員登録
- ログイン
- ログアウト
- ユーザー詳細取得
- ユーザー情報編集
- Todo作成
- Todo編集
- Todo削除
- Todo詳細取得
- Todo一覧取得
- いいね作成
- いいね削除

## ER 図

```mermaid
erDiagram
    User {
        CHAR(28) uid PK "ユーザーID"
        VARCHAER(15) nickname "ニックネーム"
        VARCHAR(255) email "メールアドレス"
    }

    TodoStatus {
        status CHAR(2) PK "ステータス"
    }

    Todo {
        CHAR(36) id PK "TodoID"
        CHAR(28) owner FK "所有ユーザー"
        CHAR(2) status FK "ステータス"
        VARCHAER(30) title "タイトル"
        TEXT content "内容"
        DATETIME due_datetime "期限日時"
        DATETIME created_at "作成日時"
        DATETIME updated_at "更新日時"
    }

    Goodluck {
        CHAR(28) user FK "ユーザー"
        CHAR(36) todo FK "Todo"
    }

    User ||--o{ Todo :"１人のユーザーは<br>N個のTodoを持てる。"
    TodoStatus ||--o{ Todo :"１つのステータスは<br>N個のTodoから設定され得る。"
    User ||--o{ Goodluck :"1人のユーザーは<br>N回いいねができる。"
    Todo ||--o{ Goodluck :"1個のTodoは<br>N回いいねをされ得る。"
```
