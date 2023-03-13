{{- /*
# Copyright 2021 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/}}
#!/bin/bash
function set_proxy() {
  {{- if .proxy }}
    {{- if .proxy.httpProxy }}
  export HTTP_PROXY={{ .proxy.httpProxy | quote }}
  export http_proxy=${HTTP_PROXY}
    {{- end }}
    {{- if .proxy.httpsProxy }}
  export HTTPS_PROXY={{ .proxy.httpsProxy | quote }}
  export https_proxy=${HTTPS_PROXY}
    {{- end }}
    {{- if .proxy.noProxy }}
  export NO_PROXY={{ .proxy.noProxy | join "," | quote }}
  export no_proxy=${NO_PROXY}
    {{- end }}
  {{- else }}
  unset HTTP_PROXY http_proxy HTTPS_PROXY https_proxy NO_PROXY no_proxy
  {{- end }}
}

function bootstrap_debian_based() {
  export LANG=C
  set_proxy
  apt update
  export DEBIAN_FRONTEND=noninteractive
  until apt install jq netcat-openbsd curl -y; do
    echo "Error installing packages"
    apt update
    sleep 10
  done
  mkdir -p /var/lib/bashible/
}

function check_xfs() {
  for FS_NAME in $(mount -l -t xfs | awk '{ print $1 }'); do
    if command -v xfs_info >/dev/null && xfs_info $FS_NAME | grep -q ftype=0; then
       >&2 echo "XFS file system with ftype=0 was found ($FS_NAME). This may cause problems (https://www.suse.com/support/kb/doc/?id=000020068), please fix it and try again."
       exit 1
    fi
  done
}
