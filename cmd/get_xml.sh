#!/usr/bin/env bash

ncftpget -u "${EPG_FTP_USERNAME}" -p "${EPG_FTP_PASSWORD}" ftp.epgservice.ru ./ TV_Pack.xml
