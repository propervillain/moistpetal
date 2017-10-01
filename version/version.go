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

//go:generate go run make_version.go

package version

import (
	"fmt"
	"time"
)

// These strings will be overwritten by an init function in
// created by make_version.go during the release process.
var (
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = 0
	BuildTime    time.Time
	CommitSHA    = "dev"
	TimeFormat   = "02-Jan-06"
)

// Version displays the moistpetal code version.
func Version() string {
	s := fmt.Sprintf("%d.%d.%d+%s (%s)",
		MajorVersion,
		MinorVersion,
		PatchVersion,
		CommitSHA,
		BuildTime.Format(TimeFormat),
	)
	return s
}
