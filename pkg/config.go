package pkg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBHost      string           `json:"dbHost"`
	DBName      string           `json:"dbName"`
	DBUser      string           `json:"dbUser"`
	DBPassword  string           `json:"dbPassword"`
	DBPort      int              `json:"dbPort"`
	Action      Action           `json:"action"`
	ExportTable *ExportTable     `json:"exportTable"`
	ImportTable *ImportTableConf `json:"importTable"`
}

type ExportTable struct {
	ExportPath string   `json:"exportPath"`
	Tables     []string `json:"tables"`
}

type ImportTableConf struct {
	ImportPath   string         `json:"importPath"`
	ImportTables []*ImportTable `json:"tables"`
}

type Action string

const (
	ActionImport Action = "import"
	ActionExport Action = "export"
)

func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) ImportAction() bool {
	return c.Action == ActionImport
}

func (c *Config) Valid() error {
	if c.Action == "" {
		return fmt.Errorf("invalid action")
	}

	if c.DBHost == "" || c.DBName == "" || c.DBUser == "" || c.DBPort == 0 {
		return fmt.Errorf("invalid db config")
	}

	if c.Action == ActionImport && (c.ImportTable == nil || len(c.ImportTable.ImportTables) == 0) {
		return fmt.Errorf("import tables is empty")
	}

	if c.Action == ActionExport && (c.ExportTable == nil || len(c.ExportTable.Tables) == 0) {
		return fmt.Errorf("export tables is empty")
	}

	return nil
}
