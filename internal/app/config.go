package app

type Config struct {
	Volumes []VolumeInfo
}

type VolumeInfo struct {
	Folder      string
	MountScript string
	AfterMount  string
}
