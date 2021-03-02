package ops

type Manager interface {
	AddManifest(repo, envName, appName string, manifest []byte)
	DeleteManifest(repo, envName, appName, name string)
}
