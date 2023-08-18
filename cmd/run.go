package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/MohamedGouaouri/envault/client"
	"github.com/MohamedGouaouri/envault/util"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Example: `
	envault run --upper --env=dev -- npm run dev
	envault run --command "first-command && second-command; more-commands..."
	`,
	Use:                   "run [any envault run command flags] -- [your application start command]",
	Short:                 "Used to inject environments variables into your application process",
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		// Check if the --command flag has been set
		commandFlagSet := cmd.Flags().Changed("command")

		// If the --command flag has been set, check if a value was provided
		if commandFlagSet {
			command := cmd.Flag("command").Value.String()
			if command == "" {
				return fmt.Errorf("you need to provide a command after the flag --command")
			}

			// If the --command flag has been set, args should not be provided
			if len(args) > 0 {
				return fmt.Errorf("you cannot set any arguments after --command flag. --command only takes a string command")
			}
		} else {
			// If the --command flag has not been set, at least one arg should be provided
			if len(args) == 0 {
				return fmt.Errorf("at least one argument is required after the run command, received %d", len(args))
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		envName, err := cmd.Flags().GetString("env")
		if err != nil {
			util.HandleError(err, "Unable to parse flag")
		}
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			util.HandleError(err, "Unable to parse flag")
		}

		secretsPath, err := cmd.Flags().GetString("path")
		if err != nil {
			util.HandleError(err, "Unable to parse flag")
		}

		toUpperCase, err := cmd.Flags().GetBool("upper")
		if err != nil {
			util.HandleError(err, "Unable to parse flag")
		}

		// Split paths
		secretsPaths := strings.Split(secretsPath, ",")
		var secrets = make(map[string]string)
		for _, path := range secretsPaths {
			s, err := client.GetAllEnvironmentVariables(token, envName, path)

			if err != nil {
				util.HandleError(err, "Could not fetch secrets", "If you are using a service token to fetch secrets, please ensure it is valid")
			}
			for k, v := range s {
				secrets[k] = v
			}
		}

		// secretsByKey := getSecretsByKeys(secrets)
		environmentVariables := make(map[string]string)

		// add all existing environment vars
		for _, s := range os.Environ() {
			kv := strings.SplitN(s, "=", 2)
			key := kv[0]
			value := kv[1]
			environmentVariables[key] = value
		}

		// now add vault secrets
		for k, v := range secrets {
			environmentVariables[k] = v
			if toUpperCase {
				environmentVariables[k] = util.MakeUpperCase(v)
			}
		}

		// turn it back into a list of envs
		var env []string
		for key, value := range environmentVariables {
			s := key + "=" + value
			env = append(env, s)
		}

		log.Debug().Msgf("injecting the following environment variables into shell: %v", env)

		if cmd.Flags().Changed("command") {
			command := cmd.Flag("command").Value.String()

			err = executeMultipleCommandWithEnvs(command, len(secrets), env)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			err = executeSingleCommandWithEnvs(args, len(secrets), env)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("token", "", "Fetch secrets using the Infisical Token")
	runCmd.Flags().StringP("env", "e", "secret", "Set the environment from which your secrets should be pulled from")
	runCmd.Flags().StringP("command", "c", "", "Chained commands to execute (e.g. \"npm install && npm run dev; echo ...\")")
	runCmd.Flags().String("path", "/", "Get secrets within a folder path")
	runCmd.Flags().Bool("upper", true, "Make secrets upper case before passing them to process")
}

// Will execute a single command and pass in the given secrets into the process
func executeSingleCommandWithEnvs(args []string, secretsCount int, env []string) error {
	command := args[0]
	argsForCommand := args[1:]
	color.Green("Injecting %v envvault secrets into your application process", secretsCount)

	cmd := exec.Command(command, argsForCommand...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	return execCmd(cmd)
}

func executeMultipleCommandWithEnvs(fullCommand string, secretsCount int, env []string) error {
	shell := [2]string{"sh", "-c"}
	if runtime.GOOS == "windows" {
		shell = [2]string{"cmd", "/C"}
	} else {
		currentShell := os.Getenv("SHELL")
		if currentShell != "" {
			shell[0] = currentShell
		}
	}

	cmd := exec.Command(shell[0], shell[1], fullCommand)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	color.Green("Injecting %v vault secrets into your application process", secretsCount)
	log.Debug().Msgf("executing command: %s %s %s \n", shell[0], shell[1], fullCommand)

	return execCmd(cmd)
}

// Credit: inspired by AWS Valut
func execCmd(cmd *exec.Cmd) error {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel)

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			sig := <-sigChannel
			_ = cmd.Process.Signal(sig) // process all sigs
		}
	}()

	if err := cmd.Wait(); err != nil {
		_ = cmd.Process.Signal(os.Kill)
		return fmt.Errorf("failed to wait for command termination: %v", err)
	}

	waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
	os.Exit(waitStatus.ExitStatus())
	return nil
}
