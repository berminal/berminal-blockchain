package common

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPulicIP(t *testing.T) {
	Convey("public ip", t, func() {
		isPublic := IsPublicIP(net.ParseIP(GetPulicIP()))
		So(isPublic, ShouldBeTrue)
	})
}
