DISTRO="$(grep '^ID=' /etc/os-release | awk -F'='  '{print tolower($2)}')"
DOCKER="$(command -v docker)"
GOLANG="$(command -v go)"
DOLPHINDB_JAEGER="$(command -v dolphindb-jaeger)"

set -e
if [ -x "$DOLPHINDB_JAEGER" ]; then
    echo "dolphindb-jaeger has been installed"
    exit
fi

# install docker
if [ -z "$DOCKER" ]; then
  echo "========== install docker =========="
  case $DISTRO in
    ubuntu|debain)
      curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
      sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
      sudo apt-get update
      sudo apt-get install -y docker-ce
      ;;
    centos)
      sudo yum install -y yum-utils
      sudo yum-config-manager --add-repo "https://download.docker.com/linux/centos/docker-ce.repo"
      sudo yum install -y docker-ce
      ;;
    *)
      echo "unknown linux distribution for installing docker: $DISTRO"
      exit
      ;;
  esac
fi

# install golang
go_path_message=""
if [ -z "$GOLANG" ]; then
  echo "========== install golang =========="
  package="go1.19.1.linux-amd64.tar.gz"
  curl "https://mirrors.ustc.edu.cn/golang/$package" -o "$package"
  sudo tar -C /usr/local -xzf "$package"
  export PATH=$PATH:/usr/local/go/bin
  go_path_message='export PATH=$PATH:/usr/local/go/bin'
fi

# install dolphindb-jaeger
echo "========== install dolphindb-jaeger =========="
GOLANG_VERSION="$(go version | grep -P '\d{1,2}\.\d{1,2}\.\d{1,3}')"
versionLTE() {
  [ "$1" = "$(echo -e "$1\n$2" | sort -V | head -1)" ]
}
tool_path_message=""
if versionLTE "GOLANG_VERSION" "1.16.0"; then
  sudo git clone https://github.com/TangliziGit/dolphindb-jaeger /usr/local/dolphindb-jaeger
  cd /usr/local/dolphindb-jaeger
  sudo go build main.go -o dolphindb-jaeger
  sudo ln -s /usr/local/dolphindb-jaeger/dolphindb-jaeger /usr/bin/dolphindb-jaeger
else
  go env -w GOPROXY=https://goproxy.cn,direct
  go install github.com/TangliziGit/dolphindb-jaeger@latest
  tool_path_message='export PATH=$PATH:'"$(go env GOPATH)"
fi

echo "========== Done ==========\n"

if [ ! -z "$tool_path_message$go_path_message" ]; then
  echo "NOTE: please put scripts below into your ~/.bashrc to find golang commands:"
fi

if [ ! -z "$go_path_message" ]; then
  echo "$go_path_message"
fi
if [ ! -z "$tool_path_message" ]; then
  echo "$tool_path_message"
fi
