package resolve

import "fmt"

func genPackageNameServiceName(packageName string, serviceName string) string {
	return fmt.Sprintf("%v.%v", packageName, serviceName)
}
