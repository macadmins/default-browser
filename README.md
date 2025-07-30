# default-browser

This is a tool that will set the default browser for the current user on a macOS device.

## Usage

```shell
/opt/macadmins/bin/default-browser --identifier com.google.chrome
```

To set other browsers as the default, use the following identifiers:

- Google Chrome (the default if no identifier is passed): `com.google.chrome`
- Safari: `com.apple.safari`
- Firefox: `org.mozilla.firefox`
- MS Edge: `com.microsoft.edgemac`

To set the default browser for another user, run within a root context and specify `--user`. The user account must exist.

```shell
/opt/macadmins/bin/default-browser --identifier com.google.chrome --user tim.apple
```

## Known issues

### System Settings may not work correctly

If System Settings doesn't show all sections correctly after running, this tool, restart the machine. This is likely a timing issue with Launch Services, but we haven't reproduced it consistently. If you see that a restart doesn't fix the issue, you way wish to use the `---no-rebuild-launchservices` flag. This will mean that you need to reboot after running the tool.

![System Settings screenshot](assets/system_settings.png)
