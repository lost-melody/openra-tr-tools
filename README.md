# OpenRA Tr Tools

Tools to help translate unit description strings of _OpenRA_.

> Note: this is a temporary work around, before multi language system implemented in _OpenRA_.

## Usage

- Extract strings from rules:

  ```sh
  # read files from `rules` and export strings into `patch.yaml`
  openra-tr-tools extract CAmod/mods/ca/rules/*.yaml -o patch.yaml
  ```

- Translate strings in `patch.yaml`:

  ```sh
  nvim patch.yaml
  ```

- Patch strings to rules:

  ```sh
  # read files from `rules`, patch with `patch.yaml`, and export into `ourput`
  openra-tr-tools patch -p patch.yaml CAmod/mods/ca/rules/*.yaml -o output/
  ```

## Installation

- Install _Go_ toolchain from [go.dev](https://go.dev/dl/).
- Clone and build `openra-tr-tools`:

  ```sh
  git clone https://github.com/lost-melody/openra-tr-tools.git
  cd openra-tr-tools
  go build
  ```

- Install command completions (optional):
  - Bash:

  ```sh
  pacman -S bash-completion
  mkdir -p ~/.local/share/bash-completion/completions/
  openra-tr-tools completion bash >~/.local/share/bash-completion/completions/openra-tr-tools.bash
  ```

  - Fish:

  ```sh
  mkdir -p ~/.config/fish/completions/
  openra-tr-tools completion fish >~/.config/fish/completions/openra-tr-tools.fish
  ```
