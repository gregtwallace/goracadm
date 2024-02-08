module github.com/gregtwallace/goracadm

go 1.19

replace github.com/gregtwallace/goracadm/cmd/goracadm-cert => /pkg/cmd/goracadm-cert

replace github.com/gregtwallace/goracadm/cmd/racadm => /pkg/cmd/racadm

replace github.com/gregtwallace/goracadm/pkg/app => /pkg/app

replace github.com/gregtwallace/goracadm/pkg/idrac => /pkg/idrac

require github.com/peterbourgon/ff/v4 v4.0.0-alpha.4
