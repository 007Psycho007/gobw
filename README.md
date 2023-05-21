# GoBW
Bitwarden TUI written in Go for Linux. GoBW is a wrapper around the Bitwarden CLI. 
It uses the Bubblegum Framework as TUI. 

# Featues
[x] Copy Username and Password
[x] Basic Search Features using the Name of the entry. 

# Usage 
1. Download the Release File.  
2. Unpack the Archive: `tar -xzf <releasefile>.tar.gz`
3. Execute the binary: `./gobw`

# Keys
<kbd>k</kbd>/<kbd>j</kbd> and <kbd>up</kbd>/<kbd>down</kbd>: Move selection up or down
<kbd>h</kbd>/<kbd>l</kbd> and <kbd>left</kbd>/<kbd>right</kbd>: Move to next/previous page
<kbd>/</kbd>: Search through list
<kbd>q</kbd>: Quit Program/Clear clipboard and return to list
<kbd>Alt</kbd><kbd>Enter</kbd>: Copy Username
<kbd>Enter</kbd>: Copy Password

# Security 
This program will never write the password onto the screen or memory and will overwrite the clipboard after the 10 seconds. Be aware that the password might be still be saved in the clipboard history. 
It will however save the Token provided by Bitwarden CLI after logging in or unlocking the vault during runtime. Make sure to quit the program after you are done using it.

# Know Issues
- SSO Login is not directly supported. Login using `bw login --sso` and start gobw afterwards. 
- Entries without a Name are not selectable. Make sure to always set a Name for a entry
 
# Bubblegum Framework
This Project uses the [Bubblegum Framework](https://github.com/charmbracelet/bubbletea) by Charmbracelet
