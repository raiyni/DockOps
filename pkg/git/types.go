package git

// RepoConfig represents a configuration for a repo
type RepoConfig struct {
	URL         string `example:"https://github.com/portainer/portainer.git"`
	Ref         string `example:"refs/heads/branch_name"`
	ComposeFile string `example:"docker-compose.yml"`
	Hash        string `example:"bc4c183d756879ea4d173315338110b31004b8e0"`
}

type BaseOption struct {
	RepositoryUrl string
}

type BasicAuth struct {
	Username string
	Password string
}

type SshAuth struct {
	KeyPath     string
	KeyPassword string
}

type RepoAuth struct {
	BasicAuth
	SshAuth
}

type FetchOption struct {
	BaseOption
	RepoAuth
	ReferenceName string
}

type CloneOption struct {
	FetchOption
	Depth int
}
