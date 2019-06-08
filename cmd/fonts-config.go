package main

import (
	"fmt"
	"github.com/marguerite/util/fileutils"
	"github.com/openSUSE/fonts-config/lib"
	"github.com/urfave/cli"
	"os"
	"os/user"
	"path/filepath"
)

// Version fonts-config's version
const Version string = "20190608"

func removeUserSetting(prefix string, verbosity int) error {
	for _, f := range []string{
		filepath.Join(prefix, "fontconfig/fonts-config"),
		filepath.Join(prefix, "fontconfig/rendering-options.conf"),
		filepath.Join(prefix, "fontconfig/family-prefer.conf"),
	} {
		err := fileutils.Remove(f, lib.VerbosityDebug)
		if err != nil {
			return err
		}
	}
	return nil
}

func yastInfo() {
	// compatibility only, no actual use.
	fmt.Printf("Involved Files\n" +
		"  rendering config: /etc/fonts/conf.d/10-rendering-options.conf\n" +
		"  java fontconfig properties: /usr/lib*/jvm/jre/lib/fontconfig.SUSE.properties\n" +
		"  user sysconfig file: fontconfig/fonts-config\n" +
		"  metric compatibility avail: /usr/share/fontconfig/conf.avail/30-metric-aliases.conf\n" +
		"  metric compatibility bw symlink: /etc/fonts/conf.d/31-metric-aliases-bw.conf\n" +
		"  metric compatibility config: /etc/fonts/conf.d/30-metric-aliases.conf\n" +
		"  local family list: /etc/fonts/conf.d/58-family-prefer-local.conf\n" +
		"  metric compatibility symlink: /etc/fonts/conf.d/30-metric-aliases.conf\n" +
		"  user family list: fontconfig/family-prefer.conf\n" +
		"  java fontconfig properties template: /usr/share/fonts-config/fontconfig.SUSE.properties.template\n" +
		"  rendering config template: /usr/share/fonts-config/10-rendering-options.conf.template\n" +
		"  sysconfig file: /etc/sysconfig/fonts-config\n" +
		"  user rendering config: fontconfig/rendering-options.conf\n")
}

func main() {
	var userMode, remove, force, ttcap, enableJava, quiet, verbose, debug, autohint, bw, bwMono, ebitmaps, info bool
	var hintstyle, lcdfilter, rgba, ebitmapsLang, emojis string

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Display version and exit.",
	}
	app := cli.NewApp()
	app.Usage = "openSUSE fontconfig presets generator."
	app.Description = "openSUSE fontconfig presets generator."
	app.Version = Version
	app.Authors = []cli.Author{
		{Name: "Marguerite Su", Email: "marguerite@opensuse.org"},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "user, u",
			Usage:       "Run fontconfig setup for user.",
			Destination: &userMode,
		},
		cli.BoolFlag{
			Name:        "remove-user-setting, r",
			Usage:       "Remove current user's fontconfig setup.",
			Destination: &remove,
		},
		cli.BoolFlag{
			Name:        "force, f",
			Usage:       "Force the update of all generated files even if it appears unnecessary according to the time stamps",
			Destination: &force,
		},
		cli.BoolTFlag{
			Name:        "quiet, q",
			Usage:       "Work silently, unless an error occurs.",
			Destination: &quiet,
		},
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "Print some progress messages to standard output.Print some progress messages to standard output.",
			Destination: &verbose,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Print a lot of debugging messages to standard output.",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "force-hintstyle",
			Value:       "hintslight",
			Usage:       "Which `hintstyle` to use: hintfull, hintmedium, hintslight or hintnone.",
			Destination: &hintstyle,
		},
		cli.BoolFlag{
			Name:        "autohint",
			Usage:       "Use autohint even for well hinted fonts.",
			Destination: &autohint,
		},
		cli.BoolFlag{
			Name:        "force-bw",
			Usage:       "Do not use antialias.",
			Destination: &bw,
		},
		cli.BoolFlag{
			Name:        "force-bw-monospace",
			Usage:       "Do not use antialias for well instructed monospace fonts.",
			Destination: &bwMono,
		},
		cli.StringFlag{
			Name:        "lcdfilter",
			Value:       "lcddefault",
			Usage:       "Which `lcdfilter` to use: lcdnone, lcddefault, lcdlight, lcdlegacy.",
			Destination: &lcdfilter,
		},
		cli.StringFlag{
			Name:        "rgba",
			Value:       "rgb",
			Usage:       "Which `subpixel arrangement` your monitor use: none, rgb, vrgb, bgr, vbgr, unknown.",
			Destination: &rgba,
		},
		cli.BoolTFlag{
			Name:        "ebitmaps",
			Usage:       "Whether to use embedded bitmaps or not",
			Destination: &ebitmaps,
		},
		cli.StringFlag{
			Name:        "ebitmapslang",
			Value:       "ja:ko:zh-CN:zh-SG:zh-TW:zh-HK:zh-MO",
			Usage:       "Argument contains a `list` of colon separated languages, for example \"ja:ko:zh-CN\" which means \"use embedded bitmaps only for fonts supporting Japanese, Korean, or Simplified Chinese.",
			Destination: &ebitmapsLang,
		},
		cli.StringFlag{
			Name:        "emojis",
			Value:       "Noto Color Emoji",
			Usage:       "Default emoji fonts. for example\"Noto Color Emoji:Twemoji Mozilla\", glyphs from these fonts will be blacklisted for other non-emoji fonts",
			Destination: &emojis,
		},
		cli.BoolTFlag{
			Name:        "ttcap",
			Usage:       "Generate TTCap entries..",
			Destination: &ttcap,
		},
		cli.BoolTFlag{
			Name:        "java",
			Usage:       "Generate font setup for Java.",
			Destination: &enableJava,
		},
		cli.BoolFlag{
			Name:        "info",
			Usage:       "Print files used by fonts-config for YaST Fonts module.",
			Destination: &info,
		},
	}

	app.Action = func(c *cli.Context) error {

		if info {
			yastInfo()
			os.Exit(0)
		}

		// parse verbosity
		verbosity := 0
		verbosityMap := map[bool]int{quiet: lib.VerbosityQuiet, verbose: lib.VerbosityVerbose, debug: lib.VerbosityDebug}
		for k, v := range verbosityMap {
			if k {
				verbosity = v
				break
			}
		}

		currentUser, _ := user.Current()
		if !userMode && currentUser.Uid != "0" && currentUser.Username != "root" {
			fmt.Println("*** error: no root permissions; rerun with --user for user fontconfig setting.")
			os.Exit(1)
		}

		userPrefix := filepath.Join(lib.GetEnv("HOME"), ".config")

		options := lib.Options{verbosity, hintstyle, autohint, bw, bwMono, lcdfilter, rgba, ebitmaps, ebitmapsLang, "Noto Color Emoji", "", "", "", true, false, ttcap, enableJava}

		if remove {
			err := removeUserSetting(userPrefix, verbosity)
			lib.ErrChk(err)
			os.Exit(0)
		}

		config := lib.LoadOptions("/etc/sysconfig/fonts-config", lib.Options{})

		if userMode {
			config = lib.LoadOptions(filepath.Join(userPrefix, "fontconfig/fonts-config"), config)
		}

		config.Merge(options)
		config.Write(userMode)

		if verbosity >= lib.VerbosityDebug {
			if userMode {
				fmt.Printf("USER mode (%s)\n", lib.GetEnv("USER"))
			} else {
				fmt.Println("SYSTEM mode")
			}

			fmt.Printf("--- sysconfig options (read from /etc/sysconfig/fonts-config")
			if userMode {
				fmt.Printf(", %s)\n", filepath.Join(userPrefix, "fontconfig/fonts-config"))
			} else {
				fmt.Printf(")\n")
			}

			fmt.Printf(config.Bounce())
			fmt.Println("---")
		}

		if !userMode {
			err := lib.MkFontScaleDir(config, force)
			lib.ErrChk(err)
		}

		/*	# The following two calls may change files in /etc/fonts, therefore
			# they have to be called *before* fc-cache. If anything is
			# changed in /etc/fonts after calling fc-cache, fontconfig
			# will think that the cache files are out of date again. */

		err := lib.GenerateDefaultRenderingOptions(userMode, config)
		lib.ErrChk(err)

		err = lib.GenerateFamilyPreferenceLists(userMode, config)
		lib.ErrChk(err)

		err = lib.GenerateEmojiBlacklist(userMode, config)
		lib.ErrChk(err)

		if !userMode {
			lib.FcCache(config.Verbosity)
			lib.FpRehash(config.Verbosity)
			if config.GenerateJavaFontSetup {
				lib.GenerateJavaFontSetup(config.Verbosity)
			}
			lib.ReloadXfsConfig(config.Verbosity)
		}

		return nil
	}

	_ = app.Run(os.Args)
}
