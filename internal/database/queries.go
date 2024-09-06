package database

const (
	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT UNIQUE,
    passwordhash TEXT);`

	CREATE_ITEM_TABLE_STMT = `CREATE TABLE IF NOT EXISTS item (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    picture TEXT,
    quantity INTEGER,
    weight TEXT,
    qrcode TEXT,
    box_id TEXT REFERENCES box(id));`

	CREATE_BOX_TABLE_STMT = `CREATE TABLE IF NOT EXISTS box (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL, 
    description TEXT,
    picture TEXT,
    qrcode TEXT,
    outerbox_id TEXT REFERENCES box(id));`

	CREATE_ITEM_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS item_fts USING fts5(
    item_id,
    label, 
    description, 
    tokenize = 'porter'
);`

	CREATE_BOX_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS box_fts USING fts5(
    box_id,
    label, 
    description,
    tokenize = 'porter');`

	// Triggers for item
	CREATE_ITEM_AI_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ai BEFORE INSERT ON item 
    BEGIN
        INSERT INTO item_fts(item_id, label, description) 
        VALUES (new.id, new.label, new.description);
    END; `

	CREATE_ITEM_AU_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_au BEFORE UPDATE ON item 
    BEGIN
        UPDATE item_fts SET 
            label = new.label,
            description = new.description
        WHERE item_id = old.id;
    END; `

	CREATE_ITEM_AD_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ad BEFORE DELETE ON item 
    BEGIN
        DELETE FROM item_fts WHERE item_id = old.id;
    END; `

	// Triggers for box
	CREATE_BOX_AI_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ai BEFORE INSERT ON box 
    BEGIN
        INSERT INTO box_fts(box_id, label, description) 
        VALUES (new.id, new.label, new.description);
    END; `

	CREATE_BOX_AU_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_au BEFORE UPDATE ON box 
    BEGIN
        UPDATE box_fts SET 
            label = new.label,
            description = new.description
        WHERE box_id = old.id;
    END; `

	CREATE_BOX_AD_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ad BEFORE DELETE ON box 
    BEGIN
        DELETE FROM box_fts WHERE box_id = old.id;
    END; `
)
