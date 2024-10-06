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
    preview_picture TEXT,
    quantity INTEGER,
    weight TEXT,
    qrcode TEXT,
    box_id TEXT REFERENCES box(id),
    shelf_id TEXT REFERENCES shelf(id),
    area_id TEXT REFERENCES area(id)
	);`

	CREATE_BOX_TABLE_STMT = `CREATE TABLE IF NOT EXISTS box (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL, 
    description TEXT,
    picture TEXT,
    preview_picture TEXT,
    qrcode TEXT,
    outerbox_id TEXT REFERENCES box(id));`

	CREATE_SHELF_TABLE_STMT = `CREATE TABLE IF NOT EXISTS shelf (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL, 
    description TEXT,
    picture TEXT,
    preview_picture TEXT,
    qrcode TEXT,
    height REAL,
    width REAL,
    depth REAL,
    rows INTEGER,
    cols INTEGER
	);`

	CREATE_AREA_TABLE_STMT = `CREATE TABLE IF NOT EXISTS area (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    picture TEXT 
    preview_picture TEXT 
	);`

	CREATE_ITEM_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS item_fts USING fts5(
    item_id UNINDEXED,
    label, 
    description, 
    preview_picture UNINDEXED,
    tokenize = 'porter'
	);`

	CREATE_BOX_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS box_fts USING fts5(
    box_id UNINDEXED,
    label, 
    outerbox_label,
    outerbox_id UNINDEXED,
	preview_picture UNINDEXED,
    tokenize = 'porter'
	); `

	// Triggers for item
	CREATE_ITEM_INSERT_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ai BEFORE INSERT ON item 
    BEGIN
        INSERT INTO item_fts(item_id, label, description, preview_picture) 
        VALUES (new.id, new.label, new.description, new.preview_picture);
    END;`

	CREATE_ITEM_UPDATE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_au BEFORE UPDATE ON item 
    BEGIN
        UPDATE item_fts SET 
            label = new.label,
            description = new.description,
            preview_picture = new.preview_picture
        WHERE item_id = old.id;
    END; `

	CREATE_ITEM_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ad BEFORE DELETE ON item 
    BEGIN
        DELETE FROM item_fts WHERE item_id = old.id;
    END; `

	// Triggers for box
	CREATE_BOX_INSERT_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ai BEFORE INSERT ON box
    BEGIN
        INSERT INTO box_fts(box_id, label, outerbox_label, outerbox_id, preview_picture) 
        VALUES (
            new.id, 
            new.label,
            CASE 
                WHEN new.outerbox_id IS NOT NULL THEN (SELECT label FROM box WHERE id = new.outerbox_id)
                ELSE NULL 
            END,
            new.outerbox_id,
            new.preview_picture
        );
    END; `

	CREATE_BOX_UPDATE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_au BEFORE UPDATE ON box 
    BEGIN
        -- Update the original box's label
        UPDATE box_fts 
        SET 
			label = new.label, 
			preview_picture = new.preview_picture
        WHERE box_id = old.id;

        -- Update labels of boxes referencing this box as outerbox
        UPDATE box_fts
        SET outerbox_label = new.label
        WHERE outerbox_id = old.id;
    END; `

	CREATE_BOX_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ad BEFORE DELETE ON box 
    BEGIN
        DELETE FROM box_fts WHERE box_id = old.id;
    END; `
)
