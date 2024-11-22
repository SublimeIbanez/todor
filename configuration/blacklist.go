package configuration

import (
	"net/http"
	"net/url"
)

type BlacklistGroup []string

const TOPTAL_API_URL string = "https://www.toptal.com"
const TOPTAL_TYPE_PATH string = "developers/gitignore/api"

type Client struct {
	client    *http.Client
	UserAgent string
	ApiUrl    *url.URL
}

func NewClient() (*Client, error) {
	c := &Client{
		client:    http.DefaultClient,
		UserAgent: "todor",
	}

	api_url, err := url.Parse(TOPTAL_API_URL)
	if err != nil {
		return &Client{}, err
	}
	c.ApiUrl = api_url

	return c, nil
}

// https://www.toptal.com/developers/gitignore
// React
var React_Blacklist BlacklistGroup = []string{
	".DS_",
	"*.log",
	"logs",
	"**/*.backup.*",
	"**/*.back.*",
	"node_modules",
	"bower_components",
	"*.sublime*",
	"psd",
	"thumb",
	"sketch",
}

// NextJS
var NextJs_Blacklist BlacklistGroup = []string{
	// Dependencies
	"/node_modules",
	"/.pnp",
	".pnp.js",

	// Testing
	"/coverage",

	// Next.js
	"/.next/",
	"/out/",

	// Production
	"/build",

	// Misc
	".DS_Store",
	"*.pem",

	// Debug
	"npm-debug.log*",
	"yarn-debug.log*",
	"yarn-error.log*",
	".pnpm-debug.log*",

	// Local ENV Files
	".env*.local",

	// Vercel
	".vercel",

	// TypeScript
	"*.tsbuildinfo",
	"next-env.d.ts",
}

// C
var C_Blacklist BlacklistGroup = []string{
	// Prerequisites
	"*.d",

	// Object Files
	"*.o",
	"*.ko",
	"*.obj",
	"*.elf",

	// Linker output
	"*.ilk",
	"*.map",
	"*.exp",

	// Precompiled Headers
	"*.gch",
	"*.pch",

	// Libraries
	"*.lib",
	"*.a",
	"*.la",
	"*.lo",

	// Shared objects (including Windows DLLs)
	"*.dll",
	"*.so",
	"*.so.*",
	"*.dylib",

	// Executables
	"*.exe",
	"*.out",
	"*.app",
	"*.i*86",
	"*.x86_64",
	"*.hex",

	// Debug Files
	"*.dSYM/",
	"*.su",
	"*.idb",
	"*.pdb",

	// Kernel Module Compile Results
	"*.mod*",
	"*.cmd",
	".tmp_versions",
	"modules.order",
	"Module.symvers",
	"Mkfile.old",
	"dkms.conf",
}

var CMake_Blacklist BlacklistGroup = []string{
	// CMake
	"CMakeLists.txt.user",
	"CMakeCache.txt",
	"CMakeFiles",
	"CMakeScripts",
	"Testing",
	"Makefile",
	"cmake_install.cmake",
	"install_manifest.txt",
	"compile_commands.json",
	"CTestTestfile.cmake",
	"_deps",

	// External Projects
	"*-prefix/",
}

var CPP_Blacklist BlacklistGroup = []string{
	// Prerequisites
	"*.d",

	// Compiled Object files
	"*.slo",
	"*.lo",
	"*.o",
	"*.obj",

	// Precompiled Headers
	"*.gch",
	"*.pch",

	// Compiled Dynamic libraries
	"*.so",
	"*.dylib",
	"*.dll",

	// Fortran module files
	"*.mod",
	"*.smod",

	// Compiled Static libraries
	"*.lai",
	"*.la",
	"*.a",
	"*.lib",

	// Executables
	"*.exe",
	"*.out",
	"*.app",
}

var Java_Blacklist BlacklistGroup = []string{
	// Compiled class file
	"*.class",

	// Log file
	"*.log",

	// BlueJ files
	"*.ctxt",

	// Mobile Tools for Java (J2ME)
	".mtj.tmp/",

	// Package Files
	"*.jar",
	"*.war",
	"*.nar",
	"*.ear",
	"*.zip",
	"*.tar.gz",
	"*.rar",

	// virtual machine crash logs, see http://www.java.com/en/download/help/error_hotspot.xml
	"hs_err_pid*",
	"replay_pid*",
}
