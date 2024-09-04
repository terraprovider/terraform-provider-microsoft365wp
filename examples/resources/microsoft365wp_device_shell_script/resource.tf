terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}

/*
.env
export ARM_TENANT_ID='...'
export ARM_CLIENT_ID='...'
export ARM_CLIENT_SECRET='...'
*/


resource "microsoft365wp_device_shell_script" "test" {
  display_name = "TF Test"

  run_as_account                = "system"
  block_execution_notifications = true
  execution_frequency           = "PT1H"
  retry_count                   = 3

  file_name      = "openCompanyPortal.sh"
  script_content = <<EOT
#!/bin/bash
#set -x

############################################################################################
##
## Script to start Company Portal after DEP enrollment
##
###########################################

## Copyright (c) 2020 Microsoft Corp. All rights reserved.
## Scripts are not supported under any Microsoft standard support program or service. The scripts are provided AS IS without warranty of any kind.
## Microsoft disclaims all implied warranties including, without limitation, any implied warranties of merchantability or of fitness for a
## particular purpose. The entire risk arising out of the use or performance of the scripts and documentation remains with you. In no event shall
## Microsoft, its authors, or anyone else involved in the creation, production, or delivery of the scripts be liable for any damages whatsoever
## (including, without limitation, damages for loss of business profits, business interruption, loss of business information, or other pecuniary
## loss) arising out of the use of or inability to use the sample scripts or documentation, even if Microsoft has been advised of the possibility
## of such damages.
## Feedback: neiljohn@microsoft.com



# Define variables
log="$HOME/StartCompanyPortal.log"
appname="StartCompanyPortal"
startCompanyPortalifADE="true"
consoleuser=$(ls -l /dev/console | awk '{ print $3 }')

exec &> >(tee -a "$log")

if [[ -f "$HOME/Library/Logs/launchCompanyPortal" ]]; then

  echo "$(date) | Script has already run, nothing to do"
  exit 0

fi


echo ""
echo "##############################################################"
echo "# $(date) | Starting install of $appname"
echo "############################################################"
echo ""

waitForDesktop


# If this is an ADE enrolled device (DEP) we should launch the Company Portal for the end user to complete registration
if [ "$startCompanyPortalifADE" = true ]; then
  echo "$(date) | Checking MDM Profile Type"
  profiles status -type enrollment | grep "Enrolled via DEP: Yes"
  if [ ! $? == 0 ]; then
    echo "$(date) | This device is not ABM managed, exiting"
    echo "$(date) | Writing completion lock to [~/Library/Logs/launchCompanyPortal]"
	touch "$HOME/Library/Logs/launchCompanyPortal"
	exit 0;
  else
	echo "$(date) | Device is ABM Managed. launching Company Portal"
	echo "$(date) | Writing completion lock to [~/Library/Logs/launchCompanyPortal]"
	touch "$HOME/Library/Logs/launchCompanyPortal"
	open "/Applications/Company Portal.app"
  fi
fi
EOT

  assignments = [
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
