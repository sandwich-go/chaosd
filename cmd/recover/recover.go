// Copyright 2020 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package recover

import (
	"context"
	"fmt"
	"github.com/chaos-mesh/chaosd/pkg/core"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"

	"github.com/chaos-mesh/chaosd/cmd/server"
	"github.com/chaos-mesh/chaosd/pkg/server/chaosd"
	"github.com/chaos-mesh/chaosd/pkg/utils"
)

type recoverCommand struct {
	uid string
	all bool
}

func NewRecoverCommand() *cobra.Command {
	options := &recoverCommand{}
	dep := fx.Options(
		server.Module,
		fx.Provide(func() *recoverCommand {
			return options
		}),
	)

	cmd := &cobra.Command{
		Use:               "recover UID",
		Short:             "Recover a chaos experiment",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: completeUid,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				options.uid = args[0]
			}
			utils.FxNewAppWithoutLog(dep, fx.Invoke(recoverCommandF)).Run()
		},
	}
	cmd.Flags().BoolVarP(&options.all, "all", "A", false, "recover all chaos attacks")
	return cmd
}

func recoverCommandF(chaos *chaosd.Server, options *recoverCommand) {
	if options.all {
		exps, err := chaos.Search(&core.SearchCommand{All: true})
		if err != nil {
			utils.ExitWithError(utils.ExitError, err)
		}
		for _, v := range exps {
			err = chaos.RecoverAttack(v.Uid)
			if err != nil {
				utils.ExitWithError(utils.ExitError, err)
			}
			fmt.Printf("Recover uid: %s successfully\n", v.Uid)
		}
	} else {
		err := chaos.RecoverAttack(options.uid)
		if err != nil {
			utils.ExitWithError(utils.ExitError, err)
		}
		fmt.Printf("Recover uid: %s successfully\n", options.uid)
	}

	utils.NormalExit("")
}

func completeUid(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	completionCtx := newCompletionCtx()
	completionDep := fx.Options(
		server.Module,
		fx.Provide(func() *completionContext {
			return completionCtx
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := utils.FxNewAppWithoutLog(completionDep, fx.Invoke(listUid)).Start(ctx); err != nil {
			log.Error(errors.Wrap(err, "start application").Error())
		}
	}()
	var uids []string
	for {
		select {
		case uid := <-completionCtx.uids:
			if len(uid) == 0 {
				return uids, cobra.ShellCompDirectiveNoFileComp
			}
			uids = append(uids, uid)
		case err := <-completionCtx.err:
			log.Error(err.Error())
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
}
