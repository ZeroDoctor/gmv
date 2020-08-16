# gmv

The purpose of this program is to move large sets of files given a [**regular expression**](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions)

## Usage

<code>
    usage: gmv [options] [src dst]

        src - path to source files
        dst - moves files to specified destination folder

    options:
        -c		uses a config file to move file to a folder, 
                    left blank it will assume config is in the same directory
        
        -f		moves files to a specific folder. similar to dst
        -t		gets a list of files that matches users regular expression *untested*
        -h		prints out usages of this program
        -g		generates folders if they do not exisit
        -r		recursively find files in subdirectories
</code>

## Install

In order to install, download the respected files for you system in releases