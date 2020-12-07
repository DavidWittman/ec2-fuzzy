# This file was generated by GoReleaser. DO NOT EDIT.
class Ec2Fuzzy < Formula
  desc "Fuzzy search EC2 instances and SSH to them"
  homepage "https://github.com/DavidWittman/ec2-fuzzy"
  version "0.0.1"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/DavidWittman/ec2-fuzzy/releases/download/v0.0.1/ec2-fuzzy_0.0.1_Darwin_x86_64.tar.gz"
    sha256 "72651e1b2282b940951ca8f3b949a5644c5b1e6d1e51dd35788942fd0d843b78"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/DavidWittman/ec2-fuzzy/releases/download/v0.0.1/ec2-fuzzy_0.0.1_Linux_x86_64.tar.gz"
    sha256 "e70f2a5c61cc864360615a56612177c2a5d351968259dd4e80f3463ff5177cc4"
  end

  def install
    bin.install "ec2-fuzzy"
  end
end