{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.SDL2
    pkgs.SDL2_image
    pkgs.SDL2_ttf
    pkgs.SDL2_mixer
    pkgs.mesa # provides libGL.so and other OpenGL components
    pkgs.pkg-config
  ];

  shellHook = ''
    export CGO_ENABLED=1
  '';
}

