package main

import "testing"

func TestName(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"https://github.com/Canop/broot/releases/download/v1.6.3/broot_1.6.3.zip", "broot"},
		{"https://releases.hashicorp.com/terraform/1.0.4/terraform_1.0.4_linux_amd64.zip", "terraform"},
		{"https://get.helm.sh/helm-v3.6.3-linux-amd64.tar.gz", "helm"},
		{"https://github.com/golangci/golangci-lint/releases/download/v1.42.0/golangci-lint-1.42.0-linux-amd64.tar.gz", "golangci-lint"},
		{"https://github.com/ogham/exa/releases/download/v0.10.1/exa-linux-x86_64-v0.10.1.zip", "exa"},
		{"https://github.com/XAMPPRocky/tokei/releases/download/v12.1.2/tokei-mips64-unknown-linux-gnuabi64.tar.gz", "tokei"},
		{"https://github.com/XAMPPRocky/tokei/releases/download/v12.1.2/tokei-x86_64-unknown-linux-gnu.tar.gz", "tokei"},
		{"https://github.com/boyter/scc/releases/download/v3.0.0/scc-3.0.0-x86_64-unknown-linux.zip", "scc"},
	}

	for _, tt := range cases {
		got := name(tt.url)
		if got != tt.want {
			t.Errorf("base(%v): want %v, got %v", tt.url, tt.want, got)
		}
	}
}
