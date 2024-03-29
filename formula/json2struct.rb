# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Json2struct < Formula
  desc ""
  homepage ""
  version "1.8.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.8.0/json2struct_1.8.0_macOS_arm64.tar.gz"
      sha256 "8c158134bf55c2e07c9b640bbacc96a9bb00eb55fc0873629d82cdb0f9b274df"

      def install
        bin.install "json2struct"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.8.0/json2struct_1.8.0_macOS_64-bit.tar.gz"
      sha256 "5555182b29becb2d0b7723c2825a0c83d2a374611a08d0ca2a8f6865e63150c3"

      def install
        bin.install "json2struct"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.8.0/json2struct_1.8.0_linux_arm64.tar.gz"
      sha256 "435caf7ae0596a72a0cdeddd92774f95bcdb368405323db5acd05fab47af91a7"

      def install
        bin.install "json2struct"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.8.0/json2struct_1.8.0_linux_64-bit.tar.gz"
      sha256 "679d76e693c2647e01a146108fa5843a1984c43b481489c34b6d65287a5d75f1"

      def install
        bin.install "json2struct"
      end
    end
  end

  depends_on "git"
end
