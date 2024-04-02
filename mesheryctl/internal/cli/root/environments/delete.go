// Copyright 2024 Layer5, Inc.
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

package environments

import (
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteEnvironmentCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a new environments",
	Long:  `delete a new environments by providing the name and description of the environment`,
	Example: `
// delete a new environment
mesheryctl exp environment delete environmentId
// Documentation for environment can be found at:
https://docs.layer5.io/cloud/spaces/environments/
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		return nil
	},

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp environment delete \nRun 'mesheryctl exp environment delete --help' to see detailed help message"
		if len(args) != 1 {
			return errors.New(utils.EnvironmentSubError(fmt.Sprintf("accepts 1 arg(s), received %d\n%s", len(args), errMsg), "delete"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/environments/%s", baseUrl, args[0])
		req, err := utils.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		// Check if the response status code is 200
		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("Connection deleted successfully")
			return nil
		}

		return utils.ErrBadRequest(errors.New(fmt.Sprintf("failed to delete environment with id %s", args[0])))
	},
}
