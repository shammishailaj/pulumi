// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/pulumi/pulumi/pkg/backend"
	"github.com/pulumi/pulumi/pkg/engine"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
)

func newDestroyCmd() *cobra.Command {
	var debug bool
	var stack string

	var message string

	// Flags for engine.UpdateOptions.
	var analyzers []string
	var color colorFlag
	var diffDisplay bool
	var parallel int
	var preview string
	var showConfig bool
	var showReplacementSteps bool
	var showSames bool
	var nonInteractive bool

	var cmd = &cobra.Command{
		Use:        "destroy",
		SuggestFor: []string{"delete", "down", "kill", "remove", "rm", "stop"},
		Short:      "Destroy an existing stack and its resources",
		Long: "Destroy an existing stack and its resources\n" +
			"\n" +
			"This command deletes an entire existing stack by name.  The current state is\n" +
			"loaded from the associated snapshot file in the workspace.  After running to completion,\n" +
			"all of this stack's resources and associated state will be gone.\n" +
			"\n" +
			"Warning: although old snapshots can be used to recreate a stack, this command\n" +
			"is generally irreversible and should be used with great care.",
		Args: cmdutil.NoArgs,
		Run: cmdutil.RunFunc(func(cmd *cobra.Command, args []string) error {
			interactive := isInteractive(cmd)
			behavior, err := previewFlagsToBehavior(interactive, preview)
			if err != nil {
				return err
			}

			s, err := requireStack(stack, false)
			if err != nil {
				return err
			}
			proj, root, err := readProject()
			if err != nil {
				return err
			}

			m, err := getUpdateMetadata(message, root)
			if err != nil {
				return errors.Wrap(err, "gathering environment metadata")
			}

			if behavior.Interactive() {
				prompt := fmt.Sprintf("This will permanently destroy all resources in the '%s' stack!", s.Name())
				if !confirmPrompt(prompt, s.Name().String()) {
					return errors.New("confirmation declined")
				}
			}

			err = s.Destroy(
				proj, root, m,
				engine.UpdateOptions{
					Analyzers: analyzers,
					Parallel:  parallel,
					Debug:     debug,
				},
				behavior,
				backend.DisplayOptions{
					Color:                color.Colorization(),
					ShowConfig:           showConfig,
					ShowReplacementSteps: showReplacementSteps,
					ShowSameResources:    showSames,
					IsInteractive:        interactive,
					DiffDisplay:          diffDisplay,
					Debug:                debug,
				},
				cancellationScopes,
			)
			if err == context.Canceled {
				return errors.New("destroy cancelled")
			}
			return err
		}),
	}

	cmd.PersistentFlags().BoolVarP(
		&debug, "debug", "d", false,
		"Print detailed debugging output during resource operations")
	cmd.PersistentFlags().StringVarP(
		&stack, "stack", "s", "",
		"Choose a stack other than the currently selected one")
	cmd.PersistentFlags().StringVarP(
		&message, "message", "m", "",
		"Optional message to associate with the destroy operation")

	// Flags for engine.UpdateOptions.
	cmd.PersistentFlags().StringSliceVar(
		&analyzers, "analyzer", []string{},
		"Run one or more analyzers as part of this update")
	cmd.PersistentFlags().VarP(
		&color, "color", "c", "Colorize output. Choices are: always, never, raw, auto")
	cmd.PersistentFlags().BoolVar(
		&diffDisplay, "diff", false,
		"Display operation as a rich diff showing the overall change")
	cmd.PersistentFlags().IntVarP(
		&parallel, "parallel", "p", 0,
		"Allow P resource operations to run in parallel at once (<=1 for no parallelism)")
	cmd.PersistentFlags().StringVar(
		&preview, "preview", "",
		"Preview behavior. Choices are: only (dry-run), skip (no preview, just update), auto (auto-accept)")
	cmd.PersistentFlags().BoolVar(
		&showConfig, "show-config", false,
		"Show configuration keys and variables")
	cmd.PersistentFlags().BoolVar(
		&showReplacementSteps, "show-replacement-steps", false,
		"Show detailed resource replacement creates and deletes instead of a single step")
	cmd.PersistentFlags().BoolVar(
		&showSames, "show-sames", false,
		"Show resources that don't need to be updated because they haven't changed, alongside those that do")
	cmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "Disable interactive mode")

	return cmd
}
