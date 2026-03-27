{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls    # Go language server (optional, IDE support)
  ];
  shellHook = ''
    echo "DPTR dev shell ready. Go $(go version | awk '{print $3}')"
  '';
}
