# README

bl10server is a very application that makes it possible to lock/unlock and track shared-bicycles with a bl10 lock easily and reliably. This server can easily be integrated in your own solution, communicaton with this server is via gRPC.


## How to install?

1. Make sure that you have a system connected to the internet with a working golang compiler installed.
2. clone this repo in "~/go/src"
3. Run "go install"
4. Go to "~/go/bin" and run ./bl10server
5. The server is now running with a open socket on port 9020, and can be reached by the lock via <ip>:9020 or you can setup a domain that is referring to the ip.

## How to configure?
1. Make sure that your lock is fully charged.
2. Send a STATUS# sms to the phones lock, if everythinkg works fine, you should get a reply after a few seconds.
3. Setup the server with the following string: SERVER,0,<ip>,9020,0,#
4. OK should be replied by the lock.
5. Your lock will now be connected to the server.

## How to use it?
1. Make sure grpcurl is installed. https://github.com/fullstorydev/grpcurl
2. Addapt grpcurl_example.sh to use the IMEI number of your lock.
3. Run the grpcurl_example.sh script.
4. Run receive_stream_example.sh in another screen to see what data the lock is sending to the server. When you close the lock for example an update will be send over the stream and after two minutes the coordinates of the lock will be send.
