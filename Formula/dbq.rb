class Dbq < Formula
  desc "a database query tui"
  url "https://github.com/jshawl/dbq/archive/refs/tags/v2025.10.04.tar.gz"
  sha256 "d5558cd419c8d46bdc958064cb97f963d1ea793866414c025906ec15033512ed"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"dbq"
  end
end