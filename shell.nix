{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs-channels/archive/b58ada326aa612ea1e2fb9a53d550999e94f1985.tar.gz") { } }:
pkgs.mkShell rec {
  buildInputs = with pkgs; [
    libGL
    pkgconfig
    glfw
    xorg.libX11
    xorg.libXrandr
    xorg.libXinerama
    xorg.libXxf86vm
    xorg.libXi
    xorg.libXcursor
    xorg.libXext
  ];
  nativeBuildInputs = [ addOpenGLRunpath makeWrapper ];
}
