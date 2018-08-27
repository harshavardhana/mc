/*
 * Minio Client (C) 2017, 2018 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"io"
	"os"

	"github.com/minio/cli"
	"github.com/minio/mc/pkg/probe"
)

var adminProfilingFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "start",
		Usage: "Start specific profiling type, possible types are: 'cpu', 'mem', 'block'",
	},
	cli.BoolFlag{
		Name:  "stop",
		Usage: "Stop profiling",
	},
	cli.StringFlag{
		Name:  "download",
		Usage: "Download latest profiling data under the specified filename",
	},
}

var adminProfilingCmd = cli.Command{
	Name:            "profiling",
	Usage:           "Generate profiling data for debugging purposes",
	Action:          mainAdminProfiling,
	Before:          setGlobalsFromContext,
	Flags:           append(adminProfilingFlags, globalFlags...),
	HideHelpCommand: true,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS] TARGET

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}
EXAMPLES:
    1. Start CPU profiling
       $ {{.HelpName}} --start cpu myminio/

    2. Stop profiling
       $ {{.HelpName}} --stop myminio/

    3. Download latest profiling data under the specified filename
       $ {{.HelpName}} --download /tmp/profile.data myminio/

`,
}

func checkAdminProfilingSyntax(ctx *cli.Context) {
	isStart := ctx.IsSet("start")
	isStop := ctx.IsSet("stop")
	isDownload := ctx.IsSet("download")

	switch {
	case len(ctx.Args()) != 1:
		fallthrough
	case isStart && (isStop || isDownload):
		fallthrough
	case isStop && (isStart || isDownload):
		fallthrough
	case isDownload && (isStart || isStop):
		cli.ShowCommandHelpAndExit(ctx, "profiling", 1) // last argument is exit code
	}

	if isStart {
		startArg := ctx.String("start")
		switch startArg {
		case "cpu", "mem", "block":
		default:
			cli.ShowCommandHelpAndExit(ctx, "profiling", 1) // last argument is exit code
		}
	}
}

// mainAdminProfiling - the entry function of profiling command
func mainAdminProfiling(ctx *cli.Context) error {

	// Check for command syntax
	checkAdminProfilingSyntax(ctx)

	// Get the alias parameter from cli
	args := ctx.Args()
	aliasedURL := args.Get(0)

	// Create a new Minio Admin Client
	client, err := newAdminClient(aliasedURL)
	if err != nil {
		fatalIf(err.Trace(aliasedURL), "Cannot initialize admin client.")
		return nil
	}

	startArg := ctx.String("start")
	stopFlag := ctx.Bool("stop")
	downloadArg := ctx.String("download")

	switch {
	case startArg != "":
		cmdErr := client.StartProfiling(startArg)
		fatalIf(probe.NewError(cmdErr), "Unable to start profiling:")
	case stopFlag:
		cmdErr := client.StopProfiling()
		fatalIf(probe.NewError(cmdErr), "Unable to stop profiling:")
	case downloadArg != "":
		rd, cmdErr := client.DownloadProfilingData()
		fatalIf(probe.NewError(cmdErr), "Unable to download profiling data:")
		file, e := os.Create(downloadArg)
		fatalIf(probe.NewError(e), "Unable to download profiling data:")
		io.Copy(file, rd)
		rd.Close()
		file.Close()
	}

	return nil
}
