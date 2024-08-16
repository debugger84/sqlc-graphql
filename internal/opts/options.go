package opts

import (
	"encoding/json"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Options struct {
	Package                     string            `json:"package" yaml:"package"`
	Out                         string            `json:"out" yaml:"out"`
	Overrides                   []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename                      map[string]string `json:"rename,omitempty" yaml:"rename"`
	OutputModelsFileName        string            `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
	OutputQuerierFileName       string            `json:"output_querier_file_name,omitempty" yaml:"output_querier_file_name"`
	OutputCopyfromFileName      string            `json:"output_copyfrom_file_name,omitempty" yaml:"output_copyfrom_file_name"`
	OutputFilesSuffix           string            `json:"output_files_suffix,omitempty" yaml:"output_files_suffix"`
	EmitExactTableNames         bool              `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	InflectionExcludeTableNames []string          `json:"inflection_exclude_table_names,omitempty" yaml:"inflection_exclude_table_names"`
	QueryParameterLimit         *int32            `json:"query_parameter_limit,omitempty" yaml:"query_parameter_limit"`
	OmitSqlcVersion             bool              `json:"omit_sqlc_version,omitempty" yaml:"omit_sqlc_version"`
	OmitUnusedStructs           bool              `json:"omit_unused_structs,omitempty" yaml:"omit_unused_structs"`
	DefaultSchema               string            `json:"default_schema,omitempty" yaml:"default_schema"`
	//GqlModelPackage             string            `json:"gql_model_package,omitempty" yaml:"gql_model_package"`
	//GqlOut                      string            `json:"gql_out,omitempty" yaml:"gql_out"`
	GenCommonParts bool     `json:"gen_common_parts,omitempty" yaml:"gen_common_parts"`
	Exclude        []string `json:"exclude,omitempty" yaml:"exclude"`
}

type GlobalOptions struct {
	Overrides []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename    map[string]string `json:"rename,omitempty" yaml:"rename"`
}

func Parse(req *plugin.GenerateRequest) (*Options, error) {
	options, err := parseOpts(req)
	if err != nil {
		return nil, err
	}
	global, err := parseGlobalOpts(req)
	if err != nil {
		return nil, err
	}
	if len(global.Overrides) > 0 {
		options.Overrides = append(global.Overrides, options.Overrides...)
	}
	if len(global.Rename) > 0 {
		if options.Rename == nil {
			options.Rename = map[string]string{}
		}
		maps.Copy(options.Rename, global.Rename)
	}
	return options, nil
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	if options.Package == "" {
		if options.Out != "" {
			options.Package = filepath.Base(options.Out)
		} else {
			return nil, fmt.Errorf("invalid options: missing package name")
		}
	}

	for i := range options.Overrides {
		if err := options.Overrides[i].parse(req); err != nil {
			return nil, err
		}
	}

	if options.QueryParameterLimit == nil {
		options.QueryParameterLimit = new(int32)
		*options.QueryParameterLimit = 1
	}

	return &options, nil
}

func parseGlobalOpts(req *plugin.GenerateRequest) (*GlobalOptions, error) {
	var options GlobalOptions
	if len(req.GlobalOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.GlobalOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling global options: %w", err)
	}
	for i := range options.Overrides {
		if err := options.Overrides[i].parse(req); err != nil {
			return nil, err
		}
	}
	return &options, nil
}

func ValidateOpts(opts *Options) error {
	if *opts.QueryParameterLimit < 0 {
		return fmt.Errorf("invalid options: query parameter limit must not be negative")
	}

	return nil
}
