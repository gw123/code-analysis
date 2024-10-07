package code

// Define the structure for YAML parsing
type FileInfo struct {
	FileName    string   `yaml:"file_name"`
	PackageName string   `yaml:"package_name"`
	Imports     []string `yaml:"imports"`
}

type Constant struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

type Method struct {
	Name         string   `yaml:"name"`
	Params       []string `yaml:"params"`
	ReturnValues []string `yaml:"return_values"`
	Description  string   `yaml:"description"`
}

type Struct struct {
	Name    string   `yaml:"name"`
	Fields  []string `yaml:"fields"`
	Methods []Method `yaml:"methods"`
}

type ParsedYAML struct {
	FunctionDescription string   `yaml:"file_description"`
	FileInfo            FileInfo `yaml:"file_info"`
	//Constants           []Constant `yaml:"constants"`
	//Structs             []Struct   `yaml:"structs"`
	//Methods             []Method   `yaml:"methods"`
}
