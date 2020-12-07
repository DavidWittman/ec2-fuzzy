# This file was generated by GoReleaser. DO NOT EDIT.
class Ec2Fuzzy < Formula
  desc "Fuzzy search EC2 instances and SSH to them"
  homepage "https://github.com/DavidWittman/ec2-fuzzy"
  version "0.0.2"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/DavidWittman/ec2-fuzzy/releases/download/v0.0.2/ec2-fuzzy_0.0.2_Darwin_x86_64.tar.gz"
    sha256 "d7a6e2eb3f151dd4faf5021cb6e201e613270ff3cd7f935e8cdc7c066528874b"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/DavidWittman/ec2-fuzzy/releases/download/v0.0.2/ec2-fuzzy_0.0.2_Linux_x86_64.tar.gz"
    sha256 "5472adbbe6fa11325cb2cb9804750e4b5a994eff7b2da54f7519d6d665e57ee7"
  end

  def install
    bin.install "ec2-fuzzy"
  end
end
