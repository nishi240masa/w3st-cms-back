
# w3st CMS API仕様書（β版）


## 目次


- [概要](#概要)

- [対象ユーザー](#対象ユーザー)

- [テーブル設計](#テーブル設計)

- [主要機能](#主要機能)

- [API設計](#api設計)

- [差別化ポイント](#差別化ポイントmicrocmsなどと比較)



---



## 概要

**w3st CMS API** は、フロントエンド開発者向けに、直感的なAPI設計と自由度の高いスキーマ定義を提供するヘッドレスCMSです。

GUIベースでAPIを生成でき、通常のCMSよりもAPIベース開発に最適化されています。


---



## 対象ユーザー


- フロントエンドエンジニア

- Jamstack開発者

- Webアプリ開発者



---



## テーブル設計


### users

| カラム名       | 型            | 説明                       | 
|------------|--------------|--------------------------| 
| id         | UUID         | ユーザーID                   | 
| name       | VARCHAR(100) | 名前                       | 
| email      | VARCHAR(255) | メールアドレス (一意)             | 
| password   | VARCHAR(100) | ハッシュ化されたパスワード            | 
| role       | VARCHAR(50)  | ユーザーロール (例: user, admin) | 
| created_at | TIMESTAMP    | 作成日時                     | 
| updated_at | TIMESTAMP    | 更新日時                     | 



---



### api_collections


コレクション（スキーマ）を管理

| カラム名        | 型            | 説明        | 
|-------------|--------------|-----------| 
| id          | SERIAL       | コレクションID  | 
| user_id     | UUID         | 所有者ユーザーID | 
| name        | VARCHAR(100) | コレクション名   | 
| description | TEXT         | 説明        | 
| created_at  | TIMESTAMP    | 作成日時      | 
| updated_at  | TIMESTAMP    | 更新日時      | 



---



### api_fields


フィールド定義

| カラム名          | 型            | 説明                                        | 
|---------------|--------------|-------------------------------------------| 
| id            | SERIAL       | フィールドID                                   | 
| collection_id | INT          | 紐づくコレクションID                               | 
| field_id      | VARCHAR(100) | 内部フィールドキー                                 | 
| view_name     | VARCHAR(100) | 表示用ラベル                                    | 
| field_type    | VARCHAR(50)  | 型 (text, number, boolean, relation, etc.) | 
| is_required   | BOOLEAN      | 必須フラグ                                     | 
| default_value | JSONB        | デフォルト値                                    | 
| created_at    | TIMESTAMP    | 作成日時                                      | 
| updated_at    | TIMESTAMP    | 更新日時                                      | 



---



### content_entries


コンテンツデータ本体

| カラム名          | 型         | 説明       | 
|---------------|-----------|----------| 
| id            | SERIAL    | エントリID   | 
| collection_id | INT       | コレクションID | 
| data          | JSONB     | データ本体    | 
| created_at    | TIMESTAMP | 作成日時     | 
| updated_at    | TIMESTAMP | 更新日時     | 



---



### list_options


選択肢フィールド専用

| カラム名       | 型            | 説明         | 
|------------|--------------|------------| 
| id         | SERIAL       | 選択肢ID      | 
| field_id   | INT          | 紐づくフィールドID | 
| value      | VARCHAR(255) | 選択肢の値      | 
| created_at | TIMESTAMP    | 作成日時       | 
| updated_at | TIMESTAMP    | 更新日時       | 



---



### api_kind_relation


コレクション間リレーション管理

| カラム名                  | 型           | 説明                                  | 
|-----------------------|-------------|-------------------------------------| 
| id                    | SERIAL      | リレーションID                            | 
| collection_id         | INT         | 元コレクションID                           | 
| related_collection_id | INT         | 関連コレクションID                          | 
| relation_type         | VARCHAR(50) | リレーション型 (one-to-many, many-to-many) | 
| created_at            | TIMESTAMP   | 作成日時                                | 
| updated_at            | TIMESTAMP   | 更新日時                                | 



---



### api_keys 

APIキー管理（個別にキーを発行してアクセスコントロールする）

| カラム名                | 型            | 説明                     |
|---------------------|--------------|------------------------| 
| id                  | SERIAL       | APIキーID                |
| user_id             | UUID         | キーの所有者ユーザーID           |
| name                | VARCHAR(100) | キーの名前（管理用）             |
| key                 | VARCHAR(255) | 実際に発行されたAPIキー文字列（ユニーク） |
| ip_whitelist        | TEXT[]       | 許可されたIPリスト（空なら無制限）     |
| expire_at           | TIMESTAMP    | 有効期限（NULLなら無期限）        |
| revoked             | BOOLEAN      | 無効化されているか              |
| rate_limit_per_hour | INT          | 1時間あたりのリクエスト上限         |
| created_at          | TIMESTAMP    | 作成日時                   |

### api_key_collections

APIキー単位でアクセス許可されるコレクションを紐付ける中間テーブル

| カラム名                | 型            | 説明                     |
|---------------------|--------------|------------------------| 
| id                  | SERIAL       | APIキーID                |
| user_id             | UUID         | キーの所有者ユーザーID           |
| name                | VARCHAR(100) | キーの名前（管理用）             |
| key                 | VARCHAR(255) | 実際に発行されたAPIキー文字列（ユニーク） |
| ip_whitelist        | TEXT[]       | 許可されたIPリスト（空なら無制限）     |
| expire_at           | TIMESTAMP    | 有効期限（NULLなら無期限）        |
| revoked             | BOOLEAN      | 無効化されているか              |
| rate_limit_per_hour | INT          | 1時間あたりのリクエスト上限         |
| created_at          | TIMESTAMP    | 作成日時                   |



---



## 主要機能


- コレクション作成・編集・削除

- フィールド定義・編集

- コンテンツ（エントリ）CRUD

- リレーション管理（循環検出あり）

- APIキー管理（IP制限・レート制限）

- オプション選択フィールドサポート

- 自動タイムスタンプ更新



---



## API設計


### コレクション作成 API（tx使用）

**POST**  `/api/collections`
**Request Body**


```json
{
  "user_id": "UUID",
  "name": "コレクション名",
  "description": "説明",
  "fields": [
    {
      "field_id": "product_name",
      "view_name": "商品名",
      "field_type": "text"
    },
    {
      "field_id": "price",
      "view_name": "価格",
      "field_type": "number"
    }
  ]
}
```

**Response**


```json
{
  "collection_id": 1,
  "message": "Collection created successfully"
}
```


※ この一連の処理は**トランザクション(tx)**でまとめる



---




## 差別化ポイント（microCMSなどと比較）


- フィールド単位でリレーションが設定できる

- スキーマ循環防止が組み込まれている

- IP制限・レート制限付きAPIキー発行

- コレクション単位でアクセス制御可能

- 複雑なスキーマ構成（JSONBによる柔軟なデータ格納）

- 高速なカスタマイズ性（GUI設計予定）

- ローカル開発モード（オフラインサポート）を計画中



---



# 🚀 今後追加予定（Future Work）


- APIログ管理（リクエストログ）

- Webhook通知

- バージョニング管理

- マルチユーザー共同編集機能

- SSO (Google, GitHub認証)



---
