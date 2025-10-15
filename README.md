# w3st CMS API仕様書（β版）

## 目次

- [概要](#概要)
- [対象ユーザー](#対象ユーザー)
- [テーブル設計](#テーブル設計)
- [主要機能](#主要機能)
- [APIの使い方](#apiの使い方)
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

### projects

プロジェクト管理

| カラム名             | 型            | 説明         |
|--------------------|--------------|------------|
| id                 | SERIAL       | プロジェクトID |
| name               | VARCHAR(100) | プロジェクト名 |
| description        | TEXT         | 説明         |
| rate_limit_per_hour| INT          | 1時間あたりのレート制限 |
| created_at         | TIMESTAMP    | 作成日時     |
| updated_at         | TIMESTAMP    | 更新日時     |

---

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

| カラム名        | 型     | 説明         |
|-------------|-------|------------|
| api_key_id  | INT   | APIキーID    |
| collection_id | INT   | コレクションID |

---

### audit_logs

監査ログを記録

| カラム名     | 型            | 説明       |
|----------|--------------|----------|
| id       | UUID         | ログID     |
| user_id  | UUID         | ユーザーID   |
| action   | VARCHAR(50)  | アクション    |
| resource | VARCHAR(255) | リソース     |
| created_at | TIMESTAMP    | 作成日時     |
| details  | TEXT         | 詳細       |

---

### media_assets

メディアアセット管理

| カラム名     | 型            | 説明     |
|----------|--------------|--------|
| id       | UUID         | メディアID |
| name     | VARCHAR(255) | 名前     |
| type     | VARCHAR(50)  | タイプ    |
| path     | TEXT         | パス     |
| size     | BIGINT       | サイズ    |
| user_id  | UUID         | ユーザーID |
| created_at | TIMESTAMP    | 作成日時   |
| updated_at | TIMESTAMP    | 更新日時   |

---

### user_permissions

ユーザー権限管理

| カラム名      | 型            | 説明     |
|-----------|--------------|--------|
| id        | UUID         | 権限ID   |
| user_id   | UUID         | ユーザーID |
| permission | VARCHAR(50)  | 権限     |
| resource  | VARCHAR(255) | リソース   |
| created_at | TIMESTAMP    | 作成日時   |
| updated_at | TIMESTAMP    | 更新日時   |

---

### content_versions

コンテンツバージョン管理

| カラム名      | 型         | 説明       |
|-----------|-----------|----------|
| id        | UUID      | バージョンID |
| content_id | UUID      | コンテンツID |
| version   | INT       | バージョン番号 |
| data      | JSONB     | データ      |
| user_id   | UUID      | ユーザーID   |
| created_at | TIMESTAMP | 作成日時     |
| updated_at | TIMESTAMP | 更新日時     |

---

### system_alerts

システムアラート管理

| カラム名     | 型            | 説明         |
|----------|--------------|------------|
| id        | SERIAL       | アラートID   |
| alert_type| VARCHAR(50)  | アラートタイプ |
| severity  | VARCHAR(50)  | 深刻度       |
| title     | VARCHAR(255) | タイトル     |
| message   | TEXT         | メッセージ   |
| project_id| INT          | プロジェクトID |
| is_active | BOOLEAN      | アクティブか   |
| is_read   | BOOLEAN      | 既読か       |
| metadata  | JSONB        | メタデータ   |
| created_at| TIMESTAMP    | 作成日時     |
| updated_at| TIMESTAMP    | 更新日時     |

---

## 主要機能

- プロジェクト管理
- コレクション作成・編集・削除
- フィールド定義・編集
- コンテンツ（エントリ）CRUD
- リレーション管理（循環検出あり）
- APIキー管理（IP制限・レート制限）
- オプション選択フィールドサポート
- 自動タイムスタンプ更新
- メディアアセット管理
- ユーザー権限管理
- コンテンツバージョン管理
- 監査ログ記録
- システムアラート管理

---

## APIの使い方

このセクションでは、w3st CMS APIの基本的な使い方をステップバイステップで説明します。すべてのAPIリクエストには適切な認証が必要です。

### 1. ユーザー登録とログイン

まず、ユーザーとして登録し、JWTトークンを取得します。

#### ユーザー登録
```bash
POST /users/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### ログイン
```bash
POST /users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

レスポンスからJWTトークンを取得し、以後のリクエストのAuthorizationヘッダーに `Bearer <token>` を設定してください。

### 2. プロジェクトの作成

プロジェクトを作成します。

```bash
POST /api/projects
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "My Project",
  "description": "Description of the project",
  "rate_limit_per_hour": 1000
}
```

### 3. コレクションの作成

コンテンツを管理するためのコレクション（スキーマ）を作成します。

```bash
POST /api/collections
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "Products",
  "description": "Product catalog"
}
```

### 4. フィールドの追加

作成したコレクションにフィールドを定義します。

```bash
POST /api/collections/{collectionId}/fields
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "field_id": "name",
  "view_name": "Product Name",
  "field_type": "text",
  "is_required": true
}
```

### 5. エントリの作成と管理

コレクションにコンテンツエントリを追加します。

#### エントリ作成
```bash
POST /api/collections/{collectionId}/entries
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "data": {
    "name": "Sample Product"
  }
}
```

#### エントリ取得（SDK API）
```bash
GET /collections/{collectionId}/entries
X-API-Key: <your-api-key>
```

### 6. APIキーの発行

公開APIアクセス用のAPIキーを作成します。

```bash
POST /api/api-keys
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "Public API Key",
  "collection_ids": [1, 2],
  "ip_whitelist": ["192.168.1.1"],
  "rate_limit_per_hour": 1000
}
```

### 7. メディアアセットの管理

画像などのメディアファイルをアップロードします。

#### アップロード
```bash
POST /api/media
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "product-image.jpg",
  "type": "image/jpeg",
  "path": "/uploads/product-image.jpg",
  "size": 1024000
}
```

#### 一覧取得
```bash
GET /api/media
Authorization: Bearer <your-jwt-token>
```

### 8. 権限管理

ユーザーの権限を管理します。

#### 権限付与
```bash
POST /api/permissions/grant
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "permission": "read",
  "resource": "collection:1"
}
```

#### 権限チェック
```bash
GET /api/permissions/check?permission=read&resource=collection:1
Authorization: Bearer <your-jwt-token>
```

### 9. コンテンツバージョン管理

エントリのバージョンを管理します。

#### バージョン作成
```bash
POST /api/versions
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "content_id": "uuid-of-content",
  "data": {"name": "Updated Product"}
}
```

#### バージョン一覧
```bash
GET /api/versions/{contentID}
Authorization: Bearer <your-jwt-token>
```

#### バージョン復元
```bash
POST /api/versions/{contentID}/restore/{versionID}
Authorization: Bearer <your-jwt-token>
```

### 10. 監査ログの確認

システムのアクティビティログを確認します。

#### ログ記録
```bash
POST /api/audit
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "action": "create",
  "resource": "collection:1",
  "details": "Created new collection"
}
```

#### ログ取得
```bash
GET /api/audit/user
Authorization: Bearer <your-jwt-token>
```

### 11. システムアラートの管理

システムアラートを管理します。

#### アラート作成
```bash
POST /api/system-alerts
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "alert_type": "info",
  "severity": "low",
  "title": "System Update",
  "message": "System will be updated tonight",
  "project_id": 1
}
```

#### アクティブアラート取得
```bash
GET /api/system-alerts/active
Authorization: Bearer <your-jwt-token>
```

詳細なAPI仕様については `api-document.yaml` を参照してください。

---

## API設計

詳細なAPI仕様は `api-document.yaml` を参照してください。

---

## 差別化ポイント（microCMSなどと比較）

- フィールド単位でリレーションが設定できる
- スキーマ循環防止が組み込まれている
- IP制限・レート制限付きAPIキー発行
- コレクション単位でアクセス制御可能
- 複雑なスキーマ構成（JSONBによる柔軟なデータ格納）
- 高速なカスタマイズ性（GUI設計予定）
- プロジェクトベースの管理
- ローカル開発モード（オフラインサポート）を計画中
- メディアアセット管理機能
- 詳細な権限管理システム
- コンテンツバージョン管理
- 監査ログによるセキュリティ強化
- システムアラート管理

## 認証フロー

```flowchart TD
    A[リクエスト受信] --> B{認証方式?}
    B -->|JWT| C[JwtAuthMiddleware]
    B -->|APIキー| D[ApiKeyAuthMiddleware]
    B -->|Auth0| E[Auth0AuthMiddleware]
    
    C --> F[ValidateToken]
    F --> G{有効?}
    G -->|Yes| H[コンテキストにuserIDセット]
    G -->|No| I[401エラー]
    
    D --> J[ValidateApiKey]
    J --> K{有効?}
    K -->|Yes| L[JWTトークン生成]
    K -->|No| I
    L --> M[JWTパース検証]
    M --> N{有効?}
    N -->|Yes| O[コンテキストにuserID/projectID/collectionIdsセット]
    N -->|No| I
    
    E --> P[Auth0トークン検証]
    P --> Q{有効?}
    Q -->|Yes| R[コンテキストにuserID/email/nameセット]
    Q -->|No| I
    
    H --> S[レート制限チェック]
    O --> S
    R --> S
    
    S --> T{制限内?}
    T -->|Yes| U[次のハンドラー]
    T -->|No| V[429エラー]

```

---

# 🚀 今後追加予定（Future Work）

- APIログ管理（リクエストログ）
- Webhook通知
- バージョニング管理
- マルチユーザー共同編集機能

---
