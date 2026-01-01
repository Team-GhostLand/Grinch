EN|[PL](README_PL.md)

# Grinch
`~ Grinch: Git x Modrinth ~`

## What is this?
This project aims to improve collaborateability when authoring Mrpacks using Modrinth's own launcher (by making it less of a pain to host Mrpacks on GitHub and other Git hosts), as well as generally improve the experience of authoring modpacks on Modrinth (by introducing proper support for server overrides and optional/unsupported mods - which, although can be expressed (according to Mrpack specification), lack any support for modpack creators in Modrinth's own launcher).

## But why?
\[TODO: Make some info-graphics of modpack-making workflows and put them here]

## Installation
If you're on ARM/x86 Linux (in any form - natively (eg. Ubuntu, Arch, etc.), on Android with something like Termux, or on Windows with WSL), use [Spectre](https://github.com/Team-GhostLand/Spectre):
```sh
export PROJECT_NAME="grinch" SCRIPT_NAME="scripts/installer" INSTALL_PATH="/bin/grinch" GIT="https://github.com/Team-GhostLand/Grinch.git" && curl -fsSL https://raw.githubusercontent.com/Team-GhostLand/Spectre/master/universal-installer-scaffolding.sh | sudo -E bash
```
For Mac (UNTESTED!), native Windows (tested - and CONFIRMED BROKEN! at the moment; someone with a Windows machine willing to help would be appreciated) or RISC-V Linux (untested, but it's likely to work, given how no CGo libraries are used by Grinch) - we also have pre-built binaries in the [Releases](https://github.com/Team-GhostLand/Grinch/releases/) section. There is no installer for those, however. Add them to your `$PATH` in whatever way is canonical for your OS.

## Usage
Please see the [GitHub Wiki](https://github.com/Team-GhostLand/Grinch/wiki) for docs.