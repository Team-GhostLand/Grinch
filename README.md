# Grinch
`~ Grinch: Git x Modrinth ~`
This project aims to improve collaborateability when authoring Mrpacks using Modrinth's own launcher (by making it less of a pain to host Mrpacks on GitHub - or other on Git hosts), as well as generally improve the experience of authoring modpacks on Modrinth (by introducing proper support for server overrides and optional/unsupported mods - which, although can be expressed (according to Mrpack specification), lack any support for modpack creators in Modrinth's own launcher).

## Installation
If on Linux, use Spectre:
```sh
export PROJECT_NAME="grinch" SCRIPT_NAME="scripts/installer" INSTALL_PATH="/bin/grinch" GIT="https://github.com/Team-GhostLand/Grinch.git" && curl -fsSL https://raw.githubusercontent.com/Team-GhostLand/Spectre/master/universal-installer-scaffolding.sh | sudo -E bash
```
For Windows and Mac (BOTH UNTESTED!), we have pre-built binaries in the Releases section. There is no installer, however. Add them to your `$PATH` in whatever way is canonical for your OS.

## Usage
Uhhhh... Good luck. No, seriously. There are no docs *yet*. This was developed for internal use, so our docs are „just asking questions on Discord”. Of course, we plan to make this tool generally availble (Why else go through the hassle of publishing a GitHub release?), but we just don't have the time for documentation writing right now (GhostLand 7 is right around the corner, so that's what the team is focused on - this is also why we gave up halfway through making the docs for Spectre), especially because it's rather unlikely that users will just randomly appear a week after the GitHub repo was set up. Nevertheless, if this DOES happen (or if it's been more than a week - more like a year, and we simply forgot about the docs), please don't hesitate to open an issue on the repo.

*Btw, there are some in-app help messages, but a lot of them won't make sense without understanding the underlying context of Grinch's fundamental concpets.*