package database

const (
	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT UNIQUE,
    passwordhash TEXT);`

	// Item
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

	CREATE_ITEM_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS item_fts USING fts5(
    id UNINDEXED,
    label, 
    description, 
    preview_picture UNINDEXED,
    box_id UNINDEXED,
    box_label,
	shelf_id UNINDEXED,
	shelf_label,
	area_id UNINDEXED,
	area_label,
    tokenize = 'porter'
	);`

	CREATE_ITEM_INSERT_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ai BEFORE INSERT ON item 
    BEGIN
        INSERT INTO item_fts(id, label, description, preview_picture, box_id, box_label, shelf_id, shelf_label, area_id, area_label) 
        VALUES (
			new.id,
			new.label,
			new.description,
			new.preview_picture,
			new.box_id,
			(SELECT label FROM box WHERE box.id = new.box_id),
			new.shelf_id,
			(SELECT label FROM shelf WHERE shelf.id = new.shelf_id),
			new.area_id,
			(SELECT label FROM area WHERE area.id = new.area_id)
		);
    END;`

	CREATE_ITEM_UPDATE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_au BEFORE UPDATE ON item 
    BEGIN
        UPDATE item_fts SET 
            label = new.label,
            description = new.description,
            preview_picture = new.preview_picture,
			box_id = new.box_id, 
			box_label = (SELECT label FROM box WHERE box.id = new.box_id),
			shelf_id = new.shelf_id, 
			shelf_label = (SELECT label FROM shelf WHERE shelf.id = new.shelf_id),
			area_id = new.area_id, 
			area_label = (SELECT label FROM area WHERE area.id = new.area_id)
        WHERE id = new.id;
    END; `

	CREATE_ITEM_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ad BEFORE DELETE ON item 
    BEGIN
        DELETE FROM item_fts WHERE id = old.id;
    END; `

	// Box
	CREATE_BOX_TABLE_STMT = `CREATE TABLE IF NOT EXISTS box (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL, 
    description TEXT,
    picture TEXT,
    preview_picture TEXT,
    qrcode TEXT,
    box_id TEXT REFERENCES box(id),
    shelf_id TEXT REFERENCES shelf(id),
    area_id TEXT REFERENCES area(id)
	); `

	CREATE_BOX_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS box_fts USING fts5(
    id UNINDEXED,
    label, 
    box_id UNINDEXED,
    box_label,
	shelf_id UNINDEXED,
	shelf_label,
	area_id UNINDEXED,
	area_label,
	preview_picture UNINDEXED,
    tokenize = 'porter'
	); `

	CREATE_BOX_INSERT_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ai BEFORE INSERT ON box
    BEGIN
        INSERT INTO box_fts(id, label, box_label, box_id, shelf_id, shelf_label, area_id, area_label, preview_picture) 
        VALUES (
            new.id, 
            new.label,
            CASE 
                WHEN new.box_id IS NOT NULL THEN (SELECT label FROM box WHERE id = new.box_id)
                ELSE NULL 
            END,
            new.box_id,
            new.shelf_id,
            CASE 
                WHEN new.shelf_id IS NOT NULL THEN (SELECT label FROM shelf WHERE id = new.shelf_id)
                ELSE NULL 
            END,
            new.area_id,
            CASE 
                WHEN new.area_id IS NOT NULL THEN (SELECT label FROM area WHERE id = new.area_id)
                ELSE NULL 
            END,
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
			box_id = new.box_id,
			box_label = (SELECT label FROM box WHERE box.id = new.box_id),
			shelf_id = new.shelf_id, 
			shelf_label = (SELECT label FROM shelf WHERE shelf.id = new.shelf_id),
			area_id = new.area_id, 
			area_label = (SELECT label FROM area WHERE area.id = new.area_id),
			preview_picture = new.preview_picture
        WHERE id = old.id;

        -- Update labels of boxes referencing this box as outerbox
        UPDATE box_fts
        SET box_label = new.label
        WHERE box_id = old.id;

    END; `

	CREATE_BOX_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ad BEFORE DELETE ON box 
    BEGIN
        DELETE FROM box_fts WHERE id = old.id;
    END; `

	// Shelf
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
    cols INTEGER,
	area_id TEXT REFERENCES area(id)
	);`

	CREATE_SHELF_TABLE_STMT_FTS = `CREATE VIRTUAL TABLE IF NOT EXISTS shelf_fts USING fts5(
    id UNINDEXED,
    label,
    area_id UNINDEXED,
    area_label,
    preview_picture UNINDEXED,
    tokenize = 'unicode61'
	);`

	CREATE_SHELF_INSERT_TRIGGER = `
	CREATE TRIGGER IF NOT EXISTS shelf_ai BEFORE INSERT ON shelf
	BEGIN
		INSERT INTO shelf_fts(id, label, area_label, area_id, preview_picture)
		VALUES (
			new.id,
			new.label,
			(SELECT label FROM area WHERE id = new.area_id),
			new.area_id,
			new.preview_picture
		);
	END;`

	CREATE_SHELF_UPDATE_TRIGGER = `
	CREATE TRIGGER IF NOT EXISTS shelf_au BEFORE UPDATE ON shelf
	BEGIN
		UPDATE shelf_fts SET
			label = new.label,
			area_label = (SELECT label FROM area WHERE id = new.area_id),
			area_id = new.area_id,
			preview_picture = new.preview_picture
		WHERE id = old.id;
	END;`

	CREATE_SHELF_DELETE_TRIGGER = `
	CREATE TRIGGER IF NOT EXISTS shelf_ad BEFORE DELETE ON shelf
	BEGIN
		DELETE FROM shelf_fts WHERE id = old.id;
	END;`

	// Area
	CREATE_AREA_TABLE_STMT = `CREATE TABLE IF NOT EXISTS area (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    picture TEXT,
    preview_picture TEXT 
	);`
)
