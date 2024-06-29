// Written by Mikhail P. Ortiz-Lunyov
//
// Version 0.0.0.8-alpha
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// This script is cross platform.

package main

// Import packages
import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"time"
)

// Script-level fields
// // Operating system
const OS_TYPE string = runtime.GOOS

// // Version numbers and names
const DEV_CYCLE string = "-beta"
const SHORT_VERSION_NUM string = "0.0.0.9"
const VERSION_NAME string = "June 28th 2024"
const LONG_VERSION_NUM string = "v" + SHORT_VERSION_NUM + DEV_CYCLE + " (" + VERSION_NAME + ")"

// // Number of package managers per type
const OF_PKG_NUM int = 15
const AL_PKG_NUM int = 4

// // String array of package managers
var OFFICIAL_PKG_MANAGERS [OF_PKG_NUM]string = [OF_PKG_NUM]string{
	// Linux
	"apt",                  // 0  [Verified] Debian
	"dnf",                  // 1  [Verified] Red-Hat
	"transactional-update", // 2  [Verified*]OpenSUSE immutable
	"zypper",               // 3  [Verified**]OpenSUSE
	"yum",                  // 4  [Verified] Legacy Red-Hat
	"rpm-ostree",           // 5  [Verified] Red-Hat immutable
	"apk",                  // 6  [Verified] Alpine Linux
	"swupd",                // 7  [Verified] Clear Linux
	"pacman",               // 8  [] Arch Linux
	"Pkg_add",              // 9  [] OpenBSD
	"pkg",                  // 10 [] FreeBSD
	"eopkg",                // 11 [] Solus Linux
	"slackpkg",             // 12 [] Slackware Linux
	"xpbs",                 // 13 [] Void Linux
	// *Currently does not work on non-root execution
	// **Is NOT detected on non-root execution on OpenSUSE MicroOS, fails due to transactional-update

	// Windows
	"winget", // 14 [Verified***] Winget
	// ***Does not work on first-time execution. Needs "y" piped in first ["y" | winget upgrade --all]
	// ***Additionally, non-admin execution requires user to be present to approve admin prompts
}
var ALTERNATIVE_PKG_MANAGERS [AL_PKG_NUM]string = [AL_PKG_NUM]string{
	"brew",    // 0 [Verified] Homebrew
	"snap",    // 1 [] Snap
	"choco",   // 2 [Verified*] Chocolatey
	"flatpak", // 3 [Verified] Flatpak
	// *Is NOT detected on non-root execution
}

// // Critical variables
var rootUse string
var debugFlag bool = true

// Prints Exit Statement
func ExitStatement() {
	fmt.Println("\n\t* I hope this program was useful for you!")
	fmt.Println("\t* Please give this project a star on GitHub!")
}

// Prints Version statement
func PrintVersion() {
	fmt.Println(" = = =")
	fmt.Println("Update_Full-GO " + LONG_VERSION_NUM)
}

// Prints Flags statement
func PrintFlags(verbosity int) {
	// Print intro
	if verbosity >= 2 {
		fmt.Println("There are two types flags available: informational, and functional")
	}
	// Informational flags
	if verbosity >= 1 {
		fmt.Println("\tInformational (overrides all functional flags):")
	}
	fmt.Println("--help     | -h : Prints this help message")
	fmt.Println("--version  | -v : Prints version statement")
	fmt.Println("--flags    | -f : Prints all available flags (default verbosity 0)")
	fmt.Println("--warranty | -w : Prints the warranty seciton from the GNU Public License v3")
	fmt.Println("--debug    | -d : Prints more verbose technical output for debugging")
	// Functional flags
	if verbosity >= 1 {
		fmt.Println("\tFunctional:")
	}
	fmt.Println("--manual-all | -ma : Makes user manually select options presented by")
	fmt.Println("--alt-only   | -ao : Only updates from alternative package managers (see definition)")
	fmt.Println("--custom-domain | -cd : Adds an additional domain to test on top of raw.githubusercontent.com")
	fmt.Println("--official-only | -oo : Only updates from official package managers (see definition)")
	fmt.Println("--yum-update | -yu : Uses Yum over Dnf, if exists or is applicable")
}

// Prints Help statement
func PrintHelp(flagVerbosity int) {
	// Print version
	PrintVersion()
	fmt.Println(" = = =")
	fmt.Println("This Go script allows for Full updates on a variety of OSs, including Linux, Windows, and other flavours of UNIX")
	// Begin describing available flags
	PrintFlags(flagVerbosity)
	fmt.Println("Exit codes:")
	fmt.Println("0: Successful operation of script")
	fmt.Println("1: Error on behalf of USER")
	fmt.Println("3: Error on behalf of DEVELOPER")
	fmt.Println("4: Other Error (environmental, incompatible, etc)")
	fmt.Println("130: Cancelled by USER")
	fmt.Println()
}

// Prints Warranty statement
func PrintWarranty() {
	fmt.Println(" = = =")
	fmt.Println("An excerpt from the GNU Public License v3, regarding any warranty of this program:")
	fmt.Println("THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY APPLICABLE LAW.")
	fmt.Println("EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT HOLDERS AND/OR OTHER PARTIES")
	fmt.Println("PROVIDE THE PROGRAM “AS IS” WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESSED OR")
	fmt.Println("IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND")
	fmt.Println("FITNESS FOR A PARTICULAR PURPOSE. THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE")
	fmt.Println("OF THE PROGRAM IS WITH YOU. SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE")
	fmt.Println("COST OF ALL NECESSARY SERVICING, REPAIR OR CORRECTION.")
}

// Prints values of variables
func DebugVariablePrint(varName string, isBool bool, boolVar bool, intVar int, stringVar string, errVar error, byteVar []byte, response *http.Response) {
	switch debugFlag {
	// End method early to not spend any more time
	case false:
		return
	case true:
		fmt.Print("DEBUG= ")
		fmt.Print(varName + ": ")
		// Define what to print in event of different variables
		switch isBool {
		case true:
		default:
			fmt.Print(boolVar)
		}
		switch intVar {
		case -1:
		default:
			fmt.Print(intVar)
		}
		switch stringVar {
		case "null":
		default:
			fmt.Print(stringVar)
		}
		switch errVar {
		case nil:
		default:
			fmt.Println(errVar)
		}
		switch byteVar {
		case nil:
		default:
			fmt.Print(byteVar)
		}
		switch response {
		case nil:
		default:
			fmt.Print(response)
		}

		fmt.Println()
	}
}

// // Extracts Checksum-Checker and runs it
// func ChecksumCheck() {
// }

// Method to send a GET request to confirm internet connection
func UseGetRequest(domain string, errChan chan error) error {
	response, err := http.Get("https://" + domain)
	// Print DEBUG statements, if applicable
	DebugVariablePrint("response", false, false, -1, "null", nil, nil, response)
	DebugVariablePrint("err", false, false, -1, "null", err, nil, nil)
	// If successgul, returns nil
	return err
}

// Tests internet connectivity using PING (This will eventually be replaced with built-in Go tools)
func NetworkTest(domain string, errChan chan error) {
	switch domain {
	case "N/A":
		errChan <- nil
	default:
		// Temporary solution is using ping before using go's built-in packages
		// Forward result of UsePing() method to errChan
		errChan <- UseGetRequest(domain, errChan)
	}
}

// Method to abstract creation of a 2D slice
func ReturnSliceCreator(commandsAmount int, tokenCount int) [][]string {
	// Make new empty 2D slice
	returnSlices := make([][]string, commandsAmount)
	// Initialize each section with a new slice
	for i := range returnSlices {
		returnSlices[i] = make([]string, tokenCount)
	}

	// Return generated returnSlices[][] slice
	return returnSlices
}

// Method to find appropriate actions with specific package managers
// New package managers are added here!
func PkgManagerActions(pkgNum int, official bool) ([][]string, int, int) {
	// Initialise variables
	var returnSlices [][]string
	var commandsAmount int
	var tokenCount int
	switch official {
	// Official package managers
	case true:
		switch pkgNum {
		// Apt package manager
		case 0:
			commandsAmount = 5
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "dist-upgrade"
			returnSlices[2] = []string{"-f", "install"}
			returnSlices[3][0] = "autoremove"
			returnSlices[4][0] = "autoclean"

		// Dnf & Yum package manager
		case 1, 4:
			commandsAmount = 3
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "check-update"
			returnSlices[1][0] = "update"
			returnSlices[2][0] = "autoremove"
		// OpenSUSE immutable
		case 2:
			commandsAmount = 2
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			// returnSlices[0][0] // May need to set something
			returnSlices[1][0] = "patch"
		// Zypper package manager
		case 3:
			commandsAmount = 5
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "list-updates"
			returnSlices[1][0] = "patch-check"
			returnSlices[2][0] = "update"
			returnSlices[3][0] = "patch"
			returnSlices[4][0] = "purge-kernels"
		// Rpm-Ostree
		case 5:
			commandsAmount = 3
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "cancel"
			returnSlices[1] = []string{"upgrade", "--check"}
			returnSlices[2][0] = "upgrade"
		// Apk
		case 6:
			commandsAmount = 3
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "upgrade"
			returnSlices[2][0] = "fix"
		// Clear Linux
		case 7:
			commandsAmount = 2
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "check-update" // Returns exit code 1 if no update is available
			returnSlices[1][0] = "update"
		// Arch Linux
		case 8:
		// OpenBSD
		case 9:
			commandsAmount = 1
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "-Uuvm"
		// FreeBSD
		case 10:
			commandsAmount = 5
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "upgrade"
			returnSlices[2][0] = "autoremove"
			returnSlices[3][0] = "clean"
			returnSlices[4] = []string{"audit", "-F"}
		// Solus Linux
		case 11:
			commandsAmount = 2
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update-repo"
			returnSlices[1][0] = "upgrade"
		// Slackware Linux
		case 12:
			commandsAmount = 4
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "install-new"
			returnSlices[2][0] = "upgrade-all"
			returnSlices[3][0] = "clean-system"
		// Void Linux
		case 13:
		// Winget
		case 14:
			commandsAmount = 1
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "upgrade"
			returnSlices[0][1] = "--all"
		default:
			// This case must NEVER appear. If it does, it is a result of developer error
			fmt.Println("ERROR: PKG NUMBER", pkgNum, "not found")
			fmt.Println("QUITting!!")
			os.Exit(3)
		}

	// Alternative package managers
	case false:
		switch pkgNum {
		// Brew package manager
		case 0:
			commandsAmount = 3
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "upgrade"
			returnSlices[1][1] = "-v"
			returnSlices[2][0] = "cleanup"
			returnSlices[2][1] = "-v"
		// Snap package manager
		case 1:
			commandsAmount = 1
			tokenCount = 1
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "refresh"
		// Chocolatey package manager
		case 2:
			commandsAmount = 1
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "upgrade"
			returnSlices[0][1] = "all"
		// Flatpak package manager
		case 3:
			commandsAmount = 2
			tokenCount = 2
			returnSlices = ReturnSliceCreator(commandsAmount, tokenCount)
			// Set actions
			returnSlices[0][0] = "update"
			returnSlices[1][0] = "uninstall"
			returnSlices[1][1] = "--unused"
		}
	}

	// Return 2D slice, x number and y number
	return returnSlices, commandsAmount, tokenCount
}

// Method to execute updates from specific package managers, depending on the number
func ExecutePkgManagers(pkgNum int, official bool, manFlag bool) {
	// Initialise variables
	var err error
	var stdout []byte
	var finalActionSlice []string
	// Retrieve package manager-specific data
	pkgManActions, actionCount, tokenCount := PkgManagerActions(pkgNum, official)

	// DEBUG statement to see if official manager is used
	DebugVariablePrint("official", true, official, -1, "null", nil, nil, nil)

	// // Iterate, adding command arguments as needed
	for i := 0; i < actionCount; i++ {
		// Clear finalActionSlice for next iteration
		finalActionSlice = []string{}

		// Define what type of package managers to use
		var pkgManToUse string
		switch official {
		case true:
			pkgManToUse = OFFICIAL_PKG_MANAGERS[pkgNum]
		case false:
			pkgManToUse = ALTERNATIVE_PKG_MANAGERS[pkgNum]
		}

		// Add more items to finalActionSlice[] slice, as needed
		switch rootUse {
		case "": // If root, no sudo/doas needed
		default:
			finalActionSlice = append(finalActionSlice, pkgManToUse)
		}

		// Add specific actions for the package manager
		for j := 0; j < tokenCount; j++ {
			switch pkgManActions[i][j] {
			case "": // Skip empty lines
			default:
				// Print debug statement if needed
				switch debugFlag {
				case true:
					fmt.Print(" " + pkgManActions[i][j])
				}

				finalActionSlice = append(finalActionSlice, pkgManActions[i][j])
			}
		}

		// Add "-y" flags as needed
		switch manFlag {
		case false:
			// Add "-y" flag
			finalActionSlice = append(finalActionSlice, "-y" /*manualResponse*/)

			// Remove "-y" flag according to defined exceptions
			// Some commands accept empty "" flags, while others do not.
			//  In this case, it is better to just remove the last flag entirely
			switch official {
			// Official package managers
			case true:
				switch pkgNum {
				// Apt package manager
				case 0:
					switch i {
					case 0: // update
						finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
					}
				// Zypper package manager
				case 3:
					switch i {
					case 0, 1, 4: // list-updates
						finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
					}
				// transactional-update Immutable package manager
				case 2, 5, 6, 14:
					// finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
					finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
				}
			// Alternative package managers
			case false:
				switch pkgNum {
				// Brew package manager
				case 0:
					finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
				// Snap package manager
				case 1:
					switch i {
					case 0: // refresh
						finalActionSlice = finalActionSlice[:len(finalActionSlice)-1]
					}
				}
			}
		}

		// DEBUG statement to check critical variables
		DebugVariablePrint("rootUse", false, false, -1, rootUse, nil, nil, nil)
		DebugVariablePrint("manualResponse", false, false, -1, finalActionSlice[len(finalActionSlice)-1], nil, nil, nil)
		DebugVariablePrint("Slice LENGTH", false, false, len(finalActionSlice), "null", nil, nil, nil)

		// Execute commands, depending on rootUse
		switch rootUse {
		case "":
			stdout, err = exec.Command(pkgManToUse, finalActionSlice...).Output()
		default:
			stdout, err = exec.Command(rootUse, finalActionSlice...).Output()
		}

		// Get error messages, and work accordingly
		switch err {
		case nil:
			fmt.Println(string(stdout))
		default:
			fmt.Println(string(stdout))
			fmt.Println(err)
		}
	}
}

// Method to check for specific package manager
func PkgManCheck(pkgNum int, official bool) bool {
	// Initialise variables
	var stdout []byte
	var err error
	// First, check if official or not
	switch official {
	case true:
		stdout, err = exec.Command(OFFICIAL_PKG_MANAGERS[pkgNum], "--help").Output()
	case false:
		stdout, err = exec.Command(ALTERNATIVE_PKG_MANAGERS[pkgNum], "--help").Output()
	}
	// Return true or false based off of result
	switch err {
	case nil:
		switch official {
		case true:
			DebugVariablePrint("FOUND PACKAGE MANAGER", false, false, -1, OFFICIAL_PKG_MANAGERS[pkgNum], nil, nil, nil)
		case false:
			DebugVariablePrint("FOUND PACKAGE MANAGER", false, false, -1, ALTERNATIVE_PKG_MANAGERS[pkgNum], nil, nil, nil)
		}
		return true
	default:
		DebugVariablePrint("err", false, false, -1, "null", err, nil, nil)
		DebugVariablePrint("stdout", false, false, -1, string(stdout), nil, nil, nil)
		return false
	}
}

// Method to check for existance of package managers
func PkgManBegin(aoFlag bool, ooFlag bool, manFlag bool, yFlag bool) error {
	// DEBUG statement to print parameter Statuses
	DebugVariablePrint("AOFLAG", true, aoFlag, -1, "null", nil, nil, nil)
	DebugVariablePrint("OOFlag", true, ooFlag, -1, "null", nil, nil, nil)
	DebugVariablePrint("MANFLAG", true, manFlag, -1, "null", nil, nil, nil)
	DebugVariablePrint("YUMFLAG", true, yFlag, -1, "null", nil, nil, nil)

	// Initialise varibles
	var typeCheck int = 0
	var typeToIterate int = OF_PKG_NUM
	// var pkgLoop int
	var officialPkgMan bool = true

	// Change default values if aoFlag is active
	switch aoFlag {
	case true:
		typeCheck = 1
		officialPkgMan = false
		typeToIterate = AL_PKG_NUM
	}

	// Loop through package managers, official first, then alternative
	for i := typeCheck; i < 2; i++ {
		// // DEBUG statement to print i status
		DebugVariablePrint("i", false, false, i, "null", nil, nil, nil)
		// pkgLoop = 0
		for i2 := 0; i2 < typeToIterate+1; i2++ {
			// // DEBUG statement to print i2 status
			DebugVariablePrint("i2", false, false, i2, "null", nil, nil, nil)
			// Use specific actions for different managers, when applicable
			switch i2 {
			// In case of missing official package manager, return error
			case typeToIterate:
				if !officialPkgMan && aoFlag {
					// TODO: Figure out system of returning an error in event of all alternative package managers attempted,
					//  but was forced by -ao flag.
					// return errors.New("missing alternative package managers, forced by -ao flag")
				} else if officialPkgMan {
					return errors.New("missing official package manager")
				}
			default:
				result := PkgManCheck(i2, officialPkgMan)
				switch result {
				case true:
					// Add exception for Yum, if Dnf exists
					switch i2 {
					case 1:
						switch yFlag {
						case true:
							// Check if YUM exists
							result := PkgManCheck(4, officialPkgMan)
							switch result {
							case true:
								DebugVariablePrint("USING YUM over DNF", false, false, 01, "null", nil, nil, nil)
								i2 = 4
							case false:
								fmt.Println("-yu / --yum-update flag used, but YUM does NOT exist")
								fmt.Println("Using DNF instead")
							}
						}
					}

					// Execute package managers
					fmt.Print("\t* Using package manager [")
					switch officialPkgMan {
					case true:
						fmt.Print(OFFICIAL_PKG_MANAGERS[i2])
					case false:
						fmt.Print(ALTERNATIVE_PKG_MANAGERS[i2])
					}
					fmt.Println("] on " + OS_TYPE)
					ExecutePkgManagers(i2, officialPkgMan, manFlag)

					// If official package manager, break loop after execution
					switch officialPkgMan {
					case true:
						i2 = typeToIterate + 1
					}
				}
			}
		}

		// Define conditions for next iteration
		switch typeCheck {
		case 0:
			switch ooFlag {
			case true:
				// End loop before checking alternative package managers
				i = 2
			default:
				// Prepare for alternative package managers
				officialPkgMan = false
				typeToIterate = AL_PKG_NUM
			}
		}
	}

	// if everything works, return nil
	return nil
}

// Define actions to take based on flags
func ActionsForFlags(aoFlag bool, ooFlag bool, cdFlag string) error {
	// Initialise variables
	var err error

	// Check if aoFlag and ooFlag are both true
	if aoFlag && ooFlag {
		return errors.New("incompatible arguments [-ao && -oo]")
	}

	// Begin network test
	// // Create a new channel that funnels errors
	errChan := make(chan error, 2)
	fmt.Println("* Testing connection to [raw.githubusercontent.com]")
	switch cdFlag {
	case "N/A":
		NetworkTest("raw.githubusercontent.com", errChan)
		if err = <-errChan; err != nil {
			fmt.Println("!!Error when testing domain [raw.githubusercontent.com]...")
			return err
		} else {
			fmt.Println("* Network test with domain [raw.githubusercontent.com]  successful!")
		}
	default:
		fmt.Println("* Testing connection to [" + cdFlag + "]")
		// Concurrently run two instances of NetworkTest() method
		go NetworkTest("raw.githubusercontent.com", errChan)
		go NetworkTest(cdFlag, errChan)
		// // Loop through both errChannels
		for i := 0; i < 2; i++ {
			if err = <-errChan; err != nil {
				// Define specific error message
				switch i {
				case 0:
					fmt.Println("!!Error when testing domain [raw.githubusercontent.com]")
					return err
				case 1:
					fmt.Println("!!Error when testing domain [" + cdFlag + "]")
					return err
				}
			} else {
				// Define specific success message
				switch i {
				case 0:
					fmt.Println("* Network test with domain [raw.githubusercontent.com]  successful!")
				case 1:
					fmt.Println("* Network test with domain [" + cdFlag + "] successful!")
				}
			}
		}
	}

	// Returns nil if all is well
	return nil
}

// Check if the executing user is root or not, and what tools are available
func IsExecutorRoot(username string) (string, error) {
	// Work, according to OS_TYPE
	switch OS_TYPE {
	case "windows": // Do nothing for now, More research needed
		return "", nil
	default: // UNIX-based
		switch username {
		// If username is "root", then simply continue
		case "root":
			fmt.Println("* Script is run as root")
			return "", nil
		// Otherwise, check for sudo or doit
		default:
			fmt.Println("* Script not executed as root, checking if user " + username + " has sudo/doas permission...")
			// Check SUDO
			stdout, err := exec.Command("sh", "-c", "sudo -l | grep ALL").Output()
			switch err {
			case nil:
				switch string(stdout) {
				case "":
					fmt.Println("* no sudo detected..")
					DebugVariablePrint("stdout", false, false, -1, string(stdout), nil, nil, nil)
					DebugVariablePrint("err", false, false, -1, "null", err, nil, nil)
				default:
					fmt.Println("\t* User " + username + " has sudo permissions, continueing...")
					return "sudo", nil
				}
			default:
				fmt.Println(err)
				fmt.Println("* sudo not found...")
				fmt.Println(stdout)
			}
			// Check DOAS
			err = exec.Command("doas").Run()
			switch err {
			case nil:
				stdout, err = exec.Command("sh", "-c", "groups $(whoami) | grep wheel").Output()
				switch string(stdout) {
				case "":
					fmt.Println("* no doas detected..")
					DebugVariablePrint("stdout", false, false, -1, string(stdout), nil, nil, nil)
					DebugVariablePrint("err", false, false, -1, "null", err, nil, nil)
				default:
					fmt.Println("\t* User " + username + " has doas permissions, continueing...")
					return "doas", nil
				}
			}

			// By this point, neither are found, so throw an error to be caught later
			return "ERROR", errors.New("missing root priviledges")
		}
	}
}

// Method to easily clear screen, using terminal tools
func ClearScreen() {
	// Initialise vairables
	var clearCommand *exec.Cmd
	switch OS_TYPE {
	case "windows":
		clearCommand = exec.Command("cls")
	default:
		// We can assume the default to be UNIX-based
		clearCommand = exec.Command("clear")
	}
	// Saves output of command and sets it as system's
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()
}

// Main method
func main() {
	// Declare variables
	// // Start counting time
	timeBegin := time.Now()

	// Defer exit statement
	defer ExitStatement()

	// // Get flags
	// // // -ma / --manual-all
	allManualShort := flag.Bool("ma", false, "Manually review and approve each update")
	allManualLong := flag.Bool("manual-all", false, "See above")
	// // // -ao / --alt-only
	altOnlyShort := flag.Bool("ao", false, "Skip alternative package managers")
	altOnlyLong := flag.Bool("alt-only", false, "See above")
	// // // -cd / --custom-domain
	customDomainShort := flag.String("cd", "N/A", "Test an additional domain for internet connectivity")
	// customDomainLong := flag.String("custom-domain", "N/A", "See above")
	// // // -oo / --disable-alt-managers
	officialOnlyShort := flag.Bool("oo", false, "Skip updates from any unofficial package managers")
	officialOnlyLong := flag.Bool("official-only", false, "See above")
	// // // -yu / --yum-update
	yumUpdateShort := flag.Bool("yu", false, "Uses legacy Yum instead of Dnf on Red-Hat Linux systems")
	yumUpdateLong := flag.Bool("yum-update", false, "See above")
	// // // -h / --help
	helpShort := flag.Bool("h", false, "Prints help message")
	helpLong := flag.Bool("help", false, "See above")
	// // // -v / --version
	versionShort := flag.Bool("v", false, "Print version and quit script")
	versionLong := flag.Bool("version", false, "See above")
	// // // -f / --flags
	flagsShort := flag.Bool("f", false, "Print all available flags and quit script")
	flagsLong := flag.Bool("flags", false, "Print all available flags and quit script")
	// // // -w / --warranty
	warrantyShort := flag.Bool("w", false, "Print warranty and quit script")
	warrantyLong := flag.Bool("warranty", false, "See above")
	// // // -d / --debug
	debugShort := flag.Bool("d", false, "Print extra debugging statements")
	debugLong := flag.Bool("debug", false, "See above")
	// // // Parse flage
	flag.Parse()
	// // // Combine flags as needed
	allManualFlag := *allManualShort || *allManualLong
	altOnlyFlag := *altOnlyShort || *altOnlyLong
	officialOnlyFlag := *officialOnlyShort || *officialOnlyLong
	yumUpdateFlag := *yumUpdateShort || *yumUpdateLong
	helpFlag := *helpShort || *helpLong
	warrantyFlag := *warrantyShort || *warrantyLong
	versionFlag := *versionShort || *versionLong
	flagsFlag := *flagsShort || *flagsLong
	customDomainFlag := *customDomainShort // TODO: Figure out combination system
	debugFlag = *debugShort || *debugLong

	// // // If informational flags are run (-h, -v, -f, -w), act on those first
	if helpFlag || versionFlag || warrantyFlag || flagsFlag {
		// Run versionFlag
		switch versionFlag {
		case true:
			PrintVersion()
		}
		// Run helpFlag
		switch helpFlag {
		case true:
			PrintHelp(2)
		}
		// Run warrantyFlag
		switch warrantyFlag {
		case true:
			PrintWarranty()
		}
		// Run flagsFlag (if helpFlag is not used)
		if flagsFlag && !helpFlag {
			PrintFlags(0)
		}
		// Exit with error code 2
		ExitStatement()
		os.Exit(0)
	}

	// Get user information
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("!!Username NOT found! :")
		fmt.Println(err)
		os.Exit(3) // TODO: Set up an AllError method
	}
	executingUser := currentUser.Username

	// Clear screen
	ClearScreen()

	// Check for root permissions
	rootUse, err = IsExecutorRoot(executingUser)
	switch err {
	case nil: // Do nothing, continue
	default:
		fmt.Println("!!User [", executingUser, "] does NOT have ROOT priviledges")
		fmt.Println(err)
		os.Exit(1)
	}

	// Take initial actions based on the flags provided, including filtering, printing, etc
	switch ActionsForFlags(altOnlyFlag, officialOnlyFlag, customDomainFlag) {
	case nil: // Do nothing, continue
	default:
		fmt.Println(err)
		os.Exit(1)
	}

	// Write status on allManualFlag variable
	switch debugFlag {
	case true:
		fmt.Println("DEBUG= allManualFlag:", allManualFlag)
	}

	// Run package manager checker/runner
	pkgManErr := PkgManBegin(altOnlyFlag, officialOnlyFlag, allManualFlag, yumUpdateFlag)
	switch pkgManErr {
	case nil:
	default:
		fmt.Println("!!", pkgManErr)
		os.Exit(1)
	}

	// Print finishing time
	fmt.Println(time.Since(timeBegin))
}
