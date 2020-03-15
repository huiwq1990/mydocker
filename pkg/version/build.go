// Copyright (c) 2017 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

// 通过makefile中的ldflagers动态设置值 https://segmentfault.com/a/1190000008323048
var (
	// commitFromGit is a constant representing the source version that
	// generated this build. It should be set during build via -ldflags.
	commitSHA string
	// versionFromGit is a constant representing the version tag that
	// generated this build. It should be set during build via -ldflags.
	latestVersion string
	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	date string
)

// Info holds build information
type Info struct {
	GitCommit  string `json:"gitCommit"`
	GitVersion string `json:"GitVersion"`
	BuildDate  string `json:"BuildDate"`
}

// Get creates and initialized Info object
func Get() Info {
	return Info{
		GitCommit:  commitSHA,
		GitVersion: latestVersion,
		BuildDate:  date,
	}
}