package subcmd

import "github.com/fatih/color"

// Standard welcome banner.
func Banner() string {
	s := "\n" +
		"                 o          o                 o          8\n" +
		"                            8                 8          8\n" +
		" ooYoYo. .oPYo. o8  .oPYo. o8P .oPYo. .oPYo. o8P .oPYo.  8\n" +
		" 8' 8  8 8    8  8  Yb..    8  8    8 8oooo8  8  .oooo8  8\n" +
		" 8  8  8 8.   8  8    'Yb.  8  8    8 8.      8  8    8  8\n" +
		" 8  8  8 `Yooo'  8  `YooP'  8  8YooP' `Yooo'  8  `YooP8  8\n" +
		"                               8                          \n" +
		"                               8                          \n"
	return color.GreenString(s)
}

// Standard exit banner.
func Farewell() string {
	s := "" +
		"     _ _  __                                   __  _ _    \n" +
		"    ( | )/_/                                   \\_\\( | )   \n" +
		" __( >O< )               " + color.RedString(".Goodbye.") + "               ( >O< )__\n" +
		" \\_\\(_|_)                                         (_|_)/_/\n"
	return s
}
