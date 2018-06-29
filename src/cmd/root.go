package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/G1itchZero/ZeroGo/site_manager"
	"github.com/G1itchZero/ZeroGo/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// workingDir is the path we're executing from
var workingDir string
var localNet bool

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

		zeronetAddress := "1FAiQ7MddvavaRF6b47fPEY4nSBVJUbCXf"
		if localNet {
			zeronetAddress = "133gv4M9kx5oWP1yUkK9MRFSnEVQZAghmt"
		}

		fmt.Println("Retrieving IPFS hash from ZeroNet address", zeronetAddress)

		ipfsHash, err := getIPFSHash(zeronetAddress)
		if err != nil {
			fmt.Printf("Unable to get IPFS hash from ZeroNet: %s\n", err)
			fmt.Print("Press enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		}
		// HACK: Pause here to wait for the zeronet lib to complete its output
		// TODO: Should actually be fixed in ZeroGo library itself
		time.Sleep(time.Second * 2)
		ipfsAddress := fmt.Sprintf("https://ipfs.io/ipfs/%s", string(ipfsHash))
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

		// Cleanup
		err = os.RemoveAll(utils.GetDataPath())
		if err != nil {
			fmt.Println(
				"Warning: unable to remove temporary ZeroNet data directory:",
				err)
		}

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

	rootCmd.PersistentFlags().BoolVar(&localNet, "testnet", false, "retrieve testnet seedlist")
}

// getIPFSHash retrieves the IPFS hash from ZeroNet
func getIPFSHash(zeronetAddress string) (string, error) {

	os.MkdirAll(utils.GetDataPath(), 0777)
	utils.CreateCerts()
	log.SetLevel(log.ErrorLevel)

	// If the site has been downloaded, remove and grab the latest version
	sitePath := path.Join(utils.GetDataPath(), zeronetAddress)
	exists, err := utils.Exists(sitePath)
	if err != nil {
		return "", err
	}
	if exists {
		err = os.RemoveAll(sitePath)
		if err != nil {
			return "", err
		}
	}

	siteManager := site_manager.NewSiteManager()
	site := siteManager.Get(zeronetAddress)
	site.Wait()
	// The IPFS hash is stored in the file ipfs.hash
	ipfsHashBytes, err := site.GetFile("ipfs.hash")
	if err != nil {
		return "", err
	}
	ipfsHash := string(ipfsHashBytes)
	ipfsHash = strings.TrimSpace(ipfsHash)

	return ipfsHash, nil
}
