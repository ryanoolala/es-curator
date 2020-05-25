package main

// Configurations exported
type Configurations struct {
	ESHost   string
	Snapshot SnapshotConfigurations
	Index    []Indexes
}

// Indexes exported
type Indexes struct {
	name         string
	monthsWindow int16
}

// SnapshotConfigurations exported
type SnapshotConfigurations struct {
	name     string
	repoType string
	settings SnapshotSettings
}

// SnapshotSettings exported
type SnapshotSettings struct {
	bucket  string
	Region  string
	roleArn string
}
