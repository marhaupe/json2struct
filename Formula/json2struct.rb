# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Json2struct < Formula
  desc ""
  homepage ""
  version "1.9.2"

  depends_on "git"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.9.2/json2struct_Darwin_x86_64.tar.gz"
      sha256 "a576360bf68fcd6f83916ab7c6533bbaeeb0b9a97ccbf531ac6555259940e7f0"

      def install
        bin.install "json2struct"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/marhaupe/json2struct/releases/download/v1.9.2/json2struct_Darwin_arm64.tar.gz"
      sha256 "433d7fc57bfcf0318e4e6d92c878fca5f89fffc500d591732ac65b117b680468"

      def install
        bin.install "json2struct"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/marhaupe/json2struct/releases/download/v1.9.2/json2struct_Linux_x86_64.tar.gz"
        sha256 "6253c24dfe07e6409d14f24838d59d8e8f15c13faee841b5f11726b47a92c7b9"

        def install
          bin.install "json2struct"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/marhaupe/json2struct/releases/download/v1.9.2/json2struct_Linux_arm64.tar.gz"
        sha256 "65b8ab9e0226863df644f39df92a1d9fc50956634d3f9dbbcb0bfd4b2036fac9"

        def install
          bin.install "json2struct"
        end
      end
    end
  end
end
