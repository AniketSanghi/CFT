# CFT
Coder Friendly Translator

-> [ Under Progress ]

Presently It does mount a file system on your RAM which is mirror of already existing file system you supply.

Presenlty its only Read-only.
Write function is partially implemented, You can only append files now.

# How to Run

First open a directory where you want my repository to be cloned at , then 

                      git clone https://github.com/AniketSanghi/CFT.git
                      
Now follow this syntax to run your command

                      ./CloningRef mountlocation fileToMount
                      
- mount location : Provide a location in your system where it will Mount your file.
- fileToMount : This is the file whose mirror you wanna mount in the mount location.

Eg:                  
                    
                    ./CloningRef ~/Desktop/Trial ./TrialFolder
                      
Here I have mounted my TrialFolder which exists in this directory to my Desktop creating a folder named Trial. Now go to your desktop and open the mounted file! What you see? A mirror of your TrialFolder :)
                    
# Just for FUN

Wanna Check the tree Structure of both the repositories from terminal ? 

- For OSX (Using HomeBrew, for other methods GOOGLE!)

		brew install tree
	
- For Linux

		sudo apt install tree
    
Then write the following to see tree structure of any folder
      
                  tree folder_name

