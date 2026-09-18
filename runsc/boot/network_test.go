// Copyright 2026 The gVisor Authors.
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

package boot

import (
	"testing"
)

func TestPrimaryInterface(t *testing.T) {
	testCases := []struct {
		name string
		args *CreateLinksAndRoutesArgs
		want string
	}{
		{
			name: "nil args",
			args: nil,
			want: "",
		},
		{
			name: "v4 gateway in fdbased",
			args: &CreateLinksAndRoutesArgs{
				Defaultv4Gateway: DefaultRoute{Name: "eth1"},
				FDBasedLinks:     []FDBasedLink{{Name: "eth0"}, {Name: "eth1"}},
			},
			want: "eth1",
		},
		{
			name: "v4 gateway not in fdbased with eth0 fallback",
			args: &CreateLinksAndRoutesArgs{
				Defaultv4Gateway: DefaultRoute{Name: "other0"},
				FDBasedLinks:     []FDBasedLink{{Name: "eth0"}, {Name: "eth1"}},
			},
			want: "eth0",
		},
		{
			name: "v6 gateway in fdbased",
			args: &CreateLinksAndRoutesArgs{
				Defaultv6Gateway: DefaultRoute{Name: "eth2"},
				FDBasedLinks:     []FDBasedLink{{Name: "eth1"}, {Name: "eth2"}},
			},
			want: "eth2",
		},
		{
			name: "v6 gateway not in fdbased, single link fallback",
			args: &CreateLinksAndRoutesArgs{
				Defaultv6Gateway: DefaultRoute{Name: "other0"},
				FDBasedLinks:     []FDBasedLink{{Name: "myif0"}},
			},
			want: "myif0",
		},
		{
			name: "single link in fdbased",
			args: &CreateLinksAndRoutesArgs{
				FDBasedLinks: []FDBasedLink{{Name: "tap0"}},
			},
			want: "tap0",
		},
		{
			name: "multiple links with no gateway and no eth0",
			args: &CreateLinksAndRoutesArgs{
				FDBasedLinks: []FDBasedLink{{Name: "tap0"}, {Name: "tap1"}},
			},
			want: "",
		},
		{
			// The gateway resolves to an XDP link, which cannot be proxied over
			// UDS. hasFDBasedLink() must reject it rather than returning a name
			// that has no corresponding FDBasedLink.
			name: "v4 gateway matches an XDP link only",
			args: &CreateLinksAndRoutesArgs{
				Defaultv4Gateway: DefaultRoute{Name: "xdp0"},
				XDPLinks:         []XDPLink{{Name: "xdp0"}},
				FDBasedLinks:     []FDBasedLink{{Name: "tap0"}, {Name: "tap1"}},
			},
			want: "",
		},
		{
			// Same, but eth0 is available as a fallback.
			name: "v4 gateway matches an XDP link with eth0 fallback",
			args: &CreateLinksAndRoutesArgs{
				Defaultv4Gateway: DefaultRoute{Name: "xdp0"},
				XDPLinks:         []XDPLink{{Name: "xdp0"}},
				FDBasedLinks:     []FDBasedLink{{Name: "eth0"}, {Name: "eth1"}},
			},
			want: "eth0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PrimaryInterface(tc.args)
			if got != tc.want {
				t.Errorf("PrimaryInterface(%+v) = %q, want %q", tc.args, got, tc.want)
			}
		})
	}
}
