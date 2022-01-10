/*
Copyright © 2022 Drew Stinnett <drew@drewlink.com>

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
	"github.com/apex/log"
	"github.com/drewstinnett/vaultx/pkg/vaultx"
	"github.com/spf13/cobra"
)

// saveContextCmd represents the saveContext command
var saveContextCmd = &cobra.Command{
	Use:   "save NAME",
	Short: "Save context to state file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]
		if targetName == "environment" {
			log.Fatal("Invalid name. Cannot use 'environment' as a name.")
		}
		currentCtx, err := vaultx.GetCurrentContext(targetName)
		CheckErr(err, "Could not figure out the current context. Make sure that is working before attempting to save")
		log.Log.WithFields(log.Fields{
			"name": targetName,
		}).Info("Saving context to file")
		// Rename here
		currentCtx.Name = targetName
		err = vaultx.SaveContext("", currentCtx)
		CheckErr(err, "Could not save context")
	},
}

func init() {
	contextCmd.AddCommand(saveContextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveContextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveContextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
