package calldepth

import "k8s.io/klog/v2"

// Putting these functions into a separate file makes it possible to validate that
// their source code file is *not* logged because of WithCallDepth(1).

func myInfo(depth int, msg string) {
	klog.InfoDepth(depth+1, msg)
}

func myInfo2(msg string) {
	myInfo(1, msg)
}
