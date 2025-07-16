package parser

type DeviceProp struct {
	Model          string `json:"model"`
	Brand          string `json:"brand"`
	Name           string `json:"name"`
	AndroidVersion string `json:"android_version"`
	Fingerprint    string `json:"fingerprint"`
	SDK            string `json:"sdk"`
	SecurityPatch  string `json:"security_patch"`
	Processor      string `json:"processor"`
	ABI            string `json:"abi"`
}

func NewDeviceProp() IParser {
	return &DeviceProp{}
}

func (d *DeviceProp) Parse(rawData string) error {
	propsData := parseListData(rawData, []string{
		"[ro.product.model]",
		"[ro.product.brand]",
		"[ro.product.name]",
		"[ro.product.build.fingerprint]",
		"[ro.build.fingerprint]",
		"[ro.build.version.release]",
		"[ro.build.version.sdk]",
		"[ro.build.version.security_patch]",
		"[ro.product.board]",
		"[ro.product.cpu.abi]",
	}, true)
	d.Model = getParseListDataValue(propsData, "[ro.product.model]")
	d.Brand = getParseListDataValue(propsData, "[ro.product.brand]")
	d.Name = getParseListDataValue(propsData, "[ro.product.name]")
	d.Fingerprint = getParseListDataValue(propsData, "[ro.product.build.fingerprint]")
	if d.Fingerprint == "" {
		d.Fingerprint = getParseListDataValue(propsData, "[ro.build.fingerprint]")
	}
	d.AndroidVersion = getParseListDataValue(propsData, "[ro.build.version.release]")
	d.SDK = getParseListDataValue(propsData, "[ro.build.version.sdk]")
	d.SecurityPatch = getParseListDataValue(propsData, "[ro.build.version.security_patch]")
	d.Processor = getParseListDataValue(propsData, "[ro.product.board]")
	d.ABI = getParseListDataValue(propsData, "[ro.product.cpu.abi]")
	return nil
}
