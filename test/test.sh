#workdir linc

#3.1
go build .
sudo ./linc run -it /bin/sh

#3.2,3.3
sudo ./linc run -it -m 400m stress --vm-bytes 800M --vm 1 
top -o %MEM