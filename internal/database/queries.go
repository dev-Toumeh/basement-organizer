package database

const (
	BASIC_INFO_ID              = "id"
	BASIC_INFO_LABEL           = "label"
	BASIC_INFO_DESCRIPTION     = "description"
	BASIC_INFO_PICTURE         = "picture"
	BASIC_INFO_PREVIEW_PICTURE = "preview_picture"
	BASIC_INFO_QRCODE          = "qrcode"

	// single string with all columns of basic info which is present in every table
	ALL_BASIC_INFO_COLS string = "" +
		BASIC_INFO_ID + "," +
		BASIC_INFO_LABEL + "," +
		BASIC_INFO_DESCRIPTION + "," +
		BASIC_INFO_PICTURE + "," +
		BASIC_INFO_PREVIEW_PICTURE + "," +
		BASIC_INFO_QRCODE
	ALL_BASIC_INFO_COLS_LEN = 6

	FTS_ID              = BASIC_INFO_ID
	FTS_LABEL           = BASIC_INFO_LABEL
	FTS_DESCRIPTION     = BASIC_INFO_DESCRIPTION
	FTS_PREVIEW_PICTURE = BASIC_INFO_PREVIEW_PICTURE
	FTS_BOX_ID          = "box_id"
	FTS_BOX_LABEL       = "box_label"
	FTS_SHELF_ID        = "shelf_id"
	FTS_SHELF_LABEL     = "shelf_label"
	FTS_AREA_ID         = "area_id"
	FTS_AREA_LABEL      = "area_label"

	// single string with all columns of fts table
	ALL_FTS_COLS string = "" +
		FTS_ID + "," +
		FTS_LABEL + "," +
		FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + "," +
		FTS_BOX_ID + "," +
		FTS_BOX_LABEL + "," +
		FTS_SHELF_ID + "," +
		FTS_SHELF_LABEL + "," +
		FTS_AREA_ID + "," +
		FTS_AREA_LABEL

	// to use in create item, box, shelf, area statements
	CREATE_BASIC_INFO_BLOCK string = "" +
		BASIC_INFO_ID + " TEXT PRIMARY KEY," +
		BASIC_INFO_LABEL + " TEXT NOT NULL," +
		BASIC_INFO_DESCRIPTION + " TEXT," +
		BASIC_INFO_PICTURE + " TEXT," +
		BASIC_INFO_PREVIEW_PICTURE + " TEXT," +
		BASIC_INFO_QRCODE + " TEXT"

	// to use in create fts table (fts_item, fts_box, fts_shelf, fts_area) statements
	CREATE_FTS_BLOCK = "" +
		FTS_ID + " UNINDEXED," +
		FTS_LABEL + "," +
		FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + " UNINDEXED," +
		FTS_BOX_ID + " UNINDEXED," +
		FTS_BOX_LABEL + "," +
		FTS_SHELF_ID + " UNINDEXED," +
		FTS_SHELF_LABEL + "," +
		FTS_AREA_ID + " UNINDEXED," +
		FTS_AREA_LABEL

	// to use inside insert trigger statements for item and box tables
	CREATE_ITEM_BOX_INSERT_TRIGGER_VALUES_BLOCK = "" +
		"new." + BASIC_INFO_ID + "," +
		"new." + BASIC_INFO_LABEL + "," +
		"new." + BASIC_INFO_DESCRIPTION + "," +
		"new." + BASIC_INFO_PREVIEW_PICTURE + "," +
		"new." + FTS_BOX_ID + "," +
		"CASE " +
		"	WHEN new." + FTS_BOX_ID + " IS NOT NULL THEN (SELECT " + BASIC_INFO_LABEL + " FROM box WHERE id = new." + FTS_BOX_ID + ")" +
		"	ELSE NULL " +
		"END," +
		"new." + FTS_SHELF_ID + "," +
		"CASE " +
		"	WHEN new." + FTS_SHELF_ID + " IS NOT NULL THEN (SELECT " + BASIC_INFO_LABEL + " FROM shelf WHERE id = new." + FTS_SHELF_ID + ")" +
		"	ELSE NULL " +
		"END," +
		"new." + FTS_AREA_ID + "," +
		"CASE " +
		"	WHEN new." + FTS_AREA_ID + " IS NOT NULL THEN (SELECT " + BASIC_INFO_LABEL + " FROM area WHERE id = new." + FTS_AREA_ID + ")" +
		"	ELSE NULL " +
		"END"

	// to use inside update trigger statements
	UPDATE_TRIGGER_BLOCK string = "" +
		FTS_LABEL + " = new." + FTS_LABEL + "," +
		FTS_DESCRIPTION + " = new." + FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + " = new." + FTS_PREVIEW_PICTURE + "," +
		FTS_BOX_ID + " = new." + FTS_BOX_ID + "," +
		FTS_BOX_LABEL + " = (SELECT " + BASIC_INFO_LABEL + " FROM box WHERE box." + BASIC_INFO_ID + " = new." + ITEM_BOX_ID + ")," +
		FTS_SHELF_ID + " = new." + FTS_SHELF_ID + "," +
		FTS_SHELF_LABEL + " = (SELECT " + BASIC_INFO_LABEL + " FROM shelf WHERE shelf." + BASIC_INFO_ID + " = new." + ITEM_SHELF_ID + ")," +
		FTS_AREA_ID + " = new." + FTS_AREA_ID + "," +
		FTS_AREA_LABEL + " = (SELECT " + BASIC_INFO_LABEL + " FROM area WHERE area." + BASIC_INFO_ID + " = new." + ITEM_AREA_ID + ")"

	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT UNIQUE,
    passwordhash TEXT);`

	// Item
	// ITEM_TABLE_NAME = "item"

	ITEM_QUANTITY = "quantity"
	ITEM_WEIGHT   = "weight"
	ITEM_BOX_ID   = FTS_BOX_ID
	ITEM_SHELF_ID = FTS_SHELF_ID
	ITEM_AREA_ID  = FTS_AREA_ID
	ALL_ITEM_COLS = "" +
		ALL_BASIC_INFO_COLS + "," +
		ITEM_QUANTITY + "," +
		ITEM_WEIGHT + "," +
		ITEM_BOX_ID + "," +
		ITEM_SHELF_ID + "," +
		ITEM_AREA_ID

	CREATE_ITEM_TABLE_STMT string = "CREATE TABLE IF NOT EXISTS item (" +
		CREATE_BASIC_INFO_BLOCK + "," +
		ITEM_QUANTITY + " INTEGER," +
		ITEM_WEIGHT + " TEXT," +
		ITEM_BOX_ID + " TEXT REFERENCES box(id)," +
		ITEM_SHELF_ID + " TEXT REFERENCES shelf(id)," +
		ITEM_AREA_ID + " TEXT REFERENCES area(id)" +
		");"

	CREATE_ITEM_TABLE_STMT_FTS = "CREATE VIRTUAL TABLE IF NOT EXISTS item_fts USING fts5(" +
		CREATE_FTS_BLOCK + "," +
		"tokenize = 'porter'" +
		");"

	CREATE_ITEM_INSERT_TRIGGER = "" +
		"CREATE TRIGGER IF NOT EXISTS item_ai BEFORE INSERT ON item " +
		"BEGIN " +
		"    INSERT INTO item_fts(" + ALL_FTS_COLS + ")" +
		"    VALUES (" + CREATE_ITEM_BOX_INSERT_TRIGGER_VALUES_BLOCK + ");" +
		"END;"

	CREATE_ITEM_UPDATE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_au BEFORE UPDATE ON item 
    BEGIN
        UPDATE item_fts SET ` +
		UPDATE_TRIGGER_BLOCK +
		"WHERE " + BASIC_INFO_ID + "= new." + BASIC_INFO_ID + ";" +
		"END; "

	CREATE_ITEM_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS item_ad BEFORE DELETE ON item 
    BEGIN
        DELETE FROM item_fts WHERE id = old.id;
    END; `

	// Box
	ALL_BOX_COLS = ALL_BASIC_INFO_COLS + "," +
		FTS_BOX_ID + "," +
		FTS_SHELF_ID + "," +
		FTS_AREA_ID

	ALL_BOX_COLS_LEN = ALL_BASIC_INFO_COLS_LEN + 3

	CREATE_BOX_TABLE_STMT = "CREATE TABLE IF NOT EXISTS box (" +
		CREATE_BASIC_INFO_BLOCK + "," +
		FTS_BOX_ID + " TEXT REFERENCES box(" + BASIC_INFO_ID + ")," +
		FTS_SHELF_ID + " TEXT REFERENCES shelf(" + BASIC_INFO_ID + ")," +
		FTS_AREA_ID + " TEXT REFERENCES area(" + BASIC_INFO_ID + ")" +
		"); "

	CREATE_BOX_TABLE_STMT_FTS = "CREATE VIRTUAL TABLE IF NOT EXISTS box_fts USING fts5(" +
		CREATE_FTS_BLOCK + "," +
		"tokenize = 'porter'" +
		");"

	CREATE_BOX_INSERT_TRIGGER = "" +
		"CREATE TRIGGER IF NOT EXISTS box_ai BEFORE INSERT ON box " +
		"BEGIN " +
		"    INSERT INTO box_fts(" + ALL_FTS_COLS + ")" +
		"    VALUES (" + CREATE_ITEM_BOX_INSERT_TRIGGER_VALUES_BLOCK + ");" +
		"END; "

	CREATE_BOX_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS box_au BEFORE UPDATE ON box 
    BEGIN
        -- Update the original box's label
        UPDATE box_fts SET ` +
		UPDATE_TRIGGER_BLOCK +
		"WHERE " + BASIC_INFO_ID + "= new." + BASIC_INFO_ID + ";" +

		// "-- Update labels of boxes referencing this box as outerbox " +
		"UPDATE box_fts SET " +
		FTS_BOX_LABEL + " = new." + FTS_LABEL + " WHERE box_id = old.id;" +
		"END;"

	CREATE_BOX_DELETE_TRIGGER = `
    CREATE TRIGGER IF NOT EXISTS box_ad BEFORE DELETE ON box 
    BEGIN
        DELETE FROM box_fts WHERE id = old.id;
    END; `

	// Shelf
	SHELF_HEIGHT  = "height"
	SHELF_WIDTH   = "width"
	SHELF_DEPTH   = "depth"
	SHELF_ROWS    = "rows"
	SHELF_COLS    = "cols"
	SHELF_AREA_ID = "area_id"

	ALL_SHELF_COLS = ALL_BASIC_INFO_COLS + "," +
		SHELF_HEIGHT + "," +
		SHELF_WIDTH + "," +
		SHELF_DEPTH + "," +
		SHELF_ROWS + "," +
		SHELF_COLS + "," +
		SHELF_AREA_ID

	CREATE_SHELF_TABLE_STMT = "CREATE TABLE IF NOT EXISTS shelf (" + CREATE_BASIC_INFO_BLOCK + ", " +
		SHELF_HEIGHT + " REAL," +
		SHELF_WIDTH + " REAL," +
		SHELF_DEPTH + " REAL," +
		SHELF_ROWS + " INTEGER," +
		SHELF_COLS + " INTEGER," +
		SHELF_AREA_ID + " TEXT REFERENCES area(" + BASIC_INFO_ID + ")" +
		");"

	CREATE_SHELF_TABLE_STMT_FTS = "CREATE VIRTUAL TABLE IF NOT EXISTS shelf_fts USING fts5(" +
		CREATE_FTS_BLOCK + "," +
		"tokenize = 'unicode61'" +
		");"

	CREATE_SHELF_INSERT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS shelf_ai BEFORE INSERT ON shelf
	BEGIN
	    INSERT INTO shelf_fts(` +
		FTS_ID + "," +
		FTS_LABEL + "," +
		FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + "," +
		FTS_AREA_ID + "," +
		FTS_AREA_LABEL + ")" +
		`VALUES (` +
		"	new." + FTS_ID + "," +
		"	new." + FTS_LABEL + "," +
		"	new." + FTS_DESCRIPTION + "," +
		"	new." + FTS_PREVIEW_PICTURE + "," +
		"	new." + FTS_AREA_ID + "," +
		"	(SELECT " + BASIC_INFO_LABEL + " FROM area WHERE area." + BASIC_INFO_ID + " = new." + ITEM_AREA_ID + ")" +
		");" +
		"END;"

	CREATE_SHELF_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS shelf_au BEFORE UPDATE ON shelf
	BEGIN
		UPDATE shelf_fts SET ` +
		FTS_LABEL + " = new." + FTS_LABEL + "," +
		FTS_DESCRIPTION + " = new." + FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + " = new." + FTS_PREVIEW_PICTURE + "," +
		FTS_AREA_ID + " = new." + FTS_AREA_ID + ", " +
		FTS_AREA_LABEL + " = (SELECT " + BASIC_INFO_LABEL + " FROM area WHERE area." + BASIC_INFO_ID + " = new." + ITEM_AREA_ID + ")" +
		"WHERE " + BASIC_INFO_ID + "= new." + BASIC_INFO_ID + ";" +
		"END;"

	CREATE_SHELF_DELETE_TRIGGER = "CREATE TRIGGER IF NOT EXISTS shelf_ad BEFORE DELETE ON shelf " +
		"BEGIN " +
		"	DELETE FROM shelf_fts WHERE " + FTS_ID + " = old." + FTS_ID + ";" +
		"END;"

	// Area
	ALL_AREA_COLS = ALL_BASIC_INFO_COLS

	CREATE_AREA_TABLE_STMT = "CREATE TABLE IF NOT EXISTS area (" + ALL_AREA_COLS + ");"

	CREATE_AREA_TABLE_STMT_FTS = "CREATE VIRTUAL TABLE IF NOT EXISTS area_fts USING fts5(" +
		CREATE_FTS_BLOCK + "," +
		"tokenize = 'unicode61'" +
		");"

	CREATE_AREA_INSERT_TRIGGER = `CREATE TRIGGER IF NOT EXISTS area_ai BEFORE INSERT ON area
	BEGIN
	    INSERT INTO area_fts(` +
		FTS_ID + "," +
		FTS_LABEL + "," +
		FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + ")" +
		`VALUES (` +
		"	new." + FTS_ID + "," +
		"	new." + FTS_LABEL + "," +
		"	new." + FTS_DESCRIPTION + "," +
		"	new." + FTS_PREVIEW_PICTURE +
		");" +
		"END;"

	CREATE_AREA_UPDATE_TRIGGER = `CREATE TRIGGER IF NOT EXISTS area_au BEFORE UPDATE ON area
	BEGIN
		UPDATE area_fts SET ` +
		FTS_LABEL + " = new." + FTS_LABEL + "," +
		FTS_DESCRIPTION + " = new." + FTS_DESCRIPTION + "," +
		FTS_PREVIEW_PICTURE + " = new." + FTS_PREVIEW_PICTURE + " " +
		"WHERE " + BASIC_INFO_ID + "= new." + BASIC_INFO_ID + ";" +
		"END;"

	CREATE_AREA_DELETE_TRIGGER = "CREATE TRIGGER IF NOT EXISTS area_ad BEFORE DELETE ON area " +
		"BEGIN " +
		"	DELETE FROM area_fts WHERE " + FTS_ID + " = old." + FTS_ID + ";" +
		"END;"
)
