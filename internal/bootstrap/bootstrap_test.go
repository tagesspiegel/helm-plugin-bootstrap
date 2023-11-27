package bootstrap

import (
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"helm.sh/helm/v3/pkg/chart"
)

func Test_parseChartYaml(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *chart.Metadata
		wantErr bool
	}{
		{
			name: "with valid chart.yaml",
			args: args{
				path: path.Join(chartLocation, chartName),
			},
			want: &chart.Metadata{
				Name:        chartName,
				Version:     "0.1.0",
				AppVersion:  "1.16.0",
				Description: "A Helm chart for Kubernetes",
				APIVersion:  "v2",
				Type:        "application",
			},
		},
		{
			name: "file not found",
			args: args{
				path: path.Join(chartLocation, "not-found"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseChartYaml(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseChartYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("parseChartYaml() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func Test_checkIsValidChartStructure(t *testing.T) {
	type args struct {
		chartDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with valid chart structure",
			args: args{
				chartDir: path.Join(chartLocation, chartName),
			},
			wantErr: false,
		},
		{
			name: "with invalid chart structure",
			args: args{
				chartDir: path.Join(chartLocation, "not-found"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkIsValidChartStructure(tt.args.chartDir); (err != nil) != tt.wantErr {
				t.Errorf("checkIsValidChartStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readValuesYaml(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name: "with valid values.yaml",
			args: args{
				path: path.Join(chartLocation, chartName),
			},
			wantValue: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readValuesYaml(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readValuesYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantValue && got == nil {
				t.Errorf("readValuesYaml() got = %v, want %v", got, tt.wantValue)
				return
			}
		})
	}
}

func TestBootstrap(t *testing.T) {
	type args struct {
		chartLocation string
		force         bool
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "with valid chart structure",
			args: args{
				chartLocation: path.Join(chartLocation, chartName),
				force:         false,
			},
			wantErr: false,
			wantFiles: []string{
				"Chart.yaml",
				"values.yaml",
				".helmignore",
				"templates/_helpers.tpl",
				"templates/deployment.yaml",
				"templates/hpa.yaml",
				"templates/ingress.yaml",
				"templates/NOTES.txt",
				"templates/service.yaml",
				"templates/serviceaccount.yaml",
				"templates/pdb.yaml",
				"templates/networkpolicy.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Bootstrap(tt.args.chartLocation, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, file := range tt.wantFiles {
				// check if file exists
				finfo, err := os.Stat(path.Join(tt.args.chartLocation, file))
				if err != nil {
					t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// check if file is not empty
				if finfo.Size() == 0 {
					t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
