package cmd

import (
	"fmt"
	"os"
	"os/user"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zellyn/kooky"
)

var (
	cfgFile string
	domain, name string
	short bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cookie-getter",
	Short: "Extract cookies from web browsers",
	Long: `Reads, decrypts and prints cookies from the web browser's cookie storage database.
The default output format is: <domain>/<name>: <value>
With the --short option, the tool outputs only the cookie values, one per line
`,
		Run: func(cmd *cobra.Command, args []string) {
			var cookies []*kooky.Cookie

			usr, _ := user.Current()

			cookiesFile := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Cookies", usr.HomeDir)
			cookies, err := kooky.ReadChromeCookies(cookiesFile, domain, name, time.Time{})

			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			}

			if len(cookies) == 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "No value found\n")
				os.Exit(1)
			}

			for _, cookie := range cookies {
				switch short {
				case true:
					fmt.Printf("%s\n", cookies[0].Value)
				case false:
					fmt.Printf("%s/%s: %s\n", cookie.Domain, cookie.Name, cookie.Value)
				}
			}
		},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cookie-getter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain for which to get the cookie(s)")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the cookie to get")
	rootCmd.Flags().BoolVarP(&short, "short", "s", false, "Only output cookie value(s)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cookie-getter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cookie-getter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
