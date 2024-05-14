{
  description = "Make kubernetes API changes between versions as transparent as glass.";

  # Flake inputs
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      # Systems supported
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      # Kubeglass package
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGoModule rec {
          pname = "kubeglass";
            version = "0.0.1";
            src = ./.;

            doCheck = false;

            buildFlags = [
              "-trimpath"
            ];

            vendorHash = "sha256-t564Kv23MPodQ2jbj+sjFBVJW3Jj9C3ezIKWexYcXAc=";

            meta = with nixpkgs.lib; {
              description = "Machine Agent";
              platforms = platforms.unix;
            };
        };
      });
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            go_1_22
            gotools # Go tools like goimports, godoc, and others
          ];
        };
      });
    };
}
