package lib

import (
	"fmt"
	"log"
	"strings"
)

func genBitmapLanguagesConfig(opts Options) string {
	tmp := "\t<match target=\"font\">\n"
	tmp += "\t\t<edit name=\"embeddedbitmap\" mode=\"append\">\n"
	if opts.UseEmbeddedBitmaps {
		if opts.EmbeddedBitmapsLanguages == "no" || len(opts.EmbeddedBitmapsLanguages) == 0 {
			tmp += "\t\t\t<b>true</bool>\n"
			tmp += "\t\t</edit>\n"
			tmp += "\t</match>\n"
			return tmp
		}
		tmp += "\t\t\t<bool>false</bool>\n"
		tmp += "\t\t</edit>\n"
		tmp += "\t</match>\n"
		for _, v := range strings.Split(opts.EmbeddedBitmapsLanguages, ":") {
			tmp += "\t<match target=\"font\">\n"
			tmp += "\t\t<test name=\"lang\" compare=\"contains\"><string>" + v + "</string></test>\n"
			tmp += "\t\t<edit name=\"embeddedbitmap\" mode=\"append\"><bool>true</bool></edit>\n"
			tmp += "\t</match>\n"
		}
		return tmp
	}
	tmp += "\t\t\t<bool>false</bool>\n"
	tmp += "\t\t</edit>\n"
	tmp += "\t</match>\n"
	return tmp
}

// GenRenderingOptions generates fontconfig rendering options conf
func GenRenderingOptions(userMode bool, opts Options) {
	/* # reflect fonts-config syconfig variables or
	   # parameters in fontconfig setting to control rendering */
	renderFile := GetConfigLocation("render", userMode)

	Dbg(opts.Verbosity, Debug, fmt.Sprintf("Generating %s.", renderFile))
	renderText := genRenderingOptions(opts, userMode)

	err := overwriteOrRemoveFile(renderFile, []byte(renderText), 0644)
	if err != nil {
		log.Fatalf("Can not write %s: %s\n", renderFile, err.Error())
	}
}

func genRenderingOptions(opts Options, userMode bool) string {
	config := ""
	config += genStringOptionConfig(opts.Verbosity, opts.ForceHintstyle, "Forcing hintstyle:",
		"<!-- Choose preferred common hinting style here.  -->\n<!-- Possible values: no, hitnone, hitslight, hintmedium and hintfull. -->\n<!-- Can be overridden with some other options, e. g. force_bw\n\tor force_bw_monospace => hintfull -->\n",
		"force_hintstyle", false, true)
	config += genBoolOptionConfig(opts.Verbosity, opts.ForceAutohint, "Forcing autohint:",
		"<!-- Force autohint always. -->\n<!-- If false, for well hinted fonts, their instructions are used for rendering. -->\n",
		"force_autohint", true)
	config += genBoolOptionConfig(opts.Verbosity, opts.ForceBw, "Forcing black and white:",
		"<!-- Do not use font smoothing (black&white rendering) at all.  -->\n",
		"force_bw", true)
	config += genBoolOptionConfig(opts.Verbosity, opts.ForceBwMonospace, "Forcing black and white for good hinted monospace:",
		"<!-- Do not use font smoothing for some monospaced fonts.  -->\n<!-- Liberation Mono, Courier New, Andale Mono, Monaco, etc. -->\n",
		"force_bw_monospace", true)
	config += genStringOptionConfig(opts.Verbosity, opts.UseLcdfilter, "Lcdfilter:",
		"<!-- Set LCD filter. Amend when you want use subpixel rendering. -->\n<!-- Don't forgot to set correct subpixel ordering in 'rgba' element. -->\n<!-- Possible values: lcddefault, lcdlight, lcdlegacy, lcdnone -->\n",
		"lcdfilter", true, false)
	config += genStringOptionConfig(opts.Verbosity, opts.UseRgba, "Subpixel arrangement:",
		"<!-- Set LCD subpixel arrangement and orientation.  -->\n<!-- Possible values: unknown, none, rgb, bgr, vrgb, vbgr. -->\n",
		"rgba", true, false)
	config += genBitmapLanguagesConfig(opts)
	config += genBoolOptionConfig(opts.Verbosity, opts.SearchMetricCompatible, "Search metric compatible fonts:",
		"<!-- Search for metric compatible families? -->\n",
		"search_metric_aliases", false)
	config += genUserInclude(userMode)
	if len(config) == 0 {
		return config
	}
	return genFcPreamble(userMode, "<!-- using target=\"pattern\", because we want to change pattern in 60-family-prefer.conf\n\tregarding to this setting -->\n") +
		config + FcSuffix
}

// validStringOption return false if a string is "null", has suffix "none" or just empty.
func validStringOption(opt string) bool {
	if len(opt) == 0 || opt == "null" || strings.HasSuffix(opt, "none") {
		return false
	}
	return true
}

func genStringOptionConfig(verbosity int, opt, dbgOutput, comment, editName string, cst, force bool) string {
	if !validStringOption(opt) {
		return ""
	}
	Dbg(verbosity, Debug, fmt.Sprintf(dbgOutput+" %s", opt))
	config := comment
	config += "\t<match target=\"pattern\" >\n\t\t<edit name=\""
	config += editName
	config += "\" mode=\""
	if force {
		config += "assign"
	} else {
		config += "append"
	}
	config += "\">\n\t\t\t"
	if cst {
		config += "<const>"
	} else {
		config += "<string>"
	}
	config += opt
	if cst {
		config += "</const>"
	} else {
		config += "</string>"
	}
	config += "\n\t\t</edit>\n\t</match>\n"
	return config
}

func genBoolOptionConfig(verbosity int, opt bool, dbgOutput, comment, editName string, force bool) string {
	if strings.HasPrefix(editName, "force") && !opt {
		return ""
	}
	Dbg(verbosity, Debug, fmt.Sprintf(dbgOutput+" %t", opt))
	config := comment
	config += "\t<match target=\"pattern\">\n\t\t<edit name=\""
	config += editName
	config += "\" mode=\""
	if force {
		config += "assign"
	} else {
		config += "append"
	}
	config += "\">\n\t\t\t<bool>"
	config += fmt.Sprintf("%t", opt)
	config += "</bool>\n\t\t</edit>\n\t</match>\n"
	return config
}

func genUserInclude(userMode bool) string {
	if userMode {
		return "\t<include ignore_missing=\"yes\" prefix=\"xdg\">fontconfig/rendering-options.conf</include>\n"
	}
	return ""
}
