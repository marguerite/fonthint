package lib

import (
	"log"
	"strconv"
	"strings"
)

// genFcPreamble generate fontconfig preamble
func genFcPreamble(userMode bool, comment string) string {
	config := "<?xml version=\"1.0\"?>\n<!DOCTYPE fontconfig SYSTEM \"fonts.dtd\">\n\n<!-- DO NOT EDIT; this is a generated file -->\n<!-- modify "
	config += GetConfigLocation("fc", false)
	config += " && run /usr/bin/fonts-config "
	if userMode {
		config += "-\\-user "
	}
	config += "instead. -->\n"
	config += comment
	config += "\n<fontconfig>\n\n"
	return config
}

//genFontTypeByHinting generate fontconfig font_type block based on tt hinting.
func genFontTypeByHinting(name string, hinting bool) string {
	config := "\t<match target=\"font\">\n\t\t<test name=\"family\">\n\t\t\t<string>" + name + "</string>\n\t\t</test>\n"
	config += "\t\t<edit name=\"font_type\" mode=\"assign\">\n\t\t\t<string>"
	if hinting {
		config += "TT Instructed Font"
	} else {
		config += "NON TT Instructed Font"
	}
	config += "</string>\n\t\t</edit>\n\t</match>\n\n"
	return config
}

func genBlacklistConfig(f Font) string {
	config := "\t<match target=\"scan\">\n\t\t<test name=\"family\">\n\t\t\t<string>" + f.Name[0] + "</string>\n\t\t</test>\n"
	if !(f.Width == 0 && f.Weight == 0 && f.Slant == 0) {
		if f.Width != 100 {
			config += "\t\t<test name=\"width\">\n\t\t\t<int>" + strconv.Itoa(f.Width) + "</int>\n\t\t</test>\n"
		}
		if f.Weight != 80 {
			config += "\t\t<test name=\"weight\">\n\t\t\t<int>" + strconv.Itoa(f.Weight) + "</int>\n\t\t</test>\n"
		}
		if f.Slant != 0 {
			config += "\t\t<test name=\"slant\">\n\t\t\t<int>" + strconv.Itoa(f.Slant) + "</int>\n\t\t</test>\n"
		}
	}
	config += "\t\t<edit name=\"charset\" mode=\"assign\">\n\t\t\t<minus>\n\t\t\t\t<name>charset</name>\n"
	config += genCharsetConfig(f.Charset)
	config += "\t\t\t</minus>\n\t\t</edit>\n\t</match>\n\n"
	return config
}

// genCharsetConfig convert Charset to fontconfig conf
func genCharsetConfig(c Charset) string {
	config := "\t\t\t\t<charset>\n"
	for _, v := range c {
		if strings.Contains(v, "..") {
			config += "\t\t\t\t\t<range>\n"
			for _, s := range strings.Split(v, "..") {
				config += "\t\t\t\t\t\t<int>0x" + s + "</int>\n"
			}
			config += "\t\t\t\t\t</range>\n"
		} else {
			config += "\t\t\t\t\t<int>0x" + v + "</int>\n"
		}
	}
	config += "\t\t\t\t</charset>\n"
	return config
}

func genDualAisanConfig(f Font) string {
	config := ""
	for _, name := range f.Name {
		config += "\t<match target=\"font\">\n\t\t<test name=\"family\" compare=\"contains\">\n\t\t\t<string>"
		config += name
		config += "</string>\n\t\t</test>\n"
		config += "\t\t<edit name=\"spacing\" mode=\"append\">\n\t\t\t<const>proportional</const>\n\t\t</edit>\n"
		config += "\t\t<edit name=\"globaladvance\" mode=\"append\">\n\t\t\t<bool>false</bool>\n\t\t</edit>\n\t</match>\n\n"
	}
	return config
}

func genCJKMatrixConfig(fontName string, matrix []float64, nameLangs []string, fonts Collection) string {
	config := ""

	if len(matrix) != 4 {
		log.Fatalf("Invalid matrix: %v", matrix)
	}

	// generate nothing for non-installed fonts
	if len(fonts.FindByName(fontName)) <= 0 {
		return config
	}

	for _, nameLang := range nameLangs {
		s := "\t<match target=\"font\">\n\t\t<test name=\"family\">\n\t\t\t<string>" + fontName + "</string>\n\t\t</test>\n"
		s += "\t\t<test name=\"namelang\">\n\t\t\t<string>" + nameLang + "</string>\n\t\t</test>\n"
		s += "\t\t<edit name=\"matrix\" mode=\"assign\">\n\t\t\t<times>\n\t\t\t\t<name>matrix</name>\n\t\t\t\t<matrix>\n"
		s += "\t\t\t\t\t<double>" + strconv.FormatFloat(matrix[0], 'f', -1, 64) + "</double>\n"
		s += "\t\t\t\t\t<double>" + strconv.FormatFloat(matrix[1], 'f', -1, 64) + "</double>\n"
		s += "\t\t\t\t\t<double>" + strconv.FormatFloat(matrix[2], 'f', -1, 64) + "</double>\n"
		s += "\t\t\t\t\t<double>" + strconv.FormatFloat(matrix[3], 'f', -1, 64) + "</double>\n"
		s += "\t\t\t\t</matrix>\n\t\t\t</times>\n\t\t</edit>\n\t</match>\n\n"
		config += s
	}
	return config
}

func genCJKWeightConfig(fontName string, weights [][]int, nameLangs []string, fonts Collection) string {
	config := ""

	for _, w := range weights {
		if len(w) < 3 {
			log.Fatalf("invalid weight item: %v", w)
		}
	}

	if len(fonts.FindByName(fontName)) <= 0 {
		return config
	}

	for _, nameLang := range nameLangs {
		for _, w := range weights {
			s := "\t<match target=\"font\">\n\t\t<test name=\"family\">\n\t\t\t<string>" + fontName + "</string>\n"
			s += "\t\t</test>\n\t\t<test name=\"namelang\">\n\t\t\t<string>" + nameLang + "</string>\n\t\t</test>\n"

			if w[0] != 0 {
				s += "\t\t<test name=\"weight\" compare=\"more_eq\">\n\t\t\t<int>" + strconv.FormatInt(int64(w[0]), 10) + "</int>\n\t\t</test>\n"
			}

			if w[1] != 0 {
				s += "\t\t<test name=\"weight\" compare=\"less_eq\">\n\t\t\t<int>" + strconv.FormatInt(int64(w[1]), 10) + "</int>\n\t\t</test>\n"
			}

			s += "\t\t<edit name=\"weight\" mode=\"assign\">\n\t\t\t<int>" + strconv.FormatInt(int64(w[2]), 10) + "</int>\n\t\t</edit>\n\t</match>\n\n"
			config += s
		}
	}
	return config
}

func genCJKWidthConfig(fontName string, widths []int, nameLangs []string, fonts Collection) string {
	config := ""

	if len(widths) != 2 {
		log.Fatalf("invalid weight item: %v", widths)
	}

	if len(fonts.FindByName(fontName)) <= 0 {
		return config
	}

	for _, nameLang := range nameLangs {
		s := "\t<match target=\"font\">\n\t\t<test name=\"family\">\n\t\t\t<string>" + fontName + "</string>\n\t\t</test>\n"
		s += "\t\t<test name=\"namelang\">\n\t\t\t<string>" + nameLang + "</string>\n\t\t</test>\n"
		s += "\t\t<test name=\"width\" compare=\"more_eq\">\n\t\t\t<int>" + strconv.FormatInt(int64(widths[0]), 10) + "</int>\n\t\t</test>\n"
		s += "\t\t<test name=\"width\" compare=\"less_eq\">\n\t\t\t<int>" + strconv.FormatInt(int64(widths[1]), 10) + "</int>\n\t\t</test>\n"
		s += "\t\t<edit name=\"width\" mode=\"assign\">\n\t\t\t<int>" + strconv.FormatInt(int64(widths[0]), 10) + "</int>\n\t\t</edit>\n\t</match>\n\n"
		config += s
	}
	return config
}
