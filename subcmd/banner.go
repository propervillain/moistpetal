// Copyright 2017, The Moistpetal Authors. All rights reserved.
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

package subcmd

import "github.com/fatih/color"

// Banner prints the standard moistpetal banner.
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

// Farewell prints the standard moistpetal farewell message.
func Farewell() string {
	g := color.RedString(".Goodbye.")
	s := "" +
		"     _ _  __                                   __  _ _    \n" +
		"    ( | )/_/                                   \\_\\( | )   \n" +
		" __( >O< )               " + g + "               ( >O< )__\n" +
		" \\_\\(_|_)                                         (_|_)/_/\n"
	return s
}
