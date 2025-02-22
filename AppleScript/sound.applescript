on run argv
tell application "Background Music"
	-- get the bundleID of every audio application
	set vol of (a reference to (the first audio application whose bundleID is equal to "com.apple.Safari")) to (item 1 of argv)
end tell
end run
