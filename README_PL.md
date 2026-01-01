[EN](README.md)|PL

# Grinch
`~ Grinch: Git x Modrinth ~`

## Co to jest?
Celem tego projektu jest ułatwienie kolaboracji nad Mrpackami, przy tworzeniu ich za pomocą oficjalnego launchera Modrinth (poprzez eliminację niezręczności zwykle towarzyszącej przy przesyłaniu plików Mrpack na GitHuba lub inne hostingi Git) oraz dodanie wsparcia dla niektórych funkcjonalności formatu Mrpack (szczególnie, opcjonalnych modów oraz deklaracji zmian/modów czysto-serwerowych lub tylko-na-kliencie), które według oficjalnej specyfikacji ustalonej przez samego Modrinth *powinny* być obsługiwane (np. „Jeśli mod jest opcjonalny - zapytaj o niego przy instalacji”), których to jednak ich oficjalny launcher ironicznie nie wspiera (np. wszystkie opcjonalne mody są automatycznie pobierane, bez pytania; każdy mod jest eksportowany z oznaczeniem wymagania i na kliencie i na serwerze, nawet jeśli nie powinien).

## Ale po co?
\[TODO: Make some info-graphics of modpack-making workflows and put them here]

*Konkretnie dla członków Zespołu GhostLanda: W poście na developer-forum na Discordzie jest opisane rationale za istnieniem tego projektu.*

## Instalacja
Na Linuxie o architekturze ARM/x86 (w każdej formie - natywnie (eg. Ubuntu, Arch, etc.), na Androidzie przez np. Termux, lub Windowsie przez WSL), użyj [Spectre](https://github.com/Team-GhostLand/Spectre):
```sh
export PROJECT_NAME="grinch" SCRIPT_NAME="scripts/installer" INSTALL_PATH="/bin/grinch" GIT="https://github.com/Team-GhostLand/Grinch.git" && curl -fsSL https://raw.githubusercontent.com/Team-GhostLand/Spectre/master/universal-installer-scaffolding.sh | sudo -E bash
```
Dla Maca (NIETESTOWANE!), native Windowsa (testowano i potwierdzono, że NIE DZIAŁA! w tym momencie; prosimy kogoś z Windowsem o pomoc w debugowaniu problemu) lub Linuxa na architekturze RISC-V (nietestowane, ale powinno działać, skoro Grinch nie używa żadnych bibliotek CGo) - mamy pre-built binaries w sekcji [Releases](https://github.com/Team-GhostLand/Grinch/releases/). Tamte wersje nie posiadają jednak instalatora - musisz ręcznie dodać binarki do `$PATH` twojego systemu.

*Konkretnie dla członków Zespołu GhostLanda: W poście na developer-forum na Discordzie jest bardzo szczegółowa instrukcja instalacji na Windowsie (wraz z konfiguracją WSL i wszystkim innym).*

## Usage
Proszę sprawdzić [Wiki na GitHubie](https://github.com/Team-GhostLand/Grinch/wiki) w celu znalezienia dokumentacji.

*Konkretnie dla członków Zespołu GhostLanda: W poście na developer-forum na Discordzie jest częściowa dokumentacja.*