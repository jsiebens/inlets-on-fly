package cmd

import (
	"github.com/jsiebens/inlets-on-fly/pkg/names"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
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
		Use: "create",
	}

	command.AddCommand(createTcpServerCommand())
	command.AddCommand(createHttpServerCommand())

	return command
}

func createHttpServerCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "http",
		SilenceUsage: true,
	}

	var name string
	var org string
	var region string

	command.Flags().StringVar(&name, "name", "", "")
	command.Flags().StringVar(&org, "org", "personal", "")
	command.Flags().StringVar(&region, "region", "", "")

	command.RunE = func(command *cobra.Command, args []string) error {
		configuredPorts := []Ports{{InternalPort: 8000, ExternalPorts: []int{80, 443}}}

		if name == "" {
			name = names.GetRandomName()
		}

		token, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			return err
		}

		app := App{
			Mode:          "http",
			InletsVersion: "0.9.1",
			Name:          name,
			Org:           org,
			Region:        region,
			Token:         token,
			Ports:         configuredPorts,
		}

		return createApp(&app)
	}

	return command
}

func createTcpServerCommand() *cobra.Command {
	command := &cobra.Command{
		Use:          "tcp",
		SilenceUsage: true,
	}

	var name string
	var org string
	var region string
	var ports []string

	command.Flags().StringVar(&name, "name", "", "")
	command.Flags().StringVar(&org, "org", "personal", "")
	command.Flags().StringVar(&region, "region", "", "")
	command.Flags().StringSliceVar(&ports, "ports", []string{"8000:80", "8000:443"}, "")

	command.RunE = func(command *cobra.Command, args []string) error {
		configuredPorts, err := parsePorts(ports)
		if err != nil {
			return nil
		}

		if name == "" {
			name = names.GetRandomName()
		}

		token, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			return err
		}

		app := App{
			Mode:          "tcp",
			InletsVersion: "0.9.1",
			Name:          name,
			Org:           org,
			Region:        region,
			Token:         token,
			Ports:         configuredPorts,
		}

		return createApp(&app)
	}

	return command
}
