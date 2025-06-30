package apacheBuildr

type ApacheBuildr struct {
	GroupId     string
	ArtifactId  string
	PackageType PackageType
	Classifier  string
	Version     string
}

type PackageType string

const (
	ZIP PackageType = "zip"
	TAR PackageType = "tar"
	NA  PackageType = "n/a"
)

func ValueOf(value string) PackageType {
	if value == "zip" {
		return ZIP
	} else if value == "tar" {
		return TAR
	} else {
		return NA
	}
}
