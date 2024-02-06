ip=$(dig -4 TXT +short o-o.myaddr.l.google.com @ns1.google.com | tr -d '"')
port=8096
curl "http://192.168.1.6:5000/update" \
  --header "Content-Type: application/json" \
  --request POST \
  --data @<(cat <<EOF
  {
    "secret": "YOUR-SECRET-HERE",
    "ip": "http://$ip:$port/"
  }
EOF
)
