package fly

import (
	"fmt"
	execute "github.com/alexellis/go-execute/pkg/v1"
)

func CreateApp(cwd string, name string, org string) error {
	task := execute.ExecTask{
		Cwd:         cwd,
		Command:     "flyctl",
		Args:        []string{"apps", "create", "--name", name, "--org", org},
		StreamStdio: true,
	}
	return check(task.Execute())
}

func SetRegion(cwd string, region string) error {
	task := execute.ExecTask{
		Cwd:         cwd,
		Command:     "flyctl",
		Args:        []string{"regions", "set", region},
		StreamStdio: true,
	}
	return check(task.Execute())
}

func SetSecret(cwd string, token string) error {
	task := execute.ExecTask{
		Cwd:         cwd,
		Command:     "flyctl",
		Args:        []string{"secrets", "set", fmt.Sprintf("TOKEN=%s", token)},
		StreamStdio: true,
	}
	return check(task.Execute())
}

func Deploy(cwd string) error {
	task := execute.ExecTask{
		Cwd:         cwd,
		Command:     "flyctl",
		Args:        []string{"deploy"},
		StreamStdio: true,
	}
	return check(task.Execute())
}

func check(result execute.ExecResult, err error) error {
	if err != nil {
		return err
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("invalid exit code '%v'", result.ExitCode)
	}

	return nil
}
