package plugin

// IPluginMgr is used to get data of installed plugins in the system.
type IPluginMgr interface {

	// Find a plugin by name as defined in package.yml "name" field.
	FindByName(name string) IPluginApi

	// Find a plugin by path as defined in package.yml "package" field.
	FindByPkg(pkg string) IPluginApi

	// Returns all plugins installed in the system.
	All() []IPluginApi
}