package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// InstallFlags holds the flags for the install command
type InstallFlags struct {
	Override bool
}

var installCmd = &cobra.Command{
	Use:   "install [PRESET]",
	Short: "Enable Kool preset configuration in the current working directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   runInstall,
}

var installFlags = &InstallFlags{false}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&installFlags.Override, "override", "", false, "Force replace local existing files with the default preset files")
}

func runInstall(cmd *cobra.Command, args []string) {
	var (
		presetFiles                   map[string]string
		exists, hasExistingFile       bool
		preset, fileName, fileContent string
		err                           error
		file                          *os.File
		wrote                         int
	)

	preset = args[0]
	args = args[1:]

	if presetFiles, exists = presets[preset]; !exists {
		fmt.Println("Unknown preset", preset)
		os.Exit(2)
	}

	fmt.Println("Preset ", preset, " is gonna be installed!")

	for fileName = range presetFiles {
		if !installFlags.Override {
			if _, err = os.Stat(fileName); !os.IsNotExist(err) {
				fmt.Println("Preset file", fileName, "already exists.")
				hasExistingFile = true
			}
		}
	}

	if hasExistingFile {
		fmt.Println("Some preset files already exist. In case you wanna override them, use --override.")
		os.Exit(2)
	}

	for fileName, fileContent = range presetFiles {
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

		if err != nil {
			fmt.Println("Failed to create preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		if wrote, err = file.Write([]byte(fileContent)); err != nil {
			fmt.Println("Failed to write preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		if len([]byte(fileContent)) != wrote {
			fmt.Println("Failed to write preset file", fileName, " - failed to write all bytes:", wrote)
			os.Exit(2)
		}

		if err = file.Sync(); err != nil {
			fmt.Println("Failed to sync preset file", fileName, "due to error:", err)
			os.Exit(2)
		}

		file.Close()

		fmt.Println("  Preset file", fileName, " installed.")

		fileContent = ""
	}

	fmt.Println("Preset ", preset, " installed!")

	presetFiles = nil
}
