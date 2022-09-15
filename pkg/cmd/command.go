package cmd

import (
	"fmt"
	cp "github.com/inlets/cloud-provision/provision"
	"github.com/jsiebens/inlets-on-fly/pkg/names"
	"github.com/jsiebens/inlets-on-fly/pkg/provision"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func Execute() error {
	rootCmd := rootCommand()
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(deleteCommand())
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

	var apiToken string
	var name string
	var org string
	var region string

	command.Flags().StringVar(&name, "name", "", "")
	command.Flags().StringVar(&org, "org", "", "")
	command.Flags().StringVar(&region, "region", "ams", "")
	command.Flags().StringVar(&apiToken, "api-token", "", "")

	command.RunE = func(command *cobra.Command, args []string) error {
		inletsVersion := "0.9.9"

		if apiToken == "" {
			apiToken = os.Getenv("FLY_API_TOKEN")
		}

		if apiToken == "" {
			fmt.Println("give a value --api-token or set the environment variable \"FLY_API_TOKEN\"")
			return nil
		}

		token, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			return err
		}

		hostReq := &cp.BasicHost{
			Region: region,
			Name:   names.GetRandomName(),
			Additional: map[string]string{
				"inlets-token":   token,
				"inlets-version": inletsVersion,
			},
		}

		provisioner, err := provision.NewFlyProvisioner(apiToken, org, region)
		if err != nil {
			return err
		}
		defer provisioner.Close()

		hostRes, err := provisioner.Provision(*hostReq)
		if err != nil {
			return err
		}

		fmt.Printf("Host: %s, status: %s\n", hostRes.ID, hostRes.Status)

		poll := time.Second * 2
		max := 500
		for i := 0; i < max; i++ {
			time.Sleep(poll)

			hostStatus, err := provisioner.Status(hostRes.ID)
			if err != nil {
				return err
			}

			fmt.Printf("[%d/%d] Host: %s, status: %s\n",
				i+1, max, hostStatus.ID, hostStatus.Status)

			if hostStatus.Status == "active" {

				fmt.Printf(`inlets Pro HTTPS (%s) server summary:
  IP: %s
  HTTPS Domains: %v
  Auth-token: %s

Command:

# Obtain a license at https://inlets.dev/pricing
# Store it at $HOME/.inlets/LICENSE or use --help for more options

# Where to route traffic from the inlets server
export UPSTREAM="http://127.0.0.1:8000"

inlets-pro http client --url "wss://%s:%d" \
--token "%s" \
--upstream $UPSTREAM

To delete:
  inlets-on-fly delete --id "%s"
`,
					inletsVersion,
					hostStatus.IP,
					fmt.Sprintf("%s.fly.dev", name),
					token,
					hostStatus.IP,
					8123,
					token,
					hostStatus.ID)

				return nil
			}
		}

		return nil
	}

	return command
}

func deleteCommand() *cobra.Command {
	command := &cobra.Command{
		Use: "delete",
	}

	var apiToken string
	var id string
	var org string
	var region string

	command.Flags().StringVar(&id, "id", "", "")
	command.Flags().StringVar(&org, "org", "", "")
	command.Flags().StringVar(&region, "region", "ams", "")
	command.Flags().StringVar(&apiToken, "api-token", "", "")

	command.RunE = func(command *cobra.Command, args []string) error {
		if apiToken == "" {
			apiToken = os.Getenv("FLY_API_TOKEN")
		}

		if apiToken == "" {
			fmt.Println("give a value --api-token or set the environment variable \"FLY_API_TOKEN\"")
			return nil
		}

		provisioner, err := provision.NewFlyProvisioner(apiToken, org, region)
		if err != nil {
			return err
		}
		defer provisioner.Close()

		request := cp.HostDeleteRequest{ID: id}

		return provisioner.Delete(request)
	}

	return command
}
