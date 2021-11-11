#!/bin/bash
echo "Run inside the container"
echo /shared/jenkins-send -socket-file /shared/jenkins.sock -commit 7329847328972938742 -file-name /myfile.txt
