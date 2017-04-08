//go:generate jsonenums -type=FormFieldType,LoginSearchType

package keepassrpc

/*

Most of this is translated directly from the KeePassRPC C# source:

https://github.com/kee-org/keepassrpc/blob/master/KeePassRPC/KeePassRPCService.cs

See any method tagged with `[JsonRpcMethod]` for specifics.

*/

// FormFieldType is the type of form entry field being referenced
type FormFieldType int

const (
	// FFTradio represents a radio-button selection
	FFTradio FormFieldType = iota

	// FFTusername represents a username entry
	FFTusername

	// FFTtext represents free-form text entry
	FFTtext

	// FFTpassword represents a password entry
	FFTpassword

	// FFTselect represents a drop-down selection
	FFTselect

	// FFTcheckbox represents a multiple-checkbox selection
	FFTcheckbox
)

// LoginSearchType is the type of login search being performed
type LoginSearchType int

const (
	// LSTall searches everything
	LSTall LoginSearchType = iota

	// LSTnoForms excludes forms
	LSTnoForms

	// LSTnoRealms excludes realms
	LSTnoRealms
)

const (
	// MatchAccuracyNone means no match (ie. asked to return all entries)
	MatchAccuracyNone = 0

	// MatchAccuracyDomain means same domain
	MatchAccuracyDomain = 10

	// MatchAccuracyHostname means same hostname
	MatchAccuracyHostname = 20

	// MatchAccuracyHostnameAndPort means same hostname and port
	MatchAccuracyHostnameAndPort = 30

	// MatchAccuracyClose means same URL excluding the query string
	MatchAccuracyClose = 40

	// MatchAccuracyBest means same URL including the query string (regex)
	MatchAccuracyBest = 50
)

// ApplicationMetadata describes the running instance of KeePass
type ApplicationMetadata struct {
	KeePassVersion string `json:"keePassVersion"`
	IsMono         bool   `json:"isMono"`
	NETCLR         string `json:"nETCLR"`
	NETversion     string `json:"nETversion"`
	MonoVersion    string `json:"monoVersion"`
}

// Configuration represents notable configuration of the KeePass instance
type Configuration struct {
	KnownDatabases []string `json:"knownDatabases"`
	AutoCommit     bool     `json:"autoCommit"`
}

// FormField represents a form for data entry
type FormField struct {
	Name        string        `json:"name"`
	DisplayName string        `json:"displayName"`
	Value       string        `json:"value"`
	Type        FormFieldType `json:"type"`
	ID          string        `json:"id"`
	Page        int           `json:"page"`
}

// LightEntry represents a single basic entry in the open KeePass database
type LightEntry struct {
	URLs          []string `json:"uRLs"`
	Title         string   `json:"title"`
	UniqueID      string   `json:"uniqueID"`
	UsernameValue string   `json:"usernameValue"`
	UsernameName  string   `json:"usernameName"`
	IconImageData string   `json:"iconImageData"`
}

// Entry describes a single complete entry in the open KeePass database
type Entry struct {
	*LightEntry

	HTTPRealm     string      `json:"hTTPRealm"`
	FormFieldList []FormField `json:"formFieldList"`

	// How accurately do the URLs in this entry match the URL we are looking for?
	// Higher = better match.
	// We don't consider protocol
	MatchAccuracy int `json:"matchAccuracy"`

	AlwaysAutoFill   bool `json:"alwaysAutoFill"`
	NeverAutoFill    bool `json:"neverAutoFill"`
	AlwaysAutoSubmit bool `json:"alwaysAutoSubmit"`
	NeverAutoSubmit  bool `json:"neverAutoSubmit"`
	// "KeeFox priority" = 1 (1 = 30000 relevancy score, 2 = 29999 relevancy score)
	// long autoTypeWhen "KeeFox config: autoType after page 2" (after/before or > / <) (page # or # seconds or #ms)
	// bool autoTypeOnly "KeeFox config: only autoType" This is probably redundant considering feature request #19?
	Priority int `json:"priority"`

	Parent Group    `json:"parent"`
	Db     Database `json:"db"`
}

// Group describes a group within the open KeePass database
type Group struct {
	Title         string `json:"title"`
	UniqueID      string `json:"uniqueID"`
	IconImageData string `json:"iconImageData"`
	Path          string `json:"path"`
}

// Database describes the currently-open database in KeePass
type Database struct {
	Name          string `json:"name"`
	FileName      string `json:"fileName"`
	Root          Group  `json:"root"`
	Active        bool   `json:"active"`
	IconImageData string `json:"iconImageData"`
}

// LaunchGroupEditor opens the editor on a given group
func (c *Client) LaunchGroupEditor(uuid, dbFileName string) error {
	return c.JSONRPCCtx.r.Call("LaunchGroupEditor",
		[]interface{}{uuid, dbFileName}, nil)
}

// LaunchLoginEditor opens the editor on a given login
func (c *Client) LaunchLoginEditor(uuid, dbFileName string) error {
	return c.JSONRPCCtx.r.Call("LaunchLoginEditor",
		[]string{uuid, dbFileName}, nil)
}

// GetCurrentKFConfig returns configuration information for the running KeePass
func (c *Client) GetCurrentKFConfig() (*Configuration, error) {
	var reply Configuration
	err := c.JSONRPCCtx.r.Call("GetCurrentKFConfig", nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// GetApplicationMetadata retrieves information about the running KeePass
func (c *Client) GetApplicationMetadata() (*ApplicationMetadata, error) {
	var reply ApplicationMetadata
	err := c.JSONRPCCtx.r.Call("GetApplicationMetadata", nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// GetDatabaseName retrieves the name of the currently open database
func (c *Client) GetDatabaseName() (string, error) {
	var reply string
	err := c.JSONRPCCtx.r.Call("GetDatabaseName", nil, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

// GetDatabaseFileName retrieves the filename of the currently open database
func (c *Client) GetDatabaseFileName() (string, error) {
	var reply string
	err := c.JSONRPCCtx.r.Call("GetDatabaseFileName", nil, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

// ChangeDatabase switches the active KeePass database
func (c *Client) ChangeDatabase(filename string, closeCurrent bool) error {
	return c.JSONRPCCtx.r.Call("ChangeDatabase",
		[]interface{}{filename, closeCurrent}, nil)
}

// ChangeLocation switches the active KeePass location
func (c *Client) ChangeLocation(locationID string) error {
	return c.JSONRPCCtx.r.Call("ChangeLocation", locationID, nil)
}

// GetPasswordProfiles retrieves a list of password profiles
func (c *Client) GetPasswordProfiles() ([]string, error) {
	var reply []string
	err := c.JSONRPCCtx.r.Call("GetPasswordProfiles", nil, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GeneratePassword asks KeePass to generate a new password
func (c *Client) GeneratePassword(profileName, url string) (string, error) {
	var reply string
	err := c.JSONRPCCtx.r.Call("GeneratePassword",
		[]string{profileName, url}, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

// RemoveEntry removes a specified entry from the active KeePass database
func (c *Client) RemoveEntry(uuid string) (bool, error) {
	var reply bool
	err := c.JSONRPCCtx.r.Call("RemoveEntry", uuid, &reply)
	if err != nil {
		return false, err
	}
	return reply, nil
}

// RemoveGroup removes a specified entry from the active KeePass database
func (c *Client) RemoveGroup(uuid string) (bool, error) {
	var reply bool
	err := c.JSONRPCCtx.r.Call("RemoveGroup", uuid, &reply)
	if err != nil {
		return false, err
	}
	return reply, nil
}

// AddLogin adds a new login to the database
func (c *Client) AddLogin(login *Entry, parentUUID, dbFileName string) (*Entry, error) {
	var reply Entry
	err := c.JSONRPCCtx.r.Call("AddLogin",
		[]interface{}{login, parentUUID, dbFileName}, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// AddGroup adds a new group to the database
func (c *Client) AddGroup(name, parentUUID string) (*Group, error) {
	var reply Group
	err := c.JSONRPCCtx.r.Call("AddGroup", []string{name, parentUUID}, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// UpdateLogin updates an existing login in the database
func (c *Client) UpdateLogin(login *Entry, oldLoginUUID string, urlMergeMode int, dbFileName string) (*Entry, error) {
	var reply Entry
	err := c.JSONRPCCtx.r.Call("UpdateLogin",
		[]interface{}{login, oldLoginUUID, urlMergeMode, dbFileName},
		&reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// GetParent retrieves the parent group of a specified group
func (c *Client) GetParent(uuid string) (*Group, error) {
	var reply Group
	err := c.JSONRPCCtx.r.Call("GetParent", uuid, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// GetRoot retrieves the root group of the database
func (c *Client) GetRoot() (*Group, error) {
	var reply Group
	err := c.JSONRPCCtx.r.Call("GetRoot", nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// GetAllDatabases returns all of the available KeePass databases
func (c *Client) GetAllDatabases(fullDetails bool) ([]Database, error) {
	var reply []Database
	err := c.JSONRPCCtx.r.Call("GetAllDataases", fullDetails, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetAllLogins retrieves all logins in the database
func (c *Client) GetAllLogins() ([]Entry, error) {
	var reply []Entry
	err := c.JSONRPCCtx.r.Call("GetAllLogins", nil, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetChildEntries returns all entries under a specified parent
func (c *Client) GetChildEntries(uuid string) ([]Entry, error) {
	var reply []Entry
	err := c.JSONRPCCtx.r.Call("GetChildEntries", uuid, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetChildGroups returns all groups under a specified parent
func (c *Client) GetChildGroups(uuid string) ([]Group, error) {
	var reply []Group
	err := c.JSONRPCCtx.r.Call("GetChildGroups", uuid, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// FindGroups searches the database for groups matching a pattern
// THIS DOES NOT APPEAR TO WORK. The `out Group[]` parameter doesn't make
// much sense exported via JSON-RPC. Original C# method signature:
//
// public int FindGroups(string name, string uuid, out Group[] groups)
func (c *Client) FindGroups(name, uuid string) (int, error) {
	var reply int
	err := c.JSONRPCCtx.r.Call("FindGroups", []interface{}{name, uuid, nil}, &reply)
	if err != nil {
		return -1, err
	}
	return reply, nil
}

// FindLogins searches the database for logins matching a pattern
func (c *Client) FindLogins(unsanitizedURLs []string, actionURL, httpRealm string, lst LoginSearchType, requireFullURLMatches bool, uniqueID, dbFileName, freeTextSearch, username string) ([]Entry, error) {
	args := []interface{}{
		unsanitizedURLs,
		actionURL,
		httpRealm,
		lst,
		requireFullURLMatches,
		uniqueID,
		dbFileName,
		freeTextSearch,
		username,
	}
	var reply []Entry
	err := c.JSONRPCCtx.r.Call("FindLogins", args, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// CountLogins returns the number of logins that match a pattern
func (c *Client) CountLogins(URL, actionURL, httpRealm string, lst LoginSearchType, requireFullURLMatches bool) (int, error) {
	args := []interface{}{
		URL,
		actionURL,
		httpRealm,
		lst,
		requireFullURLMatches,
	}
	var reply int
	err := c.JSONRPCCtx.r.Call("CountLogins", args, &reply)
	if err != nil {
		return -1, err
	}
	return reply, nil
}

/* These methods are automatically generated by Jayrock */

// SystemListMethods (system.listMethods) returns all available methods
func (c *Client) SystemListMethods() ([]string, error) {
	var reply []string
	err := c.JSONRPCCtx.r.Call("system.listMethods", nil, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// SystemVersion (system.version) returns the server version information
func (c *Client) SystemVersion() (string, error) {
	var reply string
	err := c.JSONRPCCtx.r.Call("system.version", nil, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

// SystemAbout (system.about) returns a summary of information about the service
func (c *Client) SystemAbout() (string, error) {
	var reply string
	err := c.JSONRPCCtx.r.Call("system.about", nil, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}
