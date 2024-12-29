{ pkgs ? import <nixpkgs> {} }:
let
  unstable = import (fetchTarball "https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz") {
  };

  libPath = with pkgs; lib.makeLibraryPath [
  ];

  in
  pkgs.mkShell {

    nativeBuildInputs =  [
      unstable.wails
      pkgs.pkg-config
      # pkgs.wrapGAppsHook3
      # pkgs.autoPatchelfHook
      # pkgs.copyDesktopItems
    ];

    buildInputs = [
      pkgs.webkitgtk_4_1
      # pkgs.libsoup_3
    ];


    shellHook = with pkgs; ''
      export XDG_DATA_DIRS=${gsettings-desktop-schemas}/share/gsettings-schemas/${gsettings-desktop-schemas.name}:${gtk3}/share/gsettings-schemas/${gtk3.name}:$XDG_DATA_DIRS;
      export GIO_MODULE_DIR="${pkgs.glib-networking}/lib/gio/modules/";
      # export GDK_BACKEND=wayland;
    '';

    # NIXOS_OZONE_WL = "1"; # Enable Ozone wayland ?
    LD_LIBRARY_PATH = libPath;
  }
