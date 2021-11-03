package cmd

import (
	"fmt"
	"github.com/jsiebens/inlets-on-fly/pkg/fly"
	"github.com/jsiebens/inlets-on-fly/pkg/names"
	"github.com/jsiebens/inlets-on-fly/pkg/templates"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Execute() error {
	rootCmd := rootCommand()
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(versionCommand())
	return rootCmd.Execute()
}

func rootCommand() *cobra.Command {
	return &cobra.Command{
		Use: "inlets-on-fly",
	}
}

func createCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "create",
		SilenceUsage: true,
	}

	var name string
	var org string
	var region string
	var ports []string

	command.Flags().StringVar(&name, "name", "", "")
	command.Flags().StringVar(&org, "org", "personal", "")
	command.Flags().StringVar(&region, "region", "", "")
	command.Flags().StringSliceVar(&ports, "ports", []string{"8080:80", "8080:443"}, "")

	command.RunE = func(command *cobra.Command, args []string) error {
		configuredPorts, err := parsePorts(ports)
		if err != nil {
			return err
		}

		if name == "" {
			name = names.GetRandomName()
		}

		token, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			return err
		}

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

		d := App{
			InletsVersion: "0.9.1",
			Name:          name,
			Token:         token,
			Ports:         configuredPorts,
		}

		if err := templates.Render(flyFile, "fly.toml.tpl", &d); err != nil {
			return err
		}

		if err := templates.Render(dockerFile, "Dockerfile.tpl", &d); err != nil {
			return err
		}

		if err := fly.CreateApp(dname, name, org); err != nil {
			return err
		}

		if region != "" {
			if err := fly.SetRegion(dname, region); err != nil {
				return err
			}
		}

		if err := fly.SetSecret(dname, token); err != nil {
			return err
		}

		if err := fly.Deploy(dname); err != nil {
			return err
		}

		if err := templates.Render(os.Stdout, "message.tpl", &d); err != nil {
			return err
		}

		return nil
	}

	return command
}

func parsePorts(values []string) ([]Ports, error) {
	var raw = map[int][]int{}

	for _, p := range values {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid value for port: %s", p)
		}

		internalPort, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		externalPort, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		externalPorts, ok := raw[internalPort]
		if !ok {
			raw[internalPort] = []int{externalPort}
		} else {
			raw[internalPort] = append(externalPorts, externalPort)
		}
	}

	var ports = []Ports{}
	for i, e := range raw {
		ports = append(ports, Ports{InternalPort: i, ExternalPorts: e})
	}
	return ports, nil
}

type App struct {
	InletsVersion string
	Name          string
	Token         string
	Ports         []Ports
}

type Ports struct {
	InternalPort  int
	ExternalPorts []int
}
