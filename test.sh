#!/bin/bash

print_style () {
    if [ "$2" == "info" ] ; then
        COLOR="96m";
    elif [ "$2" == "success" ] ; then
        COLOR="92m";
    elif [ "$2" == "error" ] ; then
        COLOR="93m";
    elif [ "$2" == "danger" ] ; then
        COLOR="91m";
    else #default color
        COLOR="0m";
    fi

    START_COLOR="\e[$COLOR";
    END_COLOR="\e[0m";

    printf "$START_COLOR%b$END_COLOR" "$1\n";
}

info () {
    print_style "$*" "info"
}
success () {
    print_style "$*" "success"
}
error () {
    print_style "$*" "error"
}
danger () {
    print_style "$*" "danger"
}


while read -r url expected_response; do
    response=$(curl -o /dev/null --silent --write-out "%{http_code}" "$url")

    if [ "$response" == "$expected_response" ];
    then
        success "$response $url";
    else
        error "$response $url";
    fi
done < url-list.txt

