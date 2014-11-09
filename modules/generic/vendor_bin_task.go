package generic

import "fmt"

type VendorBinTask struct {
	sourcePath string
}

func NewVendorBinTask(sourcePath string) *VendorBinTask {
	return &VendorBinTask{"vendor/src/" + sourcePath}
}

func (t *VendorBinTask) Dependencies() []string {
	return []string{t.sourcePath}
}

func (t *VendorBinTask) Invoke(name string) error {
	fmt.Printf("Invoking task %s with source %s\n", name, t.sourcePath)
	return nil
}
