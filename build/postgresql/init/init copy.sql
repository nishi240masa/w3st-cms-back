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

-- apiSchema テーブル
CREATE TABLE IF NOT EXISTS apiSchema (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    field_id VARCHAR(100) NOT NULL, -- フィールドID
    view_name VARCHAR(100) NOT NULL, -- 表示名
    field_type VARCHAR(50) NOT NULL, -- フィールドの型 ('text', 'number', 'boolean', etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- field_data テーブル
CREATE TABLE IF NOT EXISTS field_data (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    field_type VARCHAR(50) NOT NULL, -- データ型 ('text', 'number', etc.)
    field_value JSONB, -- データの値 (JSON形式で保存)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- list_options テーブル (選択肢)
CREATE TABLE IF NOT EXISTS list_options (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL, -- 選択肢の内容
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- api_kind_relation テーブル (リレーション情報)
CREATE TABLE IF NOT EXISTS api_kind_relation (
    id SERIAL PRIMARY KEY,
    apiSchema_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    related_id INT NOT NULL REFERENCES apiSchema(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL, -- リレーションの種類 ('parent', 'child', etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- スキマーの削除時に関連データを削除するトリガー関数
CREATE OR REPLACE FUNCTION cascade_delete_apiSchema()
RETURNS TRIGGER AS $$
BEGIN
    -- field_data を削除
    DELETE FROM field_data WHERE apiSchema_id = OLD.id;

    -- list_options を削除
    DELETE FROM list_options WHERE apiSchema_id = OLD.id;

    -- api_kind_relation の関連レコードを削除
    DELETE FROM api_kind_relation WHERE apiSchema_id = OLD.id OR related_id = OLD.id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- トリガーの設定
CREATE TRIGGER delete_related_data
AFTER DELETE ON apiSchema
FOR EACH ROW
EXECUTE FUNCTION cascade_delete_apiSchema();


-- list_options テーブルに対するトリガー関数
CREATE OR REPLACE FUNCTION validate_list_options()
RETURNS TRIGGER AS $$
DECLARE
    schema_type VARCHAR(50);
BEGIN
    -- 対応する apiSchema の field_type を取得
    SELECT field_type INTO schema_type FROM apiSchema WHERE id = NEW.apiSchema_id;

    -- field_type が 'select' または 'dropdown' でなければエラーを返す
    IF schema_type NOT IN ('select', 'dropdown') THEN
        RAISE EXCEPTION 'list_options can only be added to select or dropdown fields';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーの設定
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
    -- 循環チェック (簡略化した例)
    WITH RECURSIVE relation_path AS (
        SELECT related_id
        FROM api_kind_relation
        WHERE apiSchema_id = NEW.related_id
        UNION ALL
        SELECT r.related_id
        FROM api_kind_relation r
        INNER JOIN relation_path rp ON rp.related_id = r.apiSchema_id
    )
    SELECT EXISTS (
        SELECT 1 FROM relation_path WHERE related_id = NEW.apiSchema_id
    ) INTO is_cyclic;

    -- 循環が検出された場合はエラーをスロー
    IF is_cyclic THEN
        RAISE EXCEPTION 'Cyclic relation detected';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーの設定
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
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON apiSchema
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON field_data
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON list_options
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON api_kind_relation
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();