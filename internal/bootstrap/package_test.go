package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

var (
	chartName     = "test-chart"
	chartLocation = path.Join(os.TempDir(), "charts")
)

func TestMain(m *testing.M) {
	binPath, err := exec.LookPath("helm")
	if err != nil {
		fmt.Println("helm not found")
		os.Exit(1)
	}
	err = os.MkdirAll(chartLocation, 0755)
	if err != nil {
		fmt.Println("failed to create chart directory")
		os.Exit(1)
	}
	defer os.RemoveAll(chartLocation)

	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cf()

	err = os.Chdir(chartLocation)
	if err != nil {
		fmt.Println("failed to create chart directory")
		os.Exit(1)
	}

	if err := exec.CommandContext(ctx, binPath, "create", chartName).Run(); err != nil {
		fmt.Println("failed to create chart", err.Error())
		os.Exit(1)
	}

	if exitCode := m.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
