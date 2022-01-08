package cmd

import (
	"fmt"
	"github.com/jsiebens/inlets-on-fly/pkg/fly"
	"github.com/jsiebens/inlets-on-fly/pkg/templates"
	"os"
	"path/filepath"
)

func parsePorts(values []uint) ([]Ports, error) {
	var ports []Ports
	for _, p := range values {
		ports = append(ports, Ports{InternalPort: p, ExternalPorts: []uint{p}})
	}
	return ports, nil
}

func createApp(app *App) error {
	dname, err := os.MkdirTemp("", "inletsfly-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dname)

	fmt.Println("Temp dir name:", dname)

	flyFile, err := os.OpenFile(filepath.Join(dname, "fly.toml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	dockerFile, err := os.OpenFile(filepath.Join(dname, "Dockerfile"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if err := templates.Render(flyFile, "fly.toml.tpl", &app); err != nil {
		return err
	}

	if err := templates.Render(dockerFile, "Dockerfile.tpl", &app); err != nil {
		return err
	}

	if err := fly.CreateApp(dname, app.Name, app.Org); err != nil {
		return err
	}

	if app.Region != "" {
		if err := fly.SetRegion(dname, app.Region); err != nil {
			return err
		}
	}

	if err := fly.SetSecret(dname, app.Token); err != nil {
		return err
	}

	if err := fly.Deploy(dname); err != nil {
		return err
	}

	if err := templates.Render(os.Stdout, "message.tpl", &app); err != nil {
		return err
	}

	return nil
}
