class Gitlike < Formula
  desc "GitLike CLI with Git-like workflow for developers"
  homepage "https://github.com/bigdog156/gitlike"
  version "1.0.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-arm64"
    sha256 "d85716fd5da1345a849d913957d572ca5fc48eb91f4db7466671f524174af630"
  else
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-amd64"
    sha256 "3c7c8b60e7cb799d75ec4e3a0ec949cd3d3b8c848f7ac2eaa920dba0f1fb8322"
  end

  def install
    bin.install "gitlike-darwin-arm64" => "gitlike" if Hardware::CPU.arm?
    bin.install "gitlike-darwin-amd64" => "gitlike" if Hardware::CPU.intel?
  end

  test do
    system "#{bin}/gitlike", "--help"
    assert_match "1.0.0", shell_output("#{bin}/gitlike --version 2>&1")
  end
end
