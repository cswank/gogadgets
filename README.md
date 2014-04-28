All Gadgets in the system respond to Robot Command Language (RCL) messages.
A RCL command is a string consisting of:

      <action> <location> <device name> <for|to> <command argument> <units>

So, for example::

    turn on living room light
    turn on living room light for 30 minutes

    heat boiler to 22 C


                                                  
    1  apt-get update
    2  apt-get build-essential
    3  apt-get install build-essential
    4  apt-cache search libzmq
    5  apt-get install libzmq3
    6  apt-get install bzr
    7  git
    8  apt-get install git
    9  cd /opt/
   10  mkdir gadgets
   11  adduser craig
   12  chown craig:craig gadgets/
   13  cd gadgets/
   14  su craig
   15  sudo adduser craig sudo
   16  hostname greenhouse
   17  su craig
   18  history
root@greenhouse:




    
   15  git clone git@bitbucket.org:cswank/gogadgets.git
   16  sudo apt-get install go
   17  sudo apt-get install golang
   18  apt-cache search golang
   19  ls
   20  cd /opt/gadgets/
   21  ls
   22  cd bin/
   23  sudo apt-get install golang
   24  export GOPATH=/opt/gadgets/
   25  ls
   26  cd ../src/bitbucket.org/cswank/gogadgets/
   27  go get
   28  sudo apt-get install libzmq3-dev
   29  go get
   30  cd gogadgets/
   31  go install   