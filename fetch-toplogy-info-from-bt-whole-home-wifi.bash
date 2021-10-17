export PASSWORD="$1"
export BT_HOST="192.168.1.3"
export CURL_FLAGS=( -H 'Cookie: popup=1; key=value' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36' -L )

# Log in
export TN=$(curl "http://${BT_HOST}/login.htm" "${CURL_FLAGS[@]}" | pup 'img[src^="data:"] attr{src}' | cut -c 79- | base64 -d)
export TMP_VAL=$(curl "http://${BT_HOST}/cgi/cgi_login.js?_tn=$TN&_t=$(date +%s)&_=$(date +%s)" -H 'Referer: http://mybtdevice.home/login.htm' "${CURL_FLAGS[@]}" | grep tmp_val | sed 's/[^0-9]//g' | tr -d '[:space:]')
export HASHED_PASSWORD=$(echo -n "${PASSWORD}" | md5sum | tr -d '[:space:]-' | shasum --algorithm 512 - | sed 's/[^a-z0-9]//g' | tr -d '[:space:]')
export SUBMITTABLE_HASHED_PASSWORD=$(echo -n "${HASHED_PASSWORD}${TMP_VAL}" | md5sum | tr -d '[:space:]-' | shasum --algorithm 512 - | sed 's/[^a-z0-9]//g' | tr -d '[:space:]')
curl "http://${BT_HOST}/login.cgi" -H 'Referer: http://mybtdevice.home/login.htm' -H 'Content-Type: application/x-www-form-urlencoded' --data-raw "httoken=${TN}&url=&name=&pws=${SUBMITTABLE_HASHED_PASSWORD}" "${CURL_FLAGS[@]}" 1>&2

# Fetch vast majority of info used by web interface
TN=$(curl "http://${BT_HOST}/status_lan_device.htm" "${CURL_FLAGS[@]}" | pup 'img[src^="data:"] attr{src}' | cut -c 79- | base64 -d)
TOPOLOGY_INFO=$(curl "http://${BT_HOST}/cgi/cgi_toplogy_info.js?_tn=$TN&_t=$(date +%s)&_=$(date +%s)" -H 'Referer: http://mybtdevice.home/status_lan_device.htm' "${CURL_FLAGS[@]}")
NODES=$(echo "${TOPOLOGY_INFO}" | grep toplogy_info | sed -E 's/^.+=//g' | sed -E 's/;$//g')
STATIONS=$(echo "${TOPOLOGY_INFO}" | grep station_info | sed -E 's/^.+=//g' | sed -E 's/;$//g')
echo "${NODES} ${STATIONS}" | jq -s add
