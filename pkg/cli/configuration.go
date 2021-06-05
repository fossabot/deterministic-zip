package cli

import (
	"errors"
	flag "github.com/spf13/pflag"
)

// Configuration represents the config for the cli and may be mutated by features
type Configuration struct {
	// ZipFile is the target zip file name
	ZipFile string

	// SourceFiles contains a flat list with paths that are either fully qualified or based on the pwd
	SourceFiles []string

	// Verbose outputs very detailed information
	Verbose bool

	// Recursive includes all child folders automatically
	Recursive bool

	// Exclude file patterns from the archive
	Exclude []string

	// NoDirEntries states if no directories should be created inside the zip
	NoDirEntries bool

	// Include patterns
	Include []string

	// CompressionMethod to use for the zip archive
	CompressionMethod string

	// Quiet eliminates all outputs
	Quiet bool

	// LogFilePath contains the path to the log
	LogFilePath string

	// LogFileAppend specifies if the file should be appended
	LogFileAppend bool
}

func (conf *Configuration) addBoolFlag(field *bool, long string, short string, val bool, usage string) {
	flag.BoolVarP(field, long, short, val, usage)
}

func (conf *Configuration) addStringsFlag(field *[]string, long string, short string, val []string, usage string) {
	if short == "" {
		flag.StringSliceVar(field, long, val, usage)
	} else {
		flag.StringSliceVarP(field, long, short, val, usage)
	}
}

func (conf *Configuration) addStringFlag(field *string, long string, short string, val string, usage string) {
	flag.StringVarP(field, long, short, val, usage)
}

func (conf *Configuration) defineFlags() {
	conf.addBoolFlag(&conf.Verbose, "verbose", "v", false, "Verbose mode or print diagnostic version info.")
	conf.addBoolFlag(&conf.Recursive, "recurse-paths", "r", false, "Include all files verbose")
	conf.addBoolFlag(&conf.NoDirEntries, "no-dir-entries", "D", false, "Do not create entries in the zip archive for directories. Directory entries are created by default.")
	conf.addBoolFlag(&conf.Quiet, "quiet", "q", false, "Quiet mode; eliminate informational messages")
	conf.addStringsFlag(&conf.Exclude, "exclude", "", []string{}, "Exclude specific file patterns")
	conf.addStringsFlag(&conf.Include, "include", "i", []string{}, "Include only the specified file pattern")
	conf.addStringFlag(&conf.CompressionMethod, "compression-method", "Z", "deflate", "Set the default compression method. \nCurrently the main methods supported by zip are store and deflate. \nCompression method can be set to:\n\nstore       Setting the compression method to store forces to store entries with no compression. \n            This is generally faster than compressing entries, but results in no space savings.\n\ndeflate     This is the default method for zip. If zip determines that storing is better than deflation, the entry will be stored instead.\n")
	conf.addStringFlag(&conf.LogFilePath, "logfile-path", "", "", "Open a logfile at the given path. By default any existing file at that location is overwritten, but the -la option will result in an existing file being opened and the new log information appended to any existing information. Only warnings and errors are written to the log unless the -li option is also given, then all information messages are also written to the log.")
	conf.addBoolFlag(&conf.LogFileAppend, "log-append", "", false, "Append to existing logfile. Default is to overwrite.")
}

func (conf *Configuration) parseVarargs() error {
	remaining := flag.Args()
	if len(remaining) < 2 {
		return errors.New("specify at least the destination package and source files")
	}

	conf.ZipFile = remaining[0]
	conf.SourceFiles = remaining[1:]
	return nil
}

func (conf *Configuration) help() {
	PrintCompactInfo()
	flag.PrintDefaults()
}

// Parse the configuration from cli args
func (conf *Configuration) Parse() error {
	conf.defineFlags()

	isHelp := flag.BoolP("help", "h", false, "Show available commands")
	isVersion := flag.Bool("version", false, "Show version info")
	flag.Parse()

	if *isHelp {
		conf.help()
		return ErrAbort
	} else if *isVersion {
		PrintVersionInfo()
		return ErrAbort
	}

	return conf.parseVarargs()
}

// NewConfiguration creates a new configuration
func NewConfiguration() *Configuration {
	return &Configuration{}
}
