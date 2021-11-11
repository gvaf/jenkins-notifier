# jenkins-notifier

```
guix environment --ad-hoc go make
```

# Using chroot

```
mkdir -p chrootjail
CMD='guix pack -R -S /bin=bin -S /etc=etc -S /lib=lib bash python coreutils'
${CMD}
cp $(${CMD}) .
tar xvf *gz -C chrootjail/
mkdir ./shared
mkdir -p chrootjail/{dev,proc,sys,shared}
sudo mount --bind /dev chrootjail/dev
sudo mount --bind /proc chrootjail/proc
sudo mount --bind /sys chrootjail/sys
sudo mount --bind ./shared chrootjail/shared
sh build.sh
cp jenkins-send ./shared/
mkdir -p for-jenkins/
./jenkins-notifier -socket-file ./shared/jenkins.sock -output-dir ./for-jenkins

sudo chroot ./chrootjail /bin/bash
source /etc/profile

sudo umount ./chrootjail/dev/
sudo umount ./chrootjail/proc
sudo umount ./chrootjail/sys
sudo umount ./chrootjail/shared
chmod u+rwx -R chrootjail/ && rm -rf chrootjail/
chmod u+rwx *.tar.gz
rm *.tar.gz
rm -rf shared/
rm -rf for-jenkins/

```
