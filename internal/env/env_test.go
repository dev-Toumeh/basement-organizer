package env

import (
	"basement/main/internal/logg"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    option
		expectedErr bool
	}{
		"valid option with value": {
			input:    "option=true",
			expected: option{"option", "true"},
		},
		"valid with trailing comment": {
			input:    "option3=true # this is valid option",
			expected: option{"option3", "true"},
		},
		"empty line": {
			input:       "",
			expected:    option{},
			expectedErr: true,
		},
		"new line": {
			input:       "         \n",
			expected:    option{},
			expectedErr: true,
		},
		"multiple options in one line": {
			input:       "option1=val option2=val2",
			expected:    option{},
			expectedErr: true,
		},
		"missing '='": {
			input:       "option",
			expected:    option{},
			expectedErr: true,
		},
		"missing value": {
			input:       "option=",
			expected:    option{},
			expectedErr: true,
		},
		"missing option": {
			input:       "=true",
			expected:    option{},
			expectedErr: true,
		},
		"ignoring comment": {
			input:       "# =true",
			expected:    option{},
			expectedErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opt, err := parseLine(tt.input)
			if tt.expected != opt {
				t.Errorf("got:\"%s\" expected \"%s\"", opt, tt.expected)
			}
			if tt.expectedErr && (err == nil) {
				t.Errorf("got error: nil expected error with input \"%s\"", tt.input)
			}
			if !tt.expectedErr && (err != nil) {
				t.Errorf("got error: \"%s\" expected no error with input \"%s\"", logg.CleanLastError(err), tt.input)
			}
		})
	}

}

func TestParseWithReader(t *testing.T) {
	// config := defaultDevConfigPreset
	config := &Configuration{}
	config.Init()

	ll := `
		env=2
	debugLogsEnabled=true
	dbPath=asdfasdfkj
	# comment ignored
	useMemoryDB=true # comment with correctly used option
	defaultTableSize=0
	debugLogsEnabled
		=
	`
	r := strings.NewReader(ll)
	_, _ = parseWithReader(r, config)
	// if len(errs) != 0 {
	// 	for _, e := range errs {
	// 		t.Error(logg.CleanLastError(e))
	// 	}
	// }
}

var validConfig Configuration = Configuration{env: env_test, defaultTableSize: 1, dbPath: "/"}

func TestCheckDefaultTableSizeConstraints(t *testing.T) {
	// config := &Configuration{}

	tableConstraintsTests := map[string]struct {
		input Configuration
		// expected    *int
		expectedErr bool
	}{
		"invalid defaultTableSize 0": {
			input:       Configuration{dbPath: "/", defaultTableSize: 0},
			expectedErr: true,
		},
		"invalid defaultTableSize -1": {
			input:       Configuration{dbPath: "/", defaultTableSize: -1},
			expectedErr: true,
		},
		"valid defaultTableSize": {
			input: Configuration{dbPath: "/", defaultTableSize: 1},
		},
	}

	for name, tt := range tableConstraintsTests {
		t.Run(name, func(t *testing.T) {
			err := validateDefaultTableSize(&tt.input)
			if tt.expectedErr && (err == nil) {
				t.Errorf("got error: nil expected error with input \"%d\"", tt.input.defaultTableSize)
			}
			if !tt.expectedErr && (err != nil) {
				t.Errorf("got error: \"%s\" expected no error with input \"%d\"", logg.CleanLastError(err), tt.input.defaultTableSize)
			}
		})
	}
}

func TestCheckDBConstraints(t *testing.T) {
	dbConstraintsTests := map[string]struct {
		input       Configuration
		expectedErr bool
	}{
		"invalid db": {
			input:       Configuration{defaultTableSize: 1, dbPath: "dbpath", useMemoryDB: true},
			expectedErr: true,
		},
		"invalid db2": {
			input:       Configuration{defaultTableSize: 1, dbPath: ":memory:", useMemoryDB: false},
			expectedErr: true,
		},
		"empty db path": {
			input:       Configuration{defaultTableSize: 1, dbPath: "", useMemoryDB: false},
			expectedErr: true,
		},
		"valid db": {
			input: Configuration{defaultTableSize: 1, dbPath: ":memory:", useMemoryDB: true},
		},
		"valid db2": {
			input: Configuration{defaultTableSize: 1, dbPath: "dbpath", useMemoryDB: false},
		},
	}

	for name, tt := range dbConstraintsTests {
		t.Run(name, func(t *testing.T) {
			err := validateDBOptions(&tt.input)
			if tt.expectedErr && (err == nil) {
				t.Errorf("got error: nil expected error with input dbPath=\"%s\" useMemoryDB=\"%t\"", tt.input.dbPath, tt.input.useMemoryDB)
			}
			if !tt.expectedErr && (err != nil) {
				t.Errorf("got error: \"%s\" expected no error with input dbPath=\"%s\" useMemoryDB=\"%t\"", logg.CleanLastError(err), tt.input.dbPath, tt.input.useMemoryDB)
			}
		})
	}
}

func TestApplyLastLine(t *testing.T) {
	newConfigFromFile := `
	dbPath=:memory:
	useMemoryDB=true`

	config := Configuration{
		dbPath:      "testdb.db",
		useMemoryDB: false}
	r := strings.NewReader(newConfigFromFile)
	out, _ := parseWithReader(r, &config)
	useMemoryDB, ok := out["useMemoryDB"]
	if !ok {
		t.Error("dropped useMemoryDB field")
	}
	if useMemoryDB != "true" {
		t.Errorf("got useMemoryDB: \"%v\" expected \"false\"", useMemoryDB)
	}
}

func TestApply(t *testing.T) {
	oldEnv := env_test
	oldDefaultTableSize := 1
	oldDBPath := "/"
	oldInfoLogsEnabled := false
	oldDebugLogsEnabled := false

	v1 := Configuration{
		env:              oldEnv,
		defaultTableSize: oldDefaultTableSize,
		dbPath:           oldDBPath,
		infoLogsEnabled:  oldInfoLogsEnabled,
		debugLogsEnabled: oldDebugLogsEnabled,
	}

	newConfigFromFile := `debugLogsEnabled=true
defaultTableSize=2
	dbPath=/new/path # comment: change db path
env=1 # attempt to override env should not work
infoLogsEnabled==true # invalid


	`
	newDebugLogsEnabled := true
	newDefaultTableSize := 2
	newDbPath := "/new/path"

	config := &v1
	config.env = oldEnv
	config.infoLogsEnabled = oldInfoLogsEnabled
	config.Init()

	r := strings.NewReader(newConfigFromFile)
	out, _ := parseWithReader(r, config)
	// if len(errs) != 0 {
	// 	for _, e := range errs {
	// 		t.Error(logg.CleanLastError(e))
	// 	}
	// }

	applyParsedOptions(out, config)

	// Check bool conversion
	if config.debugLogsEnabled != newDebugLogsEnabled {
		t.Errorf("got config: \"%v\" expected \"%v\"", config.debugLogsEnabled, newDebugLogsEnabled)
	}
	// Check Int conversion
	if config.defaultTableSize != newDefaultTableSize {
		t.Errorf("got config: \"%v\" expected \"%v\"", config.defaultTableSize, newDefaultTableSize)
	}
	// Check String conversion
	if config.dbPath != newDbPath {
		t.Errorf("got config: \"%v\" expected \"%v\"", config.dbPath, newDbPath)
	}
	// Don't override env
	if config.env != oldEnv {
		t.Errorf("got config: \"%v\" expected \"%v\"", config.env, oldEnv)
	}
	// Don't override infoLogsEnabled
	if config.infoLogsEnabled != oldInfoLogsEnabled {
		t.Errorf("got config: \"%v\" expected \"%v\"", config.infoLogsEnabled, oldInfoLogsEnabled)
	}
}
