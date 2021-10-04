// Copyright Project Contour Authors
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

package cli

import (
	"fmt"
	"os"

	"github.com/robinfoe/kuberparse/pkg/version"
	"github.com/spf13/cobra"
)

func mustString(s string, err error) string {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", version.Progname, err)
		os.Exit(int(EX_CONFIG))
	}

	return s
}

func anyString(values ...string) bool {
	for _, s := range values {
		if s != "" {
			return true
		}
	}

	return false
}

// Defaults applies default settings to a cobra.Command to improve the help output.
func Defaults(c *cobra.Command) *cobra.Command {
	c.SilenceUsage = true
	c.SilenceErrors = true
	c.DisableFlagsInUseLine = true

	return c
}
