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
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

func fatalf(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use: "decrypt <dir> [<dir_1> [<dir_2> [...]]]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fatalf(cmd.UsageString())
		}

		key := parseKey(cmd)

		for _, path := range args {
			if err := walk(path, func(path string, plaintext []byte) error {
				if path[len(path)-len(".secure"):] != ".secure" {
					return nil
				}

				ciphertext, err := cryptopasta.Decrypt(plaintext, &key)
				if err != nil {
					return err
				}

				logrus.Errorf("Decrypted file \"%s\"", path)

				return ioutil.WriteFile(path[:len(path)-len(".secure")], ciphertext, 0774)
			}); err != nil {
				fatalf("An error occurred: %s", err)
			}
		}

		fmt.Println("Done")
	},
}

func parseKey(cmd *cobra.Command) [32]byte {
	key, err := cmd.Flags().GetString("key")
	if err != nil {
		fatalf("You must supply a key")
	} else if len(key) < 32 {
		fatalf("Expected key length to be 32, got %d", len(key))
	}

	var keyArray [32]byte
	copy(keyArray[:], []byte(key)[:32])

	return keyArray
}

func walk(dir string, cb func(path string, in []byte) error) error {
	return filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		input, err := ioutil.ReadFile(path)
		if err != nil {
			logrus.Errorf("Unable to read file \"%s\"", path)
			return err
		}

		return cb(path, input)
	})
}

func init() {
	RootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("key", "k", "", "Key to be used")
}
