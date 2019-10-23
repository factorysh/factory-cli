// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package main

import (
	_ "github.com/factorysh/factory-cli/cmd/container" // container commands
	_ "github.com/factorysh/factory-cli/cmd/infos"     // infos commands
	_ "github.com/factorysh/factory-cli/cmd/journal"   // journal commands
	"github.com/factorysh/factory-cli/cmd/root"
	_ "github.com/factorysh/factory-cli/cmd/runjob"  // runjob command
	_ "github.com/factorysh/factory-cli/cmd/upgrade" // upgrade commands
	_ "github.com/factorysh/factory-cli/cmd/volume"  // volume commands
)

func main() {
	root.Execute()
}
