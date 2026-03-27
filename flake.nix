{
  description = "DPTR: Daily Personal Terminal Report (Go Native Rewrite)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        packages.default = pkgs.buildGoModule {
          pname = "dptr";
          version = "1.0.0";
          
          src = ./.;

          # Update the vendorHash if go.mod/go.sum changes
          vendorHash = "sha256-mPF9CHAIhFuqrUsVdZHrKH+nBkgO5VZhoBLXI7CBHSQ=";

          meta = with pkgs.lib; {
            description = "Daily Personal Terminal Report - A terminal briefing system";
            homepage = "https://github.com/sidharthify/dptr";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    );
}
