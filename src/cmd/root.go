package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// workingDir is the path we're executing from
var workingDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retrieve-zn-ipfs-seedlist",
	Short: "Stellite ZeroNet/IPFS seedlist downloader",
	Long: `
  __ _____ ___ _   _   _ _____ ___
/' _/_   _| __| | | | | |_   _| __|
'._'. | | | _|| |_| |_| | | | | _|
|___/ |_| |___|___|___|_| |_| |___|
            ZERONET/IPFS SEEDLIST DOWNLOADER

retrieve-zn-ipfs-seedlist downloads the latest available seedlist
from ZeroNet/IPFS for the Stellite daemon to use instead of the hardcoded ones.
`,

	// By default the root command executes a download and import of the
	// latest blockchain file
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf(`
  __ _____ ___ _   _   _ _____ ___
/' _/_   _| __| | | | | |_   _| __|
'._'. | | | _|| |_| |_| | | | | _|
|___/ |_| |___|___|___|_| |_| |___|
            ZERONET/IPFS SEEDLIST DOWNLOADER
			`)
		// Clear all the spaces
		fmt.Println("")

		ipfsAddress := "https://ipfs.io/ipfs/QmctxmwEvU7oXwydATLLM5iScm1vXRtvwep9KY6KzxtT6V"
		fmt.Printf("Downloading seedlist from IPFS at '%s'\n", ipfsAddress)

		response, err := http.Get(ipfsAddress)
		if err != nil {
			fmt.Printf("Unable to get seedlist from IPFS: %s\n", err)
			fmt.Print("Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		}

		seedlistPath := filepath.Join(workingDir, "ipfs-seedlist.txt")
		fileBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Unable to read IPFS response: %s\n", err)
			fmt.Print("Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		}

		err = ioutil.WriteFile(seedlistPath, fileBytes, 0644)
		if err != nil {
			fmt.Printf("Unable to save seedlist: %s\n", err)
			fmt.Print("Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		}

		fmt.Printf("\nSeedlist saved to '%s'\n", seedlistPath)
		fmt.Println("You can start the Stellite daemon now.")

		fmt.Print("\nPress enter to continue...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""

	// Declare err to not shadow workingDir
	var err error
	workingDir, err = os.Executable()
	if err != nil {
		fmt.Println(err)
		fmt.Print("Press enter to continue...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	}
	workingDir = filepath.Dir(workingDir)

}
