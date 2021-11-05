#!/bin/bash
mkdir -p ./static/tvlogo/
urls=`cat logo.m3u|grep logo|awk -v RS=' ' '$0 ~ /^tvg-(logo|name)/ {print}' | awk -v RS='"' 'NR%2==0{print}'|awk '{if(NR%2!=0)ORS=" ";else ORS="\n"}1'`
echo -e "$urls"|while read -r url
do
    echo "$url"
    filetype=`echo "$url"| awk '{print $1}'|awk -v RS="." 'END{print}'`
    wget `echo "$url"| awk '{print $1}'` -O "./static/tvlogo/`echo "$url"| awk '{print $2}'`.${filetype}" 
done