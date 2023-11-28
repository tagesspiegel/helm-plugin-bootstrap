package bootstrap

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

var (
	ErrChartYamlIsDirectory  = errors.New("chart.yaml is a directory")
	ErrValuesYamlIsDirectory = errors.New("values.yaml is a directory")
	ErrTemplatesIsDirectory  = errors.New("templates is a directory")
	ErrChartYamlNameMissing  = errors.New("chart metadata missing")
)

// parseChartYaml reads the Chart.yaml file from the given path and returns the parsed metadata.
func parseChartYaml(path string) (*chart.Metadata, error) {
	bts, err := os.ReadFile(filepath.Join(path, chartutil.ChartfileName))
	if err != nil {
		return nil, err
	}
	metadata := &chart.Metadata{}
	err = yaml.Unmarshal(bts, metadata)
	if err != nil {
		return nil, err
	}
	if metadata.Name == "" {
		return nil, ErrChartYamlNameMissing
	}
	return metadata, nil
}

// checkIsValidChartStructure checks if the given path contains a valid chart structure.
func checkIsValidChartStructure(chartDir string) error {
	// validate Chart.yaml
	fsInfo, err := os.Stat(filepath.Join(chartDir, chartutil.ChartfileName))
	if err != nil {
		// might not exist
		return err
	}
	if fsInfo.IsDir() {
		// is a directory not a file
		return ErrChartYamlIsDirectory
	}

	// validate templates folder
	fsInfo, err = os.Stat(path.Join(chartDir, chartutil.TemplatesDir))
	if err != nil {
		return err
	}
	if !fsInfo.IsDir() {
		return ErrTemplatesIsDirectory
	}

	// validate values.yaml
	fsInfo, err = os.Stat(path.Join(chartDir, chartutil.ValuesfileName))
	if err != nil {
		return err
	}
	if fsInfo.IsDir() {
		// is a directory not a file
		return ErrValuesYamlIsDirectory
	}

	return nil
}

// readValuesYaml reads the values.yaml file from the given path and returns the raw bytes.
// We do not parse the values.yaml file because we do not want to overwrite any existing comments.
func readValuesYaml(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(path, chartutil.ValuesfileName))
}

// Bootstrap reads properties from an existing Chart.yaml and creates additional configurations for `PodDisruptionBudgets`, `NetworkPolicies`.
//
// Inside of dir, this will read the Chart.yaml file, and create additional template files.
// It will also overwrite any existing files.
//
// If dir does not exist, this will return an error.
// If Chart.yaml or the templates directory do not exist, this will return an error.
func Bootstrap(chartLocation string, force bool) error {
	// check if Chart.yaml, values.yaml and templates folder exist
	err := checkIsValidChartStructure(chartLocation)
	if err != nil {
		return err
	}
	// parse Chart.yaml
	metadata, err := parseChartYaml(chartLocation)
	if err != nil {
		return err
	}
	templateFolderLocation := filepath.Join(chartLocation, chartutil.TemplatesDir)

	// read existing values.yaml
	valuesData, err := readValuesYaml(chartLocation)
	if err != nil {
		return err
	}

	files := []struct {
		path        string
		content     []byte
		propertyKey string
		values      []byte
	}{
		{
			// pdb.yaml
			path:        filepath.Join(templateFolderLocation, PodDisruptionBudgetFileName),
			content:     []byte(fmt.Sprintf(pdbTemplate, metadata.Name)),
			propertyKey: "pdb",
			values:      []byte(pdbValuesYaml),
		},
		{
			// networkpolicy.yaml
			path:        filepath.Join(templateFolderLocation, NetworkPolicyFileName),
			content:     []byte(fmt.Sprintf(networkPolicy, metadata.Name)),
			propertyKey: "networkPolicy",
			values:      []byte(networkPolicyValuesYaml),
		},
	}
	// write files
	for _, file := range files {
		if _, err := os.Stat(file.path); err == nil {
			// There is no handle to a preferred output stream here.
			fmt.Fprintf(os.Stderr, "WARNING: File %q already exists. Overwriting.\n", file.path)
		}
		if err := os.WriteFile(file.path, file.content, 0644); err != nil {
			return err
		}
		// write values.yaml if there is content
		if len(file.values) > 0 && (!bytes.Contains(valuesData, []byte(file.propertyKey)) || force) {
			valuesData = append(valuesData, []byte(file.values)...)
		}
	}
	// write values.yaml
	if err := os.WriteFile(filepath.Join(filepath.Join(chartLocation, chartutil.ValuesfileName)), valuesData, 0644); err != nil {
		return err
	}
	return nil
}
