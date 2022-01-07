/*
Copyright © 2021 Drew Stinnett <drew@drewlink.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/drewstinnett/vaultx/internal/unsealers"

	"github.com/spf13/cobra"
)

// unsealerCmd represents the unsealer command
var unsealerCmd = &cobra.Command{
	Use:   "unsealer",
	Short: "Unseal vault using unseal data from a standard place",
	Long:  `Use either a file, 1password secret, etc to perform an unesal`,
	Run: func(cmd *cobra.Command, args []string) {
		method, err := cmd.Flags().GetString("method")
		CheckErr(err, "")
		var unsealer unsealers.Unsealer
		switch method {
		case "op":
			unsealer = &unsealers.OPUnsealer{}
		default:
			log.Fatal("Unsupported unseal method")
		}
		err = unsealer.Prerequisites()
		CheckErr(err, "")
		unsealData, err := unsealer.FetchUnsealData(map[string]interface{}{
			"path": os.Getenv("OP_UNSEAL_PATH"),
		})
		CheckErr(err, "")
		err = unsealer.Unseal(*unsealData)
		CheckErr(err, "")
	},
}

func init() {
	rootCmd.AddCommand(unsealerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unsealerCmd.PersistentFlags().String("foo", "", "A help for foo")
	unsealerCmd.PersistentFlags().StringP("method", "m", "", "Method to use for unsealing")
	unsealerCmd.MarkPersistentFlagRequired("method")
	// viper.BindPFlag("method", unsealerCmd.PersistentFlags().Lookup("method"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unsealerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
