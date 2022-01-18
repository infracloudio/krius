#!/usr/bin/env bash
#
# Installs Krius the quick way
#
# Requires curl, grep, cut, tar, uname, chmod, mv, rm.

[[ $- = *i* ]] && echo "Don't source this script!" && return 10



check_cmd() {
	command -v "$1" > /dev/null 2>&1
}

check_tools() {
	Tools=("curl" "grep" "cut" "tar" "uname" "chmod" "mv" "rm")
	for tool in ${Tools[*]}; do
		if ! check_cmd $tool; then
			echo "Aborted, missing $tool, sorry!"
			exit 6
		fi
	done
}

install_krius()
{
	trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; exit 1' ERR
	install_path="/usr/local/bin"
	krius_os="unsupported"
	krius_arch="unknown"
	krius_arm=""
	check_tools

	if [[ -n "$PREFIX" ]]; then
		install_path="$PREFIX/bin"
	fi

	# Fall back to /usr/bin if necessary
	if [[ ! -d $install_path ]]; then
		install_path="/usr/bin"
	fi
	# Not every platform has or needs sudo (https://termux.com/linux.html)
	((EUID)) && sudo_cmd="sudo"

	#########################
	# Which OS and version? #
	#########################

	krius_bin="krius"
	krius_dl_ext=".tar.gz"

	# NOTE: `uname -m` is more accurate and universal than `arch`
	# See https://en.wikipedia.org/wiki/Uname
	unamem="$(uname -m)"
	if [[ $unamem == *aarch64* ]]; then
		krius_arch="arm64"
	elif [[ $unamem == *64* ]]; then
		krius_arch="x86_64"
	elif [[ $unamem == *armv5* ]]; then
		krius_arch="arm"
		krius_arm="v5"
	elif [[ $unamem == *armv6l* ]]; then
		krius_arch="arm"
		krius_arm="v6"
	elif [[ $unamem == *armv7l* ]]; then
		krius_arch="arm"
		krius_arm="v7"
	else
		echo "Aborted, unsupported or unknown architecture: $unamem"
		return 2
	fi

	unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
	if [[ $unameu == *DARWIN* ]]; then
		krius_os="darwin"
		version=${vers##*ProductVersion:}
	elif [[ $unameu == *LINUX* ]]; then
		krius_os="Linux"
	elif [[ $unameu == *FREEBSD* ]]; then
		krius_os="freebsd"
	elif [[ $unameu == *OPENBSD* ]]; then
		krius_os="openbsd"
	elif [[ $unameu == *WIN* || $unameu == MSYS* ]]; then
		# Should catch cygwin
		sudo_cmd=""
		krius_os="windows"
		krius_bin=$krius_bin.exe
	else
		echo "Aborted, unsupported or unknown os: $uname"
		return 6
	fi

	########################
	# Download and extract #
	########################

	echo "Downloading krius for ${krius_os}/${krius_arch}${krius_arm}..."
	krius_file="krius_${krius_os}_${krius_arch}${krius_arm}${krius_dl_ext}"
	if [[ "$#" -eq 0 ]]; then
		# get latest release
		krius_tag="v0.1.0"
		krius_version="0.1.0"
	elif [[ "$#" -gt 1 ]]; then
		echo "Too many arguments."
		exit 1
	elif [ -n $1  ]; then
		# try to get passed version
		krius_tag="v$1"
		krius_version=$1
	fi

	krius_url="https://github.com/infracloudio/krius/releases/download/${krius_tag}/krius_${krius_version}_${krius_os}_${krius_arch}${krius_arm}.tar.gz"
	dl="/tmp/$krius_file"
	rm -rf -- "$dl"
	curl -fsSL "$krius_url" -o "$dl"
	echo "Extracting..." 
	case "$krius_file" in
		*.tar.gz) tar -xzf "$dl" -C "$PREFIX/tmp/" "$krius_bin" ;;
	esac
	chmod +x "$PREFIX/tmp/$krius_bin"

	echo "Putting krius in $install_path $krius_bin (may require password)"
	$sudo_cmd cp "$PREFIX/tmp/$krius_bin" "$install_path/$krius_bin"
	$sudo_cmd rm -- "$dl"

	# check installation
	$krius_bin

	echo "Successfully installed"
	trap ERR
	return 0
}

install_krius $@
