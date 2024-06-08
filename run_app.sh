#!/bin/bash

# Debug mode enables logging
DEBUG=on

# Colour (Bold, switch [1 to [0 to make it regular)
readonly RED='\033[1;31m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[1;34m'
readonly NONE='\033[0m'

# Gets real-time Date and Time in predefined format
getDateAndTime() {
  date +"%d/%m/%Y %T"
}

# Logging function (outputs logs if DEBUG var is O
logger() {
  log_level=$1
  message=$2

  if [[ ${DEBUG} == "on" ]]
  then
    case ${1} in
      error )
        printf "$(getDateAndTime) [${RED}ERROR${NONE}] >>> ${RED}${message}${NONE}\n"
        ;;
      info )
        printf "$(getDateAndTime) [${BLUE}INFO${NONE}] >>> ${message}\n"
        ;;
      warn )
        printf "$(getDateAndTime) [${YELLOW}WARNING${NONE}] >>> ${YELLOW}${message}${NONE}\n"
        ;;
      * )
        printf "${log_level}\n"
        ;;
    esac
  fi
}

checkIfConfigTomlExists() {
  file_name=$*
  logger info "Verifying if '$file_name' file exists."
  if [ -f "$file_name" ]; then
    logger info "$file_name file exists."
  else
    logger error "No $file_name file found. Ensure 'config.toml' exists in the root directory (see 'examples.toml')."
    exit 1
  fi
}

#### Run Application ####
readonly FILE_NAME='config.toml'
checkIfConfigTomlExists $FILE_NAME

logger info "Running the application"
make run
logger info "Execution successfully finished"