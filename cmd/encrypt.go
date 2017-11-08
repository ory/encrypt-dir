// Copyright © 2017 Aeneas Rekkas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/gtank/cryptopasta"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

// decryptCmd represents the decrypt command
var encryptCmd = &cobra.Command{
	Use: "encrypt <dir> [<dir_1> [<dir_2> [...]]]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fatalf(cmd.UsageString())
		}

		key := parseKey(cmd)

		for _, path := range args {
			if err := walk(path, func(path string, plaintext []byte) error {
				if path[len(path)-len(".secure"):] == ".secure" {
					return nil
				}

				ciphertext, err := cryptopasta.Encrypt(plaintext, &key)
				if err != nil {
					return err
				}

				logrus.Errorf("Encrypted file \"%s\"", path)

				return ioutil.WriteFile(path + ".secure", ciphertext, 0774)
			}); err != nil {
				fatalf("An error occurred: %s", err)
			}
		}

		fmt.Println("Done")
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("key", "k", "", "Key to be used")
}
