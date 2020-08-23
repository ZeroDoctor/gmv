# gmv

The purpose of this program is to move large sets of files given a [**regular expression**](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions)

## Usage

```
    usage: gmv [options] [src dst]

        src - path to source files
        dst - moves files to specified destination folder

    options:
        -c      uses a config file to move file to a folder, 
                    left blank it will assume config is in the same directory
        
        -f      moves files to a specific folder. similar to dst
        -t      gets a list of files that matches users regular expression *untested*
        -h      prints out usages of this program
        -g      generates folders if they do not exisit
        -r      recursively find files in subdirectories
        -i      ignore files in specified folders
        -v      prints out current version
```

To create a folder use the '[]' to mark string as folder. To places file in said folder simply place a regular expression under it.

```ini
    [docs]
    \.docx
    \.doc

    [pics]
    \.jpeg
    \.jpg

    [pics/png]
    \.png
```

This will put any word documents in a folder called 'docs' and will put jpegs in the folder called 'pics' while placing pngs in the 'pics' subfolder called 'png'

### Example

`gmv -g -r -c config.ini -i Pictures -f Output Src`

The above command will (-g) generate missing folders, (-r) recursively search folders for all matching files, (-c) use config.ini files to determine which files to what folder, (-i) ignore the 'Picture' folder when using '-r', (-f) Stores all desired folders and matching files to a folder named 'Output', and (Src) looks for files to match in a folder called 'Src'. 

`gmv -g -r Src Output -c config.ini -i Pictures`

The above command does the same as the command precceding it with the notable expection of not using the '-f' flag.


Note: 'config.ini', 'Picture', 'Output', and 'Src' are all arbitrary file/folder names.

## Install

In order to install, download the respected files for you system in releases