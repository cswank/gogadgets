All Gadgets in the system respond to Robot Command Language (RCL) messages.
A RCL command is a string consisting of:

      <action> <location> <device name> <for|to> <command argument> <units>

So, for example::

    turn on living room light
    turn on living room light for 30 minutes

    heat boiler to 22 C


as root                                                  
apt-get update
apt-get install build-essential
apt-get install libzmq3-dev
apt-get install git
mkdir gadgets
adduser craig
chown craig:craig gadgets/
cd gadgets/
sudo adduser craig sudo
hostname greenhouse
su craig
    
git clone git@bitbucket.org:cswank/gogadgets.git
sudo apt-get install golang
export GOPATH=/opt/gadgets/
cd ../src/bitbucket.org/cswank/gogadgets/
go get
cd gogadgets/
go install   