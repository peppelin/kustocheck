P=$(go tool cover -func coverage.cov | grep total | cut -d \) -f 2 | sed -r 's/%//g' | tr -s \[:space:\] ' ' | sed -r 's/\s//g')
echo "coverage is ${P} %"
if [ "`echo "${P} < 60" | bc`" -eq 1 ];
  then echo "coverage must be at least 70%";
  exit 1;
fi
