{
  description = "Make kubernetes API changes between versions as transparent as glass.";

  # Flake inputs
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      # Kubeglass package
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGoModule {
          pname = "kubeglass";
            version = "0.0.1";
            src = ./.;

            doCheck = false;

            buildFlags = [
              "-trimpath"
            ];

            vendorHash = "sha256-+gReBKh+tedt0yexXLOcLTJURk/aWWzIn3coi/utQYM=";

            meta = with nixpkgs.lib; {
              description = "Make kubernetes API changes between versions as transparent as glass.";
              platforms = platforms.all;
            };
        };
      });
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            go
            gotools # Go tools like goimports, godoc, and others
          ];
        };
      });
    };
}
