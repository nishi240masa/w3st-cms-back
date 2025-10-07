-- users テーブル
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- api_collections: スキーマ（エンティティ）を定義
CREATE TABLE api_collections (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INT NOT NULL, -- プロジェクトID
    name VARCHAR(100) NOT NULL, -- コレクション名 ex) 'ユーザー', '商品'
    description TEXT, -- 説明 ex) 'ユーザー情報を管理するコレクション'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- api_fields: 各スキーマのフィールド定義
CREATE TABLE api_fields (
    id SERIAL PRIMARY KEY,
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    field_id VARCHAR(100) NOT NULL, -- 内部的なキー ex) 'user_id', 'product_name'
    view_name VARCHAR(100) NOT NULL, -- 表示名 ex) 'ユーザーID', '商品名'
    field_type VARCHAR(50) NOT NULL, -- フィールドの型 ('text', 'number', 'boolean', etc.)
    is_required BOOLEAN DEFAULT false, -- 必須かどうか
    default_value JSONB, -- デフォルト値 (JSON形式で保存)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- field_data: 各スキーマのフィールド定義 (互換性用)
CREATE TABLE field_data (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL, -- プロジェクトID
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    field_id VARCHAR(100) NOT NULL, -- 内部的なキー ex) 'user_id', 'product_name'
    view_name VARCHAR(100) NOT NULL, -- 表示名 ex) 'ユーザーID', '商品名'
    field_type VARCHAR(50) NOT NULL, -- フィールドの型 ('text', 'number', 'boolean', etc.)
    is_required BOOLEAN DEFAULT false, -- 必須かどうか
    default_value JSONB, -- デフォルト値 (JSON形式で保存)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- content_entries: 実際のレコード（エントリ）を管理
CREATE TABLE content_entries (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL, -- プロジェクトID
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    data JSONB NOT NULL, -- エントリのデータ (JSON形式で保存)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- list_options テーブル (選択肢)
CREATE TABLE IF NOT EXISTS list_options (
    id SERIAL PRIMARY KEY,
    field_id INT NOT NULL REFERENCES api_fields(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL, -- 選択肢の値 ex) 'オプション1', 'オプション2'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS entries (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL, -- プロジェクトID
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- api_kind_relation テーブル (リレーション情報)
CREATE TABLE IF NOT EXISTS api_kind_relation (
    id SERIAL PRIMARY KEY,
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    related_collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL, -- リレーションの種類 ex) 'one-to-many', 'many-to-many'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INT NOT NULL,                         -- プロジェクトID
    name VARCHAR(100),                              -- 任意の名前（管理用）
    key VARCHAR(255) UNIQUE NOT NULL,               -- 発行されたAPIキー文字列
    collection_ids INT[] NOT NULL DEFAULT '{}',     -- アクセス可能なコレクションIDリスト
    ip_whitelist TEXT[],                            -- 許可されたIPアドレス（空配列は無制限）
    expire_at TIMESTAMP,                            -- 有効期限（NULLなら無期限）
    revoked BOOLEAN DEFAULT FALSE,                  -- 無効化フラグ
    rate_limit_per_hour INT DEFAULT 1000,           -- 1時間あたりの最大リクエスト数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE api_key_collections (
    api_key_id INT NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    collection_id INT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    PRIMARY KEY (api_key_id, collection_id)
);


-- スキマーの削除時に関連データを削除するトリガー関数
CREATE OR REPLACE FUNCTION cascade_delete_apiSchema()
RETURNS TRIGGER AS $$
BEGIN
    -- api_fields を削除
    DELETE FROM api_fields WHERE collection_id = OLD.id;

    -- list_options を削除
    DELETE FROM list_options WHERE field_id IN (SELECT id FROM api_fields WHERE collection_id = OLD.id);

    -- api_kind_relation の関連レコードを削除
    DELETE FROM api_kind_relation WHERE collection_id = OLD.id OR related_collection_id = OLD.id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- トリガーの設定
CREATE TRIGGER delete_related_data
AFTER DELETE ON api_collections
FOR EACH ROW
EXECUTE FUNCTION cascade_delete_apiSchema();


-- list_options テーブルに対するトリガー関数
CREATE OR REPLACE FUNCTION validate_list_options()
RETURNS TRIGGER AS $$
DECLARE
    f_type VARCHAR(50);
BEGIN
    SELECT field_type INTO f_type FROM api_fields WHERE id = NEW.field_id;
    IF f_type NOT IN ('select', 'dropdown') THEN
        RAISE EXCEPTION 'list_options can only be added to select or dropdown fields';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガー設定
CREATE TRIGGER validate_options
BEFORE INSERT ON list_options
FOR EACH ROW
EXECUTE FUNCTION validate_list_options();


-- relationが循環していないかチェックするトリガー関数
CREATE OR REPLACE FUNCTION check_cyclic_relation()
RETURNS TRIGGER AS $$
DECLARE
    is_cyclic BOOLEAN;
BEGIN
    WITH RECURSIVE relation_path AS (
        SELECT related_collection_id
        FROM api_kind_relation
        WHERE collection_id = NEW.related_collection_id
        UNION ALL
        SELECT r.related_collection_id
        FROM api_kind_relation r
        INNER JOIN relation_path rp ON rp.related_collection_id = r.collection_id
    )
    SELECT EXISTS (
        SELECT 1 FROM relation_path WHERE related_collection_id = NEW.collection_id
    ) INTO is_cyclic;

    IF is_cyclic THEN
        RAISE EXCEPTION 'Cyclic relation detected';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガー設定
CREATE TRIGGER prevent_cyclic_relation
BEFORE INSERT ON api_kind_relation
FOR EACH ROW
EXECUTE FUNCTION check_cyclic_relation();



-- トリガー関数の作成
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 各テーブルにトリガーを設定
CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_api_collections
BEFORE UPDATE ON api_collections
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_content_entries
BEFORE UPDATE ON content_entries
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_api_fields
BEFORE UPDATE ON api_fields
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_list_options
BEFORE UPDATE ON list_options
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp_api_kind_relation
BEFORE UPDATE ON api_kind_relation
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- media_assets テーブル
CREATE TABLE IF NOT EXISTS media_assets (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100),
    size BIGINT,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- content_versions テーブル
CREATE TABLE IF NOT EXISTS content_versions (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    content_entry_id INT REFERENCES content_entries(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    data JSONB NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (content_entry_id, version_number),
    CHECK (version_number > 0)
);

-- user_permissions テーブル
CREATE TABLE IF NOT EXISTS user_permissions (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_type VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- audit_logs テーブル
CREATE TABLE IF NOT EXISTS audit_logs (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    details JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- updated_at を持つ新しいテーブルにトリガーを設定
DROP TRIGGER IF EXISTS set_timestamp_media_assets ON media_assets;
CREATE TRIGGER set_timestamp_media_assets
BEFORE UPDATE ON media_assets
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

DROP TRIGGER IF EXISTS set_timestamp_content_versions ON content_versions;
CREATE TRIGGER set_timestamp_content_versions
BEFORE UPDATE ON content_versions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

DROP TRIGGER IF EXISTS set_timestamp_user_permissions ON user_permissions;
CREATE TRIGGER set_timestamp_user_permissions
BEFORE UPDATE ON user_permissions
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- インデックス追加
CREATE INDEX IF NOT EXISTS idx_content_versions_created_by ON content_versions(created_by);

-- user_permissions のユニーク制約（NULL セマンティクスを保持）
-- グローバル権限（resource_type, resource_id が両方 NULL）のユニーク制約
CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_permissions_global
  ON user_permissions(user_id, permission_type)
  WHERE resource_type IS NULL AND resource_id IS NULL;
-- スコープ付き権限（resource_type, resource_id が両方 NOT NULL）のユニーク制約
CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_permissions_scoped
  ON user_permissions(user_id, permission_type, resource_type, resource_id)
  WHERE resource_type IS NOT NULL AND resource_id IS NOT NULL;

-- user_permissions: resource_type と resource_id の整合性を保証する CHECK 制約
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'user_permissions_resource_scope_consistency'
      AND conrelid = 'user_permissions'::regclass
  ) THEN
    ALTER TABLE user_permissions
      ADD CONSTRAINT user_permissions_resource_scope_consistency
      CHECK (
        (resource_type IS NULL AND resource_id IS NULL) OR
        (resource_type IS NOT NULL AND resource_id IS NOT NULL)
      );
  END IF;
END
$$;

-- user_permissions の外部キー用インデックス
CREATE INDEX IF NOT EXISTS idx_user_permissions_user_id ON user_permissions(user_id);

-- audit_logs 検索用インデックス
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);