#!/usr/bin/env bash

ncftpget \
-u "${EPG_FTP_USERNAME}" \
-p "${EPG_FTP_PASSWORD}" \
ftp.epgservice.ru ./ "${EPG_DB_PATH}"
