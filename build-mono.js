#!/usr/bin/env node

const fs = require("node:fs/promises");
const path = require("node:path");
const { exec } = require("child_process");
const prod = process.env.NODE_ENV === "production";
const buildTags = "mono" + (!prod ? " dev" : "");
const ldflags = prod ? `-ldflags="-s -w"` : "";
const autogeneratedHeader =
  "/**\nThis file was generated automatically in build-mono.js, do not edit nor commit into your repo.\n*/";

/**
This function reads the main.go file under the pluginPath.
Then it will check if the build tag "go:build !mono" exists. If not it will add the build tag.
It will also create a new file called "main_mono.go" where the contents are from the "main.go" file
but it will use another build tag "go:build mono" instead and it will replace the package name with camel-case
of the base name of the pluginPatha and returns the camel-case package name.
*/
async function preparePluginMain(pluginPath) {
  const mainpath = path.join(pluginPath, "main.go");
  const pathExists = await fileExists(mainpath);
  if (!pathExists) {
    return null;
  }

  const mainContent = await fs.readFile(mainpath, "utf8");
  const buildTagregex = /\s*\/\/\s*go:build\s+(.+)/;
  const tagMatches = mainContent.match(buildTagregex);
  const existingTags = tagMatches && tagMatches[1] ? tagMatches[1] : "";
  const hasBuildTags = buildTagregex.test(mainContent);
  const hasNotMonoTag = existingTags.includes("!mono");

  let newMainContent = mainContent.slice();
  if (!hasBuildTags) {
    newMainContent = `//go:build !mono\n\n${mainContent}`;
  } else {
    if (!hasNotMonoTag) {
      // add !mono tag
      const newBuildTags = existingTags + " !mono";
      newMainContent = mainContent.replace(
        buildTagregex,
        `//go:build ${newBuildTags}`
      );
    }
  }

  await fs.writeFile(path.join(pluginPath, "main.go"), newMainContent);
  console.log(`Created ${path.join(pluginPath, "main.go")}`);

  let monoMainContent = mainContent.slice();
  if (!hasBuildTags) {
    monoMainContent = `//go:build mono\n\n${mainContent}`;
  } else {
    const newBuildTags = existingTags.replace("!mono", "") + " mono";
    monoMainContent = mainContent.replace(
      buildTagregex,
      `//go:build ${newBuildTags}`
    );
  }

  const packageRegex = /package\s+(\w+)/;
  const monoPackageName = slugify(path.basename(pluginPath));
  monoMainContent = monoMainContent.replace(
    packageRegex,
    `package ${monoPackageName}`
  );

  // remove func main () {} if exists
  const funcMainRegex = /func\s+main\s*\(\s*\)\s*\{\s*\}/g;
  monoMainContent = monoMainContent.replace(funcMainRegex, "");
  monoMainContent = `${autogeneratedHeader}\n\n${monoMainContent}`;
  await fs.writeFile(path.join(pluginPath, "main_mono.go"), monoMainContent);
  console.log(`Created ${path.join(pluginPath, "main_mono.go")}`);

  return monoPackageName;
}

/**
---------------------------------------------------------------
helper functions
---------------------------------------------------------------
*/

function slugify(str) {
  return str
    .toLowerCase()
    .replace(/[^a-z0-9]/g, "")
    .replace(/-+/g, "")
    .replace(/^-|-$/g, "");
}

async function goModule(modulePath) {
  // Construct the full path to the go.mod file
  const goModFilePath = path.join(modulePath, "go.mod");

  // Read the contents of the go.mod file
  const goModContent = await fs.readFile(goModFilePath, "utf-8");

  // Use regular expressions to extract the module name
  const moduleNameMatch = goModContent.match(/module\s+([\w\/.-]+)/);

  if (moduleNameMatch && moduleNameMatch[1]) {
    return moduleNameMatch[1];
  } else {
    throw new Error("Unable to extract module name from go.mod file.");
  }
}

async function pluginPackage(pluginPath) {
  const jsonPath = path.join(pluginPath, "plugin.json");
  const pkg = await fs.readFile(jsonPath, "utf-8").then(JSON.parse);
  return pkg.package;
}

async function fileExists(filePath) {
  try {
    // Attempt to access the file
    await fs.access(filePath, fs.constants.F_OK);
    return true;
  } catch (error) {
    // File does not exist or there was an error accessing it
    return false;
  }
}

async function makePluginInit(pluginMonoModules) {
  // exclude core in plugin-init_mono.go
  const coreModule = pluginMonoModules.find((p) => p.go_package == "core");
  pluginMonoModules = pluginMonoModules.filter((p) => p.go_package != "core");

  const importModules = pluginMonoModules
    .map((p) => {
      return `${p.go_package} "${p.go_module}"`;
    })
    .join("\n  ");

  const pluginSwitches = pluginMonoModules
    .map((p) => {
      return `\tcase "${p.plugin_package}":\n\t\t${p.go_package}.Init(p)`;
    })
    .join("\n");

  const pluginMonoInit = `
//go:build mono

${autogeneratedHeader}

package plugins

import (
  "log"
  ${importModules}
)

func (p *PluginApi) Init() error {
\tswitch p.Pkg() {
\tcase "${coreModule.plugin_package}":
\t\tlog.Println("core package, skipping plugin.Init()...")
${pluginSwitches}
\tdefault:
\t\tlog.Println("Unable to load plugin: " + p.dir)
\t}
\treturn nil
}`;

  const pluginInitMonoPath = path.join("core/plugins/plugin-init_mono.go");
  await fs.writeFile(pluginInitMonoPath, pluginMonoInit.trim());
  console.log(`Created ${pluginInitMonoPath}`);
}

async function execPromise(cmd) {
  console.log("executing: ", cmd);
  await new Promise((resolve, reject) => {
    const proc = exec(cmd, (err) => {
      if (err) {
        reject(err);
      } else {
        resolve();
      }
    });

    proc.stdout.pipe(process.stdout);
    proc.stderr.pipe(process.stderr);
  });
}

/**
---------------------------------------------------------------
main function
---------------------------------------------------------------
*/

(async function () {
  const corePath = "core";
  const pluginsDir = "plugins";
  const userPlugins = await fs
    .readdir(pluginsDir)
    .then((paths) => paths.map((p) => path.join(pluginsDir, p)));

  const pluginsPaths = [corePath, ...userPlugins];

  // list of plugin module name, path and package name
  // e.g [{path: "/path/to/plugin", go_module: "github.com/my/plugin", go_package: "mono_package_name", plugin_package: "yaml_package"}]
  const pluginMonoModules = [];

  for (const pluginDir of pluginsPaths) {
    try {
      const goPackage = await preparePluginMain(pluginDir);
      const mod = await goModule(pluginDir);
      const pkg = await pluginPackage(pluginDir);

      pluginMonoModules.push({
        path: pluginDir,
        go_module: mod,
        go_package: goPackage,
        plugin_package: pkg
      });
    } catch (e) {
      console.log(e);
    }
  }

  await makePluginInit(pluginMonoModules);

  await execPromise(
    `./go-work.sh && cd main && go build --tags="${buildTags}" ${ldflags} -trimpath -o main.app main_mono.go`
  );

  console.log("Built app successfully: main/main.app");
})();
