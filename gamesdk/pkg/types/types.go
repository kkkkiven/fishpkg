package types

import (
	"fmt"
)

//--------------DeployType--------------------
// DeployType 部署类型
type DeployType int8

var deployTypes = []string{"App", "VGame"}

func (d DeployType) String() string {
	return fmt.Sprintf("(%d)%s", d, deployTypes[d])
}
