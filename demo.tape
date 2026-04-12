Output docs/demo.gif
Set Width 900
Set Height 550
Set FontSize 15
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Padding 20
Set Shell "zsh"

# Hidden setup
Hide
Type "source ~/.zshrc"
Enter
Sleep 500ms
Type `printf '127.0.0.1  localhost\n::1        localhost\n192.168.1.50  old-server.local\n' > /tmp/demo.hosts`
Enter
Sleep 300ms
Type `rm -f ~/.hostfile/origin.hostfile ~/.hostfile/my-config.hostfile`
Enter
Sleep 300ms
Type `export HOSTFILE__HOSTS_FILE=/tmp/demo.hosts`
Enter
Sleep 300ms
Type `alias hostfile='./hostfile'`
Enter
Sleep 300ms
Type `export PROMPT='$ '`
Enter
Sleep 300ms
Type `export RPROMPT=''`
Enter
Sleep 300ms
Show
Hide
Type "clear"
Enter
Sleep 500ms
Show

# 1. Init
Type "# Take over existing hosts file"
Enter
Sleep 1.5s
Type "hostfile init"
Enter
Sleep 500ms
Type "y"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

# 2. Add entries
Type ""
Enter
Type "# Add new entries"
Enter
Sleep 1.5s
Type "hostfile add 192.168.1.100 web.local api.local"
Enter
Sleep 2s

Type "hostfile add 10.0.0.1 db.local"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

Type "hostfile show --json"
Enter
Sleep 3s

# 2.5 JSON pipeline
Type ""
Enter
Type "# Import from JSON"
Enter
Sleep 1.5s
Type `echo '{"172.16.0.1":["cache.local"]}' | hostfile merge -`
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

# 3. Search
Type ""
Enter
Type "# Search by IP or domain"
Enter
Sleep 1.5s
Type "hostfile search 192.168.1.100"
Enter
Sleep 3s

Type "hostfile search web.local"
Enter
Sleep 3s

# 4. Disable
Type ""
Enter
Type "# Disable a domain"
Enter
Sleep 1.5s
Type "hostfile disable api.local"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

# 5. Enable
Type ""
Enter
Type "# Enable it back"
Enter
Sleep 1.5s
Type "hostfile enable api.local"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

# 6. Save
Type ""
Enter
Type "# Save a snapshot"
Enter
Sleep 1.5s
Type "hostfile save my-config"
Enter
Sleep 2.5s

# 7. Remove
Type ""
Enter
Type "# Remove an entry"
Enter
Sleep 1.5s
Type "hostfile remove api.local"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 3s

# 8. Load snapshot
Type ""
Enter
Type "# Restore from snapshot"
Enter
Sleep 1.5s
Type "hostfile load my-config"
Enter
Sleep 2s

Type "hostfile show"
Enter
Sleep 4s
